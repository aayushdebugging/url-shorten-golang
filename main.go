package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type URL struct {
	ID           string    `json:"id"`
	OriginalURL  string    `json:"original_url"`
	ShortURL     string    `json:"short_url"`
	CreationDate time.Time `json:"creation_date"`
}

var urlDB = make(map[string]URL)

func generateShortURL(OriginalURL string) string {
	hasher := md5.New()
	hasher.Write([]byte(OriginalURL))
	fmt.Println("hasher:", hasher)
	data := hasher.Sum(nil)
	fmt.Println("hasher data", data)
	hash := hex.EncodeToString(data)
	fmt.Println("Encode to String:", hash)
	fmt.Println("Final String:", hash[:8])
	return hash[:8]
}

func createURL(OriginalURL string) string {
	shortURL := generateShortURL(OriginalURL)
	newURL := URL{
		ID:           shortURL,
		OriginalURL:  OriginalURL,
		ShortURL:     shortURL,
		CreationDate: time.Now(),
	}
	urlDB[newURL.ID] = newURL
	return newURL.ShortURL
}

func getURL(id string) (URL, error) {
	url, ok := urlDB[id]
	if !ok {
		return URL{}, errors.New("URL not Found")
	}
	return url, nil
}

func RootPageURL(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World") //for sending in response
}


func ShortURLHandler(w http.ResponseWriter , r *http.Request){
	var data struct{
		URL string `json:"url"`
	}
	err :=json.NewDecoder(r.Body).Decode(&data)
		if err!= nil{
			http.Error(w,"Invalid request Body",http.StatusBadRequest)
			return 
		}
		shortURL := createURL(data.URL)
		//fmt.Fprintf(w,shortURL)
		response:= struct{
			ShortURL string `json:"short_url"`
		}{ShortURL:shortURL}

		w.Header().Set("Content-Type","application/json")
		json.NewEncoder(w).Encode(response)
	}

func redirectURLHandler(w http.ResponseWriter ,r *http.Request){
	id:= r.URL.Path[len("/redirect/"):]
	url,err :=getURL(id)
	if err !=nil{
		http.Error(w,"Invalid Response",http.StatusNotFound)
	}
	http.Redirect(w,r,url.OriginalURL,http.StatusFound)
}
func main() {
	//fmt.Println("Starting URL Shortner......")
	//OriginalURL := "https://www.google.co.in"
	//generateShortURL(OriginalURL)

	//handler function
	http.HandleFunc("/", RootPageURL)
	http.HandleFunc("/shorten", ShortURLHandler)
	http.HandleFunc("/redirect/",redirectURLHandler)

	// Start the server
	fmt.Println("Starting server on PORT 8080....")
	error := http.ListenAndServe(":8080", nil)
	if error != nil {
		fmt.Println("Error on starting the Server:", error)
	}
}
