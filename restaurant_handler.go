package main

import (	
	"net/http"
	"github.com/gorilla/mux"
	"encoding/json"
	"log"
	"toc/db"
	"io/ioutil"
)

func RestaurantH_GetList(w http.ResponseWriter, r *http.Request) *apiError {
	if _, e := ValidateToken(r.Header.Get("Token")); e != nil {
		return &apiError{e, "User validation error", 500}
	}
	log.Println("RestaurantH_GetList")
	var restaurants = db.DB_GetRestaurants()
	json.NewEncoder(w).Encode(restaurants)
	return nil
}

func RestaurantH_Get(w http.ResponseWriter, r *http.Request) *apiError {
	if _, e := ValidateToken(r.Header.Get("Token")); e != nil {
		return &apiError{e, "User validation error", 500}
	}
	log.Println("RestaurantH_Get")
	params := mux.Vars(r)
	id := params["id"]
	restaurant := db.DB_GetRestaurant(id)
	if restaurant.ID == "" {
		return &apiError{nil, "Unknown restaurant", 500}
	}
	json.NewEncoder(w).Encode(restaurant)
	return nil
}

func RestaurantH_Delete(w http.ResponseWriter, r *http.Request) *apiError {
	if _, e := ValidateToken(r.Header.Get("Token")); e != nil {
		return &apiError{e, "User validation error", 500}
	}
	log.Println("RestaurantH_Delete")
	params := mux.Vars(r)
	id := params["id"]
	db.DB_DeleteRestaurant(id)
	return nil
}

func RestaurantH_Update(w http.ResponseWriter, r *http.Request) *apiError {
	log.Println("RestaurantH_Update")
 	var e error
	if _, e = ValidateToken(r.Header.Get("Token")); e != nil {
		return &apiError{e, "User validation error", 500}
	}
	b, _ := ioutil.ReadAll(r.Body)
	var restaurant db.Restaurant
	json.Unmarshal(b, &restaurant)
	if !db.DB_UpdateRestaurant(restaurant) {
		return &apiError{e, "Restaurant doesn't exist in update: " + restaurant.ID, 500}
	}
	return nil
}

func RestaurantH_Create(w http.ResponseWriter, r *http.Request) *apiError {
	var userId string
	var e error
	if userId, e = ValidateToken(r.Header.Get("Token")); e != nil {
		return &apiError{e, "User validation error", 500}
	}
	user := db.DB_GetUser(userId)
	log.Println("RestaurantH_Create")
	b, _ := ioutil.ReadAll(r.Body)
	var restaurant db.Restaurant
	json.Unmarshal(b, &restaurant)
	restaurant.AddedBy = user.UserName
	newRestaurant := db.DB_AddRestaurant(restaurant)
	json.NewEncoder(w).Encode(newRestaurant)
	return nil
}
