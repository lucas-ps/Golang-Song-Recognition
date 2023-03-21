package main

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	//"net/url"
	//"fmt"
	"strings"
)

type Audio_json struct {
	Audio string `json:"Audio"`
}

func main() {
	log.Fatal(http.ListenAndServe(":3002", Router()))
}

func Router() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/cooltown", Cooltown).Methods("POST")
	return r
}

func Cooltown(w http.ResponseWriter, r *http.Request) {
	var audio_json Audio_json
	if err := json.NewDecoder(r.Body).Decode(&audio_json); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	/* Getting song name from Search microservice */
	// Prepare params

	postBody, _ := json.Marshal(map[string]string{
		"Audio": audio_json.Audio,
	})
	responseBody := bytes.NewBuffer(postBody)

	// Send POST request
	resp, err := http.Post("http://localhost:3001/search", "application/json", responseBody)
	if err != nil {
		http.Error(w, ("Error sending request to Search microservice: " + err.Error()), http.StatusInternalServerError) /* 500 */
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err, _ := ioutil.ReadAll(resp.Body)
		http.Error(w, "Search microservice error - "+string(err), resp.StatusCode) /* 400 404 500 */
		return
	}

	// Get song name
	var name string
	err = json.NewDecoder(resp.Body).Decode(&name)
	if err != nil {
		http.Error(w, ("Error decoding Search microservice JSON response: " + err.Error()), http.StatusInternalServerError) /* 500 */
	}
	name = strings.ReplaceAll(name, " ", "+")

	/* Get audio from Tracks microservice */
	// Send GET request
	resp, err = http.Get("http://localhost:3000/tracks/" + name)
	if err != nil {
		http.Error(w, ("Error sending request to Tracks microservice: " + err.Error()), http.StatusInternalServerError) /* 500 */
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err, _ := ioutil.ReadAll(resp.Body)
		http.Error(w, "Tracks microservice error - "+string(err), resp.StatusCode) /* 201 204 400 404 500 */
		return
	}

	// Get base64 audio
	var data map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		http.Error(w, ("Error decoding Search microservice JSON response" + err.Error()), http.StatusInternalServerError)
	}

	audio := data["Audio"].(string)

	json.NewEncoder(w).Encode(audio) /* 200 */
}
