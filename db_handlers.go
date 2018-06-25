package main

import (	
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/pusher/pusher-http-go"
	"log"
	"toc/db"
)

// Get list of DBs. First item in list is active DB
func DBH_GetDBList(w http.ResponseWriter, r *http.Request) *apiError {
	log.Println("DBH_GetDBList")
	if e := DBH_IsAdmin(r); e != nil {
		return e
	}
	var dbIds []string
	dbIds = append(dbIds, db.DBId)
	for k, _ := range db.DBIs {
		dbIds = append(dbIds, k)
	}
	json.NewEncoder(w).Encode(dbIds)
	return nil
}

func DBH_Events(w http.ResponseWriter, r *http.Request) *apiError {
	log.Println("DBH_Events")
	if e := DBH_IsAdmin(r); e != nil {
		return e
	}
	events := db.DB_Events()
	json.NewEncoder(w).Encode(events)
	return nil
}

func DBH_SetDB(w http.ResponseWriter, r *http.Request) *apiError {
	if e := DBH_IsAdmin(r); e != nil {
		return e
	}
	params := mux.Vars(r)
	id := params["id"]
	db.DB_SetDBID(id)

	client := pusher.Client{
		AppId: "522816",
		Key: "6e49fd8ecac8ba034462",
		Secret: "0e9eb2d687576664499a",
		Cluster: "eu",
		Secure: true,
	}

	data := map[string]string{"message": "DB changed", "DB": id}
	client.Trigger("my-channel", "my-event", data)
	
	return nil
}

func DBH_DeleteDB(w http.ResponseWriter, r *http.Request) *apiError {
	if e := DBH_IsAdmin(r); e != nil {
		return e
	}
	params := mux.Vars(r)
	id := params["id"]
	if(id == "prod") {
		return &apiError{nil, "Cannot delete prod database", 500}
	}
	db.DB_Delete(id)
	return nil
}

func DBH_IsAdmin(r *http.Request) *apiError {
	if e := ValidateUserHasAdminRights(r.Header.Get("Token")); e != nil {
		return &apiError{e, "User validation error", 500}
	}
	return nil
}

func DBH_CreateDB(w http.ResponseWriter, r *http.Request) *apiError {
	if e := DBH_IsAdmin(r); e != nil {
		return e
	}
	params := mux.Vars(r)
	id := params["id"]
	dbType := params["type"]
	db.DB_Create(id, dbType)
	return nil
}
