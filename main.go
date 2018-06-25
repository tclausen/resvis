package main

import (	
	"log"	
	"github.com/gorilla/mux"
	"net/http"
	"github.com/rs/cors"
	"os"
	"toc/db"
)

func BuildRouter() *mux.Router {
	router := mux.NewRouter()	
	router.Handle("/api/db/{id}", appHandler(DBH_SetDB))
	router.Handle("/api/db/delete/{id}", appHandler(DBH_DeleteDB))
	router.Handle("/api/db/create/{id}/{type}", appHandler(DBH_CreateDB))
	router.Handle("/api/db", appHandler(DBH_GetDBList))
	router.Handle("/api/user", appHandler(GetUsers)).Methods("GET")
	router.Handle("/api/user", appHandler(CreateUser)).Methods("POST")
	router.Handle("/api/user/update", appHandler(UserH_Update)).Methods("POST")
	router.Handle("/api/user/fromToken", appHandler(GetUserFromToken))
	router.Handle("/api/user/login", appHandler(UserH_Login)).Methods("POST")
	router.Handle("/api/user/logout", appHandler(UserH_Logout))
	router.Handle("/api/user/get/{id}", appHandler(GetUser)).Methods("GET")
	router.Handle("/api/user/update", appHandler(DeleteUser)).Methods("POST")
	router.Handle("/api/user/delete/{id}", appHandler(DeleteUser))
	router.Handle("/api/restaurant", appHandler(RestaurantH_GetList)).Methods("GET")
	router.Handle("/api/restaurant", appHandler(RestaurantH_Create)).Methods("POST")
	router.Handle("/api/restaurant/update", appHandler(RestaurantH_Update)).Methods("POST")
	router.Handle("/api/restaurant/{id}", appHandler(RestaurantH_Get)).Methods("GET")
	router.Handle("/api/restaurant/delete/{id}", appHandler(RestaurantH_Delete)).Methods("GET")
	router.Handle("/api/resvis/public/{userid}/{visid}", appHandler(ResVisH_GetPublic)).Methods("GET")
	router.Handle("/api/resvis/add", appHandler(ResVisH_Create)).Methods("POST")
	router.Handle("/api/resvis/delete/{id}", appHandler(ResVisH_Delete)).Methods("GET")
	router.Handle("/api/resvis/update", appHandler(ResVisH_Update)).Methods("POST")
	router.Handle("/api/resvis/uploadimage", appHandler(ResVisH_UploadImage)).Methods("POST")
	router.Handle("/api/events", appHandler(DBH_Events))

	return router
}

type appHandler func(http.ResponseWriter, *http.Request) *apiError

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("------------- NEW CALL")
	if err := fn(w, r); err != nil {
		log.Println("Error sent to client: ", err.Message)
		if err.Error != nil {
			log.Println("Internal error: ", err.Error.Error())
		}
		http.Error(w, err.Message, err.Code)
	}
}

type apiError struct {
	Error   error
	Message string
	Code    int
}

func initLog(fileName string) {
	f, err := os.OpenFile(fileName, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	
	log.SetOutput(f)
	log.Println("---------------- Log init ---------------------")
}

func main() {
	// Write log to file
	//initLog("webapiLog")
	
	db.DB_Init()
	
	router := BuildRouter()
	corsHandler := cors.New(cors.Options{
		AllowedHeaders: []string{"Token"},
	})
	handler := corsHandler.Handler(router)
	println("Listening...")
	log.Fatal(http.ListenAndServe(":8082", handler))
}
