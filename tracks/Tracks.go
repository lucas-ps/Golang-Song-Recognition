package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

/* Structs needed for database and JSON operations */
type Body struct {
	Input string `json:"@input"`
}

type Track struct {
	Id    string
	Audio string
}

type Track_json struct {
	Id    string `json:"Id"`
	Audio string `json:"Audio"`
}

func main() {
	// Delete existing database
	if _, err := os.Stat("./tracks.db"); err == nil {
		if err := os.Remove("./tracks.db"); err != nil {
			fmt.Printf("Error deleting existing database: %v\n", err)
			return
		}
		//fmt.Printf("Existing database deleted\n")
	}

	// Create database file in the current directory
	var err error
	db, err = sql.Open("sqlite3", "./tracks.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Create tracks table
	_, err = db.Exec("CREATE TABLE Tracks (Name TEXT PRIMARY KEY, Audio TEXT)")
	if err != nil {
		log.Fatal(err)
	}

	// Add test record
	/*_, err = db.Exec("INSERT INTO Tracks (Name, Audio) VALUES ('test', 'test');")
	if err != nil {
		log.Fatal(err)
	}*/

	log.Fatal(http.ListenAndServe(":3000", Router()))
}

func Router() http.Handler {
	r := mux.NewRouter()
	/* Put */
	r.HandleFunc("/tracks/{Id}", Create).Methods("PUT")
	r.HandleFunc("/tracks/", NoContent).Methods("PUT")
	/* Get */
	r.HandleFunc("/tracks", List).Methods("GET")
	r.HandleFunc("/tracks/{Id}", Read).Methods("GET")
	/* Delete */
	r.HandleFunc("/tracks/{Id}", Delete).Methods("DELETE")
	return r
}

func NoContent(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNoContent) } /* 204 */

func Create(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if len(vars) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	} /* 204 */

	var track Track_json
	if err := json.NewDecoder(r.Body).Decode(&track); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) /* 400 */
		return
	}

	query, err := db.Prepare("REPLACE INTO Tracks (Name, Audio) VALUES (?, ?)")
	defer query.Close()

	//fmt.Println(track.Audio)
	_, err = query.Exec(track.Id, track.Audio)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) /* 500 */
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func List(w http.ResponseWriter, r *http.Request) {
	// SQL request
	tracks, err := db.Query("SELECT Name FROM tracks")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) /* 500 */
	}
	defer tracks.Close()

	// Collect query results into a string.
	var track_list []string
	for tracks.Next() {
		var track string
		if tracks.Scan(&track) != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError) /* 500 */
			return
		}
		track_list = append(track_list, track)
	}
	if tracks.Err() != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) /* 500 */
		return
	}

	// If no errors have occured, HTTP response with songs list is sent.
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(track_list)
}

func Read(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["Id"]

	if len(vars) == 0 {
		w.WriteHeader(http.StatusNoContent) /* 204 */
		return
	}

	var song string
	err := db.QueryRow("SELECT Audio FROM Tracks WHERE Name = ?", name).Scan(&song)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Track not found in database", http.StatusNotFound) /* 404 */
			return
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError) /* 500 */
			return
		}
	}

	/* Create JSON object with fetched song, convert to JSON bytes so it can be included in HTTP response*/
	json_response := Track{Id: name, Audio: song}
	jsonBytes, err := json.Marshal(json_response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) /* 500 */
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBytes)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["Id"]

	if name == "" {
		w.WriteHeader(http.StatusBadRequest) /* 400 */
		return
	}

	query, err := db.Exec("DELETE FROM tracks WHERE Name = ?", name)
	rowsAffected, err := query.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) /* 500 */
		return
	}
	if rowsAffected == 0 {
		http.Error(w, "Specified track not found in database", http.StatusNotFound) /* 404 */
		return
	}
	w.WriteHeader(http.StatusNoContent) /* 204 No Content */
}
