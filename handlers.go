package main

import (
	"credentials"
	"document"
	"encoding/json"
	"handler"
	"net/http"

	"github.com/gorilla/mux"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func makeMuxRouter() http.Handler {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/product/", credentials.Authenticate(handlerGetDocByID)).Methods("GET")
	muxRouter.HandleFunc("/product", credentials.Authenticate(handlerGetDoc)).Methods("GET")
	muxRouter.HandleFunc("/product", credentials.Authenticate(handlerPostDoc)).Methods("POST")
	muxRouter.HandleFunc("/product", credentials.Authenticate(handlerPutDoc)).Methods("PUT")
	muxRouter.HandleFunc("/product/", credentials.Authenticate(handlerDeleteDoc)).Methods("DELETE")
	muxRouter.HandleFunc("/login", credentials.CreateToken(verify)).Methods("POST")
	return muxRouter
}

//Retrieve all documents from database
func handlerGetDoc(w http.ResponseWriter, r *http.Request) {
	docs, err := product.FindAll()
	if err != nil {
		handler.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	handler.RespondWithJSON(w, http.StatusOK, docs)
}

//Retrieve only document matching query
func handlerGetDocByID(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	doc, err := product.FindByValue(query.Get("doc"))
	if err != nil {
		handler.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	handler.RespondWithJSON(w, http.StatusOK, doc)
}

//Post document to database
func handlerPostDoc(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var doc document.Icecream
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&doc); err != nil {
		handler.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	doc.ID = bson.NewObjectId()
	err := product.Insert(doc)
	switch {
	case mgo.IsDup(err):
		handler.RespondWithError(w, http.StatusConflict, err.Error())
		return
	case err != nil:
		handler.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	handler.RespondWithJSON(w, http.StatusCreated, doc)
}

//Update document in database
func handlerPutDoc(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var doc document.Icecream
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&doc); err != nil {
		handler.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := product.Update(doc); err != nil {
		handler.RespondWithError(w, http.StatusInternalServerError, "Unable to update: "+err.Error())
		return
	}

	handler.RespondWithJSON(w, http.StatusAccepted, map[string]string{"Result": "Successfully updated"})
}

//Delete document from database
func handlerDeleteDoc(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	err := product.Delete(query.Get("doc"))
	if err != nil {
		handler.RespondWithError(w, http.StatusInternalServerError, "Unable to delete: "+err.Error())
		return
	}

	handler.RespondWithJSON(w, http.StatusAccepted, map[string]string{"Result": "Successfully deleted"})
}
