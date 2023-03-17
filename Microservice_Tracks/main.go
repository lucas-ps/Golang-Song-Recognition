package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

/* Struct to represent the body of a request */
type Body struct {
	Input string `json:"@input"`
}

type Song_json struct {
	ID    string
	Audio string
}

func main() {
	// Connection to SQLite database file in the current directory
	var err error
	db, err = sql.Open("sqlite3", "./tracks.db")
	if err != nil {
		panic(err)
	}

	log.Fatal(http.ListenAndServe(":3000", Router()))
	defer db.Close()
}

func Router() http.Handler {
	r := mux.NewRouter()
	/* Create */
	r.HandleFunc("/tracks/{Id}", Create).Methods("PUT")
	r.HandleFunc("/tracks/", NoContent).Methods("PUT")
	/* List */
	r.HandleFunc("/tracks", List).Methods("GET")
	/* Read */
	r.HandleFunc("/tracks/{Id}", Read).Methods("GET")
	/* Delete */
	r.HandleFunc("/tracks/{Id}", Delete).Methods("DELETE")
	return r
}

func NoContent(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNoContent) } /* 204 */

func Create(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	/* fmt.Println(vars) */

	if len(vars) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	} /* 204 */

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	} /* 400 */

	/* parse the request body as JSON */
	var requestBody Body
	err = json.Unmarshal(body, &requestBody)
	if err != nil {
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
		return
	}

	// extract the value of the @input parameter from the request body
	inputValue := requestBody.Input
	fmt.Println("@input value:", inputValue)

	query, err := db.Prepare("INSERT INTO tracks (name, song) VALUES (?, ?)")
	defer query.Close()

	// TODO: this doesn't work :(
	/*_, err = query.Exec(name, audio)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) /* 500 */
	//return
	//}
	w.WriteHeader(http.StatusCreated)
}

func List(w http.ResponseWriter, r *http.Request) {
	// SQL request
	tracks, err := db.Query("SELECT name FROM tracks")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) /* 500 */
	}
	defer tracks.Close()

	// Collect query results into a string.
	var list string
	for tracks.Next() {
		var name string
		if tracks.Scan(&name) != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError) /* 500 */
			return
		}
		list = list + ", " + name
	}
	if tracks.Err() != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) /* 500 */
		return
	}

	// If no errors have occured, HTTP response with songs list is sent.
	w.Header().Set("Songs", list)
	w.WriteHeader(http.StatusOK) /* 200 */
}

func Read(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["Id"]

	if len(vars) == 0 {
		w.WriteHeader(http.StatusNoContent) /* 204 */
		return
	}

	var song string
	err := db.QueryRow("SELECT song FROM tracks WHERE name = ?", name).Scan(&song)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Track not found", http.StatusNotFound) /* 404 */
			return
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError) /* 500 */
			return
		}
	}

	/* Create JSON object with fetched song, convert to JSON bytes so it can be included in HTTP response*/
	json_response := Song_json{ID: name, Audio: song}
	jsonBytes, err := json.Marshal(json_response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) /* 500 */
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBytes)
	w.WriteHeader(http.StatusOK) /* 200 OK, return JSON object with ID and Base64 WAV */
}

func Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["Id"]

	if name == "" {
		w.WriteHeader(http.StatusBadRequest) /* 400 */
		return
	}

	query, err := db.Exec("DELETE FROM tracks WHERE name = ?", name)
	rowsAffected, err := query.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) /* 500 */
		return
	}
	if rowsAffected == 0 {
		http.Error(w, "Item not found", http.StatusNotFound) /* 404 */
		return
	}
	w.WriteHeader(http.StatusNoContent) /* 204 No Content */
}
