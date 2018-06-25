package main

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"github.com/gorilla/mux"
	"strings"
	"fmt"
	"encoding/json"
	"reflect"
	"toc/db"
)

var router *mux.Router

type tokenStruct struct {
	Token string
}

func assertEqual(t *testing.T, a interface{}, b interface{}, message string) {
	if a == b {
		return
	}
	if len(message) == 0 {
		message = fmt.Sprintf("%v != %v", a, b)
	}
	t.Fatal(message)
}

func checkStatusOKWithFail(t *testing.T, rr *httptest.ResponseRecorder) {
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func postRequest(t *testing.T, url string, body string, token string) (string, error) {
	req, err := http.NewRequest("POST", url, strings.NewReader(body))
	println("E: ", err)
	if err != nil {
		return "", err
	}
	req.Header.Set("Token", token)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	err = checkStatusOK(t, rr)
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(rr.Body.String(), "\n"), nil
}

func getRequestWithFail(t *testing.T, url string, token string) string {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Token", token)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	checkStatusOKWithFail(t, rr)
	return strings.TrimSuffix(rr.Body.String(), "\n")
}

func getRequest(t *testing.T, url string, token string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Token", token)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	e := checkStatusOK(t, rr)
	return rr.Body.String(), e
}

func testSetup() {
	println(" ====================== TEST SETUP =====================")
	router = BuildRouter()
	token, _ := login(nil, "superadmin", "monsterSecret#C#")
	// Add db
	getRequestWithFail(nil, "/api/db/create/UTtest/DBMem", token)
	// Add db
	getRequestWithFail(nil, "/api/db/create/UTtestSql/DBSqlite", token)
	// Select (make active) DB
	getRequestWithFail(nil, "/api/db/UTtestSql", token)
	db.DB_AddTestData()
	getRequestWithFail(nil, "/api/db/UTtest", token)
	db.DB_AddTestData()
}

func testTeardown() {
	println(" ====================== TEST TEARDOWN =====================")
	token, _ := login(nil, "superadmin", "monsterSecret#C#")
	getRequestWithFail(nil, "/api/db/delete/UTtest", token)
}

func checkStatusOK(t *testing.T, rr *httptest.ResponseRecorder) error {
	if status := rr.Code; status != http.StatusOK {
		//t.Errorf("Handler returned wrong status code: got %v want %v",
		//	status, http.StatusOK)
		return fmt.Errorf("Handler returned status code %v", status)
	}
	return nil
}

func login(t *testing.T, u string, p string) (string, error) {
	ret, err := postRequest(t, "/api/user/login", fmt.Sprintf("{\"username\": \"%s\", \"password\":\"%s\"}", u, p), "")
	if err != nil {
		return "", err
	}
	var token tokenStruct
	json.Unmarshal([]byte(ret), &token)
	return token.Token, nil
}

func in_array(v interface{}, in interface{}) (ok bool, i int) {
	val := reflect.Indirect(reflect.ValueOf(in))
	switch val.Kind() {
	case reflect.Slice, reflect.Array:
		for ; i < val.Len(); i++ {
			if ok = v == val.Index(i).Interface(); ok {
				return
			}
		}
	}
	return
}
