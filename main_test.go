package main

import (
	"testing"
	"os"
	"fmt"
	"encoding/json"
	"toc/db"
)

func TestMain(m *testing.M) {
	testSetup()
	// call flag.Parse() here if TestMain uses flags
	returnValue := m.Run()
	testTeardown()
	os.Exit(returnValue)
}

func TestLoginAdmin(t *testing.T) {
	_, err := login(t, "admin", "wrong code")
	if err == nil {
		t.Fatal("Expected error on wrong password")
	}
	_, err = login(t, "admin", "fire")
	if err != nil {
		t.Fatal("Expected nil in error on correct password")
	}
}

func TestLoginToc(t *testing.T) {
	login(t, "toc", "tt")
}

func TestGetUsers(t *testing.T) {
	token, _ := login(t, "toc", "tt")
	println(token)
	r := getRequestWithFail(t, "/api/user", token)

	println("body:", r)
	// Check the response body is what we expect.
//	expected := `{"alive": true}`
//	if rr.Body.String() != expected {
//		t.Errorf("handler returned unexpected body: got %v want %v",
//			rr.Body.String(), expected)
//	}
}

func TestGetUser(t *testing.T) {
	token, _ := login(t, "toc", "tt")
	getRequestWithFail(t, "/api/user/get/toc", token)
}

func TestGetRestaurants(t *testing.T) {
	token, _ := login(t, "toc", "tt")
	getRequestWithFail(t, "/api/restaurant", token)
}

func TestCreateGetRestaurant(t *testing.T) {
	token, _ := login(t, "toc", "tt")
	// Create restaurant
	testRestaurantName := "Noma"
	restaurantIn := fmt.Sprintf("{ \"name\": \"%s\", \"url\": \"www.noma.dk\" }", testRestaurantName)
	restaurantJson, e := postRequest(t, "/api/restaurant", restaurantIn, token)
	println("response", restaurantJson, e)
	var restaurant db.Restaurant
	json.Unmarshal([]byte(restaurantJson), &restaurant)
	assertEqual(t, restaurant.Name, "Noma", "")
	assertEqual(t, restaurant.Url, "www.noma.dk", "")
	assertEqual(t, restaurant.AddedBy, "toc", "")
	// Get restaurant
	r := getRequestWithFail(t, fmt.Sprintf("/api/restaurant/%s", restaurant.ID), token)
	var restaurant2 db.Restaurant
	json.Unmarshal([]byte(r), &restaurant2)
	println("body:", r)
	assertEqual(t, restaurant2.Name, "Noma", "")
	assertEqual(t, restaurant2.Url, "www.noma.dk", "")
	assertEqual(t, restaurant2.AddedBy, "toc", "")
	// Get restaurant list and check not empty
	r = getRequestWithFail(t, "/api/restaurant", token)
	println("res", r)
	if r == "{}" {
		t.Errorf("Got empty restaurant list after creating one")
	}

	// Delete restaurant
}

func TestNormalUserCannotCreateDB(t *testing.T) {
	// login
	token, _ := login(t, "toc", "tt")

	testDBName := "testUTtoc"
	// Create test DB
	_, e := getRequest(t, fmt.Sprintf("/api/db/create/%s/DBMem", testDBName),token)
	if e == nil {
		t.Fatalf("Did not get error as expected when normal user tried to create DB")
	}
}

func TestAdminCreateDB(t *testing.T) {
	// login
	token, _ := login(t, "admin", "fire")

	// get list of dbs and delete if test DB exists
	testDBName := "testUT"
	dbsJson := getRequestWithFail(t, "/api/db",token)
	println("dbs at start: ", dbsJson)
	var dbs []string
	json.Unmarshal([]byte(dbsJson), &dbs)
	if ok, _ := in_array(testDBName, dbs); ok {
		println("Warning: Test DB exists - deleting")
		getRequest(t, fmt.Sprintf("/api/db/delete/%s", testDBName),token)
	}

	// Create test DB
	getRequest(t, fmt.Sprintf("/api/db/create/%s/DBMem", testDBName),token)

	// Test that DB exists and delete if it does
	dbsJson = getRequestWithFail(t, "/api/db",token)
	println("dbs after creation: ", dbsJson)
	json.Unmarshal([]byte(dbsJson), &dbs)
	if ok, _ := in_array(testDBName, dbs); ok {
		println("Test DB exists - good!")
		getRequest(t, fmt.Sprintf("/api/db/delete/%s", testDBName),token)
	}	else {
		t.Fatalf("Test DB does not exists after creation")
	}

	// Check that DB was deleted
	dbsJson = getRequestWithFail(t, "/api/db",token)
	println("dbs after deletion: ", dbsJson)
	json.Unmarshal([]byte(dbsJson), &dbs)
	println("dbs after: ", dbsJson)
	if ok, _ := in_array(testDBName, dbs); ok {
		t.Fatalf("Test DB exists after deletetion - not good :-(")
	}	
}

func TestAdminCreateDBSQL(t *testing.T) {
	// login
	token, _ := login(t, "admin", "fire")

	// get list of dbs and delete if test DB exists
	testDBName := "testUTSQLInTest"
	dbsJson := getRequestWithFail(t, "/api/db",token)
	println("dbs at start: ", dbsJson)
	var dbs []string
	json.Unmarshal([]byte(dbsJson), &dbs)
	if ok, _ := in_array(testDBName, dbs); ok {
		println("Warning: Test DB exists - deleting")
		getRequest(t, fmt.Sprintf("/api/db/delete/%s", testDBName),token)
	}

	// Create test DB
	getRequest(t, fmt.Sprintf("/api/db/create/%s/DBSqlite", testDBName),token)

	// Test that DB exists and delete if it does
	dbsJson = getRequestWithFail(t, "/api/db",token)
	println("dbs after creation: ", dbsJson)
	json.Unmarshal([]byte(dbsJson), &dbs)
	if ok, _ := in_array(testDBName, dbs); ok {
		println("Test DB exists - good!")
		getRequest(t, fmt.Sprintf("/api/db/delete/%s", testDBName),token)
	}	else {
		t.Fatalf("Test DB does not exists after creation")
	}

	// Check that DB was deleted
	dbsJson = getRequestWithFail(t, "/api/db",token)
	println("dbs after deletion: ", dbsJson)
	json.Unmarshal([]byte(dbsJson), &dbs)
	println("dbs after: ", dbsJson)
	if ok, _ := in_array(testDBName, dbs); ok {
		t.Fatalf("Test DB exists after deletetion - not good :-(")
	}	
}

func TestAdminCreateUser(t *testing.T) {
	// login
	token, _ := login(t, "admin", "fire")

	// get list of users and delete if test user exists
	testUserName := "testUTUserName"
	usersJson := getRequestWithFail(t, "/api/user", token)
	println("users at start: ", usersJson)
	var users map[string]db.User
	json.Unmarshal([]byte(usersJson), &users)
	for _, user := range users {
		if user.UserName == testUserName {
			println("Warning: Test user exists - deleting")
		}
	}

	// Create test user
	user := fmt.Sprintf("{ \"username\": \"%s\", \"password\": \"hej\" }", testUserName)
	r, e := postRequest(t, "/api/user", user, token)
	println("response", r, e)

	// Test that user exists and delete if it does
	usersJson = getRequestWithFail(t, "/api/user", token)
	println("users after creation: ", usersJson)
	users = make(map[string]db.User)
	json.Unmarshal([]byte(usersJson), &users)
	var user1 *db.User
	for _, user := range users {
		if user.UserName == testUserName {
			println("Test user exists as expected")
			user1 = &user
			break
		}
	}
	if user1 != nil {
		println("Test user exists - good! Deleting again")
		getRequest(t, fmt.Sprintf("/api/user/delete/%s", user1.ID), token)
	}	else {
		t.Fatalf("Test user does not exists after creation")
	}

	// Check that user was deleted
	usersJson = getRequestWithFail(t, "/api/user", token)
	println("users after deletion: ", usersJson)
	users = make(map[string]db.User)
	json.Unmarshal([]byte(usersJson), &users)
	for _, user := range users {
		if user.UserName == testUserName {
			t.Fatalf("Test user exists after deletetion - not good :-(")
		}
	}
}

