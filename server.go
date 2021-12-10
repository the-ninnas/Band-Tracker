package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func main() {
	datosdeArtistas()
	dataRelation()

	fmt.Println("Starting server at port 8080.")

	mux := http.NewServeMux()
	mux.HandleFunc("/", home)

	http.Handle("/templates/", http.StripPrefix("/templates/", http.FileServer(http.Dir("templates"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	
	mux.HandleFunc("/artist", everyArtist)
	log.Fatal(http.ListenAndServe(":8080", mux))

	
}

type Artists []struct {
	ID             int                 `json:"id"`
	Image          string              `json:"image"`
	Name           string              `json:"name"`
	Members        []string            `json:"members"`
	CreationDate   int                 `json:"creationDate"`
	FirstAlbum     string              `json:"firstAlbum"`
	Relations      string              `json:"relations"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

type Relation struct {
	Index []struct {
		Id             int                 `json:"id"`
		DatesLocations map[string][]string `json:"datesLocations"`
	} `json:"index"`
}

var firstData Artists
var secondData Relation

func datosdeArtistas() {
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		fmt.Println("No artists data from request")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err := json.Unmarshal(body, &firstData); err != nil {
		fmt.Println("Not possible to unmarshal JSON")
	}
}

func dataRelation() {

	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/relation")
	if err != nil {
		fmt.Println("No data from request")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err := json.Unmarshal(body, &secondData); err != nil {
		fmt.Println("Not possible to unmarshal JSON")
	}

	for artistIndex, artistInfo := range secondData.Index {
		firstData[artistIndex].DatesLocations = artistInfo.DatesLocations

	}
}

func home(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	ts := template.Must(template.ParseFiles("./templates/index.html"))

	err := ts.Execute(w, firstData)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	} else {
		fmt.Println("Template parsed successfully.")
	}
}

func everyArtist(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/artist" {
		http.NotFound(w, r)
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 || id > len(firstData) {
		http.NotFound(w, r)
		return
	}

	files := []string{
		"./templates/artistPage.html",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
	err = ts.Execute(w, firstData[id-1])

}
