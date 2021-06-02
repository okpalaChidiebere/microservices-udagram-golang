package feeds

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

// Get all feed items
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	fds, err := AllFeedItems()
	if err != nil {
		log.Printf("Problem getting all feeds: %s", err.Error())

		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	body, _ := json.Marshal(map[string]interface{}{
		"count": len(fds),
		"rows":  fds,
	}) //stringify the go struct
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

//create a feed item
func CreateFeedItemHandler(w http.ResponseWriter, r *http.Request) {
	i, err := PostFeedItem(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	body, _ := json.Marshal(i) //stringify the go struct
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

// get a specific feeditem by Primary Key
func GetFeedItemHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		http.Error(w, http.StatusText(400), http.StatusBadRequest)
		return
	}

	i, err := GetFeedItem(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	body, _ := json.Marshal(i) //stringify the go struct
	w.Header().Set("Content-Type", "application/json")
	w.Write(body) //we return the record to the client in a sensible payload
}

// Get a signed url to put a new item in the bucket
func GetGetSignedUrlHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fn, ok := vars["fileName"]
	if !ok {
		http.Error(w, http.StatusText(400), http.StatusBadRequest)
		return
	}

	url, err := GetGetSignedUrl(fn)
	if err != nil {
		log.Printf("Problem getting signed url for %s: %s", fn, err.Error())

		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	body, _ := json.Marshal(map[string]interface{}{
		"url": url,
	})
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(body) //we return the record to the client in a sensible payload
}

func CheckDbConnectionHandler(w http.ResponseWriter, r *http.Request) {

	//assinging a processId to requests helps debug logs when we have horizontally scalled applications. This will help us differenciate processes when multiple pods are running this exact hblock of code
	pId := uuid.Must(uuid.NewV4(), nil).String()

	ctx := r.Context()

	//This uniqueID we assign to this context helps us trace of the activity through the request when debugging
	ctx = context.WithValue(ctx, "processID", pId)

	err := CheckDbConnection(ctx)
	if err != nil {
		log.Printf("%s - Unable to connect to the database: %s", pId, err.Error())

		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "Unable to connect to the database.")
		return
	}

	log.Printf("%s - Connection has been established successfully", pId)
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "Connection has been established successfully.")
}
