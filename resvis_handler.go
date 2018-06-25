package main

import (	
	"net/http"
	"github.com/gorilla/mux"
	"encoding/json"
	"log"
	"toc/db"
	"io/ioutil"
	"os"
	"os/exec"
	"io"
	"path/filepath"
	"strconv"
)

func fileExists(f string) bool {
	if _, err := os.Stat(f); err == nil {
		return true
	}
	return false
}

func ResVisH_UploadImage(w http.ResponseWriter, r *http.Request) *apiError {
	var userId string
	var e error
	if userId, e = ValidateToken(r.Header.Get("Token")); e != nil {
		return &apiError{e, "User validation error", 500}
	}
	log.Println("ResVisH_UploadImage, user " + userId)
	e = r.ParseMultipartForm(32 << 20) // 32MB is the default used by FormFile
	if e != nil {
		log.Println("Error parsing form")
		log.Println(e)
		return nil
	}
	fhs := r.MultipartForm.File
	var fileNamesOnDisk = make([]string, 0)
	folder := "../tocsigVue/resvis/static/"
	for key, _ := range fhs {
		handler := fhs[key][0]
		log.Println("Received " + key + ", " + handler.Filename)
		filename := handler.Filename
		ext := filepath.Ext(filename)
		filenameWithoutExt := filename[0:len(filename)-len(ext)]
		index := 1
		for fileExists(folder + filename) {
			filename = filenameWithoutExt + "_" + strconv.Itoa(index) + ext
			index += 1
		}
		if index != 1 {
			log.Println("File exists. Renamed to: " + filename)
		}
		fileNamesOnDisk = append(fileNamesOnDisk, filename)
		f, fileWriteErr := os.OpenFile(folder + filename, os.O_WRONLY|os.O_CREATE, 0666)
		if fileWriteErr != nil {
			log.Println("Error writing file: " + filename)
		}
		file, _ := handler.Open()
		defer file.Close()
		io.Copy(f, file)
		f.Close()
		rescaleAndRotateImage(folder + filename)
	}
	json.NewEncoder(w).Encode(fileNamesOnDisk)
	return nil
}		

func rescaleAndRotateImage(filename string) {
	cmd := exec.Command("convert", filename, "-auto-orient", "-resize", "1280", filename)
	if err := cmd.Run(); err != nil {
	  if err != nil {
			log.Println("ERROR: rotate and re-scale", err.Error())
		}
	  if _, ok := err.(*exec.ExitError); ok {
			log.Println("ERROR: rotate and re-scale", filename)
		}
	} else {
	  // Success
		log.Println("Rotated and re-scaled", filename)
	}
}

func ResVisH_GetPublic(w http.ResponseWriter, r *http.Request) *apiError {
	log.Println("ResVisH_GetPublic")
	params := mux.Vars(r)
	userid := params["userid"]
	visid := params["visid"]
	log.Println("Get public: ", userid, visid)
	ret := db.DB_GetPublicVisit(userid, visid)
	json.NewEncoder(w).Encode(ret)
	return nil
}

func ResVisH_Delete(w http.ResponseWriter, r *http.Request) *apiError {
	var userId string
	var e error
	if userId, e = ValidateToken(r.Header.Get("Token")); e != nil {
		return &apiError{e, "User validation error", 500}
	}
	log.Println("ResVisH_Delete")
	params := mux.Vars(r)
	id := params["id"]
	db.DB_DeleteResVis(userId, id)
	return nil
}

func ResVisH_Create(w http.ResponseWriter, r *http.Request) *apiError {
	var userId string
	var e error
	if userId, e = ValidateToken(r.Header.Get("Token")); e != nil {
		return &apiError{e, "User validation error", 500}
	}
	user := db.DB_GetUser(userId)
	log.Println("ResVisH_Create")
	b, _ := ioutil.ReadAll(r.Body)
	var resvis db.ResVisit
	json.Unmarshal(b, &resvis)
	newResvis := db.DB_AddResVis(user.ID, resvis)
	json.NewEncoder(w).Encode(newResvis)
	return nil
}

func ResVisH_Update(w http.ResponseWriter, r *http.Request) *apiError {
	var userId string
	var e error
	if userId, e = ValidateToken(r.Header.Get("Token")); e != nil {
		return &apiError{e, "User validation error", 500}
	}
	user := db.DB_GetUser(userId)
	log.Println("ResVisH_Update")
	b, _ := ioutil.ReadAll(r.Body)
	var resvis db.ResVisit
	json.Unmarshal(b, &resvis)
	db.DB_UpdateResVis(user.ID, resvis)
	return nil
}
