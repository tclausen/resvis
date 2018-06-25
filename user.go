package main

import (	
	"net/http"
	"github.com/gorilla/mux"
	"encoding/json"
	"log"
	"io/ioutil"
	"github.com/dgrijalva/jwt-go"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"toc/db"
	"strconv"
)

type JwtToken struct {
	Token string `json:"token"`
}

//var users []db.User

func GetUsers(w http.ResponseWriter, r *http.Request) *apiError {
	println("----- GetUsers")
	if _, e := ValidateToken(r.Header.Get("Token")); e != nil {
		return &apiError{e, "User validation error", 500}
	}
	var users = db.DB_GetUsers()
	println("users: ")
	json.NewEncoder(w).Encode(users)
	return nil
}

func GetUser(w http.ResponseWriter, r *http.Request) *apiError {
	println("----- GetUser")
	PrintVars(r)
	params := mux.Vars(r)
	id := params["id"]
	var user = db.DB_GetUserFromUsername(id)
	if user.ID == "" {
		return &apiError{nil, "User does not exist", 500}
	}
	json.NewEncoder(w).Encode(user)
	return nil
}

func sanitizeUser(u *db.User) {
	u.Password = ""
	u.Token = ""
}

func GetUserFromToken(w http.ResponseWriter, r *http.Request) *apiError {
	log.Println("----- GetUserFromToken")
//	println("headers", r.Header)
//	for k, v := range r.Header {
//		log.Printf("key=%v, value=%v", k, v)
//	}
	var id string
	var e error 
	if id, e = ValidateToken(r.Header.Get("Token")); e != nil {
		return &apiError{e, "User validation error", 500}
	}
	log.Println("User id: " + id)
	var user = db.DB_GetUserWithRestaurants(id)
	log.Println("Number of visits " + strconv.Itoa(len(user.ResVisits)))
	sanitizeUser(&user)
	log.Println("User id: " + id)
	if UserHasAdminRights(id) {
		db.DB_CreateSetting(user.ID, "Admin", "true")
	}
	db.DB_CreateSetting(user.ID, "DB", db.DBId)
	for k, v := range user.Settings {
		log.Printf("key=%v, value=%v", k, v)
	}
	json.NewEncoder(w).Encode(user)
	return nil
}

func PrintVars(r *http.Request) {
	log.Println("In Vars():")
	for k, v := range mux.Vars(r) {
		log.Printf("key=%v, value=%v", k, v)
	}
}

func DeleteUser(w http.ResponseWriter, r *http.Request) *apiError {
	log.Println("----- DeleteUser")
	if e := ValidateUserHasAdminRights(r.Header.Get("Token")); e != nil {
		return &apiError{e, "User validation error", 500}
	}
	log.Println("DeleteUser")
	params := mux.Vars(r)
	id := params["id"]
	if id == "" {
		return &apiError{nil, "User id empty in delete user", 500}
	}
	db.DB_DeleteUser(id)
	return nil
}

func CreateUser(w http.ResponseWriter, r *http.Request) *apiError {
	if e := ValidateUserHasAdminRights(r.Header.Get("Token")); e != nil {
		return &apiError{e, "User validation error", 500}
	}
	log.Println("----- CreateUser")
	b, _ := ioutil.ReadAll(r.Body)
	log.Println("Create user from ", string(b))
	
	var user db.User
	json.Unmarshal(b, &user)
	if user.UserName == "" {
		return &apiError{nil, "Username not set on user creation", 500}
	}
	if user.Password == "" {
		return &apiError{nil, "Password not set on user creation", 500}
	}
	userDB := db.DB_GetUserFromUsername(user.UserName)
	if userDB.UserName != "" {
		return &apiError{nil, "Username exists", 500}
	}
	newUser := db.DB_AddUser(user, "")
	json.NewEncoder(w).Encode(newUser)
	
	return nil
}

func UserHasAdminRights(userId string) bool {
	// Simplified user rights: Check userId against super admin
	if userId == "0" {
		log.Println("Superadmin has admin rights")
		return true
	}
	// Simplified user rights: Check user name is admin.
	user := db.DB_GetUser(userId)
	log.Println(user.UserName, userId)
	if user.UserName == "admin" || user.UserName == "toc" {
		log.Println("User has admin rights")
		return true
	}
	log.Println("User does not have admin rights")
	return false
}

func ValidateUserHasAdminRights(tokenString string) error {
 	userId, e := ValidateToken(tokenString)
	if e != nil {
		return e
	}
	if UserHasAdminRights(userId) {
		return nil
	}
	return fmt.Errorf("User does not have admin rights")
}

func ValidateToken(tokenString string) (string, error) {
	log.Println("ValidateToken. Token: ", tokenString)
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Error parsing token")
		}
		return []byte("secret"), nil
	})
	if token == nil {
		return "", fmt.Errorf("Unkown token")
	}
	var userId string
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		var user db.User
		mapstructure.Decode(claims, &user)
		log.Println("Got username from token:", user.UserName)
		userDB := db.DB_GetUserFromUsername(user.UserName)
		if userDB.Token != tokenString {
			log.Println("Token not up to date. Got vs expected:")
			log.Println(tokenString)
			log.Println(userDB.Token)
			return "", fmt.Errorf("Token not up to date")
		}
		userId = userDB.ID
	} else {
		return "", fmt.Errorf("Wrong claims in token")
	}
	log.Println("User validated")
	return userId, nil
}

func UserH_Logout(w http.ResponseWriter, r *http.Request) *apiError {
	log.Println("UserH_Logout")
	var userId string
	var e error
	if userId, e = ValidateToken(r.Header.Get("Token")); e != nil {
		return &apiError{e, "User validation error", 500}
	}
	if userId != "" {
		db.DB_SetUserToken(userId, "")
	}
	return nil
}

// Note: Update never touches visits!
func UserH_Update(w http.ResponseWriter, r *http.Request) *apiError {
	log.Println("UserH_Update")
	var e error
	if _, e = ValidateToken(r.Header.Get("Token")); e != nil {
		return &apiError{e, "User validation error", 500}
	}
	b, _ := ioutil.ReadAll(r.Body)
	var user db.User
	json.Unmarshal(b, &user)
	if !db.DB_UpdateUserIgnoreVisits(user) {
		return &apiError{e, "User doesn't exist in update: " + user.ID, 500}
	}
	return nil
}

func User_GetTokenString(username string) (string, *apiError) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
	})
	tokenString, error := token.SignedString([]byte("secret"))
	if error != nil {
		return "", &apiError{error, "Cannot sign token", 500}
	}
	return tokenString, nil
}

func UserH_Login(w http.ResponseWriter, r *http.Request) *apiError {
	log.Println("UserH_Login")
	if r.Body == nil {
		return &apiError{fmt.Errorf("No body in login request"), "User validation error", 500}
	}
	var user db.User
	_ = json.NewDecoder(r.Body).Decode(&user)
	log.Println("Login with user, pass:", user.UserName, user.Password)
	userId, errorUser := db.DB_ValidateUser(user.UserName, user.Password)
	if errorUser != nil {
		return &apiError{errorUser, "User validation error", 403}
	}
	tokenString, error := User_GetTokenString(user.UserName)
	if error != nil {
		return error
	}
	db.DB_SetUserToken(userId, tokenString)
	json.NewEncoder(w).Encode(JwtToken{Token: tokenString})
	return nil
}
