package users

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"] //FYI: the primary key for this table is the email. So we expect an email value as the id
	if !ok {
		http.Error(w, http.StatusText(400), http.StatusBadRequest)
		return
	}

	u, err := GetUserByPk(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	body, _ := json.Marshal(u) //stringify the go struct
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func CheckDbConnectionHandler(w http.ResponseWriter, r *http.Request) {

	//Assinging a processId to requests helps debug logs when we have horizontally scalled applications. This will help us differenciate processes when multiple pods are running this exact hblock of code
	pId := uuid.Must(uuid.NewV4(), nil).String()

	ctx := r.Context()

	//This uniqueID we assign to this context helps us trace of the activity through the request when debugging
	ctx = context.WithValue(ctx, "processID", pId)

	err := CheckDbConnection(ctx)
	if err != nil {
		log.Printf("%s - Unable to connect to the database: %s", pId, err.Error())

		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "Unable to connect to the database.")
		return
	}

	log.Printf("%s - Connection has been established successfully", pId)
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "Connection has been established successfully.")
}
