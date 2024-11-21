package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

// creating a runtime DB model
type URL struct {
	ID           string    `json:"id"`
	OriginalURL  string    `json:"original_url"`
	ShortURL     string    `json:"short_url"`
	CreationDate time.Time `json:"creation_time"`
}

var urlDB = make(map[string]URL)

func generateShortUrl(OriginalURL string) string {
	hasher := md5.New()
	hasher.Write([]byte(OriginalURL)) //it converts the original URL string to a byte slice

	//aftyer converting the URL to byte slice we will sum the bytes
	data := hasher.Sum(nil) //this will do the sum of the slice bytes

	//now we will encode all the sum to string
	hash := hex.EncodeToString(data) //here hex is used to encode the sum of byteslice to string
	//this will return a long string but we only need a shorter one

	//so now we will return only the first 8 characters
	return hash[:8]
	//we can also return more than 8 or less than 8
	//if we return less than 8 : suppose 5 - then there will be a possibility of two different URLs to have the same same first 5 characters
	//so to be safe we used 8 because it is very rare for the first 8 characters of two different URL be same
}

// now saving the URL structure in the DB
func createURL(originalURL string) string {
	shortURL := generateShortUrl(originalURL)

	id := shortURL //taking the short URL as the id for simplicity
	urlDB[id] = URL{
		ID:           id,
		OriginalURL:  originalURL,
		ShortURL:     shortURL,
		CreationDate: time.Now(),
	}
	return shortURL
}

// Now when someone wants to find the original URL using the shortened URL
func getURL(id string) (URL, error) {
	url, ok := urlDB[id]

	if !ok {
		return URL{}, errors.New("URL not found")
	}

	return url, nil
}

// creating handle
func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "HEllo")
}

func shortUrlHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		URL string `json:"url"`
	}

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid Request body", http.StatusBadRequest)
		return
	}

	shortURL := createURL(data.URL)
	//fmt.Fprintf(w, shortURL)

	response := struct {
		ShortURL string `json:"short_url"`
	}{ShortURL: shortURL}

	w.Header().Set("Contenmt-type", "application/json")
	json.NewEncoder(w).Encode(response)
}

//redirect url handler{}

func redirectUrlHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/redirect/"):]
	url, err := getURL(id)
	if err != nil {
		http.Error(w, "Invalid Request", http.StatusNotFound)
	}

	http.Redirect(w, r, url.OriginalURL, http.StatusFound)
}

func main() {

	//Register the handler function to handle all the requests to the root URL
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/shorten", shortUrlHandler)
	http.HandleFunc("/redirect/", redirectUrlHandler)

	//starting the http server
	fmt.Println("Starting server on the port 8080")
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatal(err)
	}
}
