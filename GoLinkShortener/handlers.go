package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
)

/*

 */
type LinkShortnerAPI struct {
	db *Database
}

type UrlMapping struct {
	ShortUrl string `json:shorturl`
	LongUrl  string `json:longurl`
}

type APIResponse struct {
	StatusMessage string `json:statusmessage`
}

func NewUrlLinkShortenerAPI() *LinkShortnerAPI {
	return &LinkShortnerAPI{
		db: NewDatabase(),
	}
}

func (api *LinkShortnerAPI) UrlRoot(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("index.html")
	t.Execute(w, nil)
}

func (api *LinkShortnerAPI) UrlCreate(w http.ResponseWriter, r *http.Request) {
	reqBodyStruct := new(UrlMapping)
	responseEncoder := json.NewEncoder(w)
	// decode json
	if err := json.NewDecoder(r.Body).Decode(&reqBodyStruct); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err := responseEncoder.Encode(&APIResponse{StatusMessage: err.Error()}); err != nil {
			fmt.Fprintf(w, "Error occured while processing post request %v \n", err.Error())
		}
		return
	}
	// write to database
	err := api.db.Put(reqBodyStruct.ShortUrl, reqBodyStruct.LongUrl )
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		if err := responseEncoder.Encode(&APIResponse{StatusMessage: err.Error()}); err != nil {
			fmt.Fprintf(w, "Error %s occured while trying to add the url \n", err.Error())
		}
		return
	}
	responseEncoder.Encode(&APIResponse{StatusMessage: "Ok"})
}

func (api *LinkShortnerAPI) UrlShow(w http.ResponseWriter, r *http.Request) {
	//retrieve the variable from the request
	vars := mux.Vars(r)
	url := vars["shorturl"]
	if len(url) > 0 {
		lUrl, err := api.db.Get(url)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "404 - Entry %s was not found in the database.\n", url)
			return
		}
		//Ensure we are dealing with an absolute path
		http.Redirect(w, r, lUrl, http.StatusFound)
	}
}
