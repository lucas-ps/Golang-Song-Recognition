package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

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
	r.HandleFunc("/tracks/{name}/{audio}", Create).Methods("PUT")
	r.HandleFunc("/tracks/{name}", BadRequest).Methods("PUT")
	r.HandleFunc("/tracks/", NoContent).Methods("PUT")
	/* List */
	r.HandleFunc("/tracks", List).Methods("GET")
	/* Read */
	r.HandleFunc("/tracks/{name}", Read).Methods("GET")
	/* Delete */
	r.HandleFunc("/tracks/{name}", Delete).Methods("DELETE")
	return r
}

func BadRequest(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusBadRequest) } /* 400 */
func NoContent(w http.ResponseWriter, r *http.Request)  { w.WriteHeader(http.StatusNoContent) }  /* 204 */

func Create(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	audio := vars["audio"]

	query, err := db.Prepare("INSERT INTO tracks (name, song) VALUES (?, ?)")
	defer query.Close()

	_, err = query.Exec(name, audio)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) /* 500 */
		return
	}
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
	/*vars := mux.Vars(r)
	  name := vars["name"]*/

	w.WriteHeader(http.StatusOK)                  /* 200 OK, return JSON object with ID and Base64 WAV */
	w.WriteHeader(http.StatusNotFound)            /* 404 Not Found */
	w.WriteHeader(http.StatusInternalServerError) /* 500 Internal Server Error */
}

func Delete(w http.ResponseWriter, r *http.Request) {
	/*vars := mux.Vars(r)
	  name := vars["name"]*/

	w.WriteHeader(http.StatusNoContent) /* 204 No Content */
	w.WriteHeader(http.StatusNotFound)  /* 404 Not Found */
}
