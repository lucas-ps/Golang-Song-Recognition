package main

import (
	"encoding/json"
	//"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
)

const (
	KEY = "f10fd1bd494f2787828293277f77e884" /* Insert own audd.io API key here */
)

type Audio_json struct {
	Audio string `json:"Audio"`
}

type Response struct {
	Status string `json:"status"`
	Result struct {
		Artist string `json:"artist"`
		Title  string `json:"title"`
	} `json:"result"`
}

func main() {
	log.Fatal(http.ListenAndServe(":3001", Router()))
}

func Router() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/search", Search).Methods("POST")
	return r
}

func Search(w http.ResponseWriter, r *http.Request) {
	var audio_json Audio_json
	if err := json.NewDecoder(r.Body).Decode(&audio_json); err != nil {
		http.Error(w, "Error decoding audio: "+err.Error(), http.StatusBadRequest) /* 400 */
		return
	}

	//fmt.Println(audio_json.Audio)

	// Prepare params
	params := url.Values{}
	params.Add("api_token", KEY)
	params.Add("audio", audio_json.Audio)

	// Send POST request
	resp, err := http.PostForm("https://api.audd.io/?", params)
	if err != nil {
		http.Error(w, ("Error creating request: " + err.Error()), http.StatusInternalServerError) /* 500 */
		return
	}
	defer resp.Body.Close()

	// Decode response
	var data map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Error reading response from audd.io", http.StatusInternalServerError) /* 500 */
		return
	}

	if data["status"].(string) != "success" {
		http.Error(w, "Audd.io was unable to find a song in the provided audio", http.StatusNotFound) /* 404 */
		return
	}

	result := data["result"].(map[string]interface{})
	title := result["title"].(string)

	json.NewEncoder(w).Encode(title) /* 200 */
}
