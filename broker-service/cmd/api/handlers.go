package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RequestPayload struct {
	Action string `json:"action"`
	Auth  AuthPayload `json:"auth,omitempty"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:  false,
		Message: "Hit the Broker endpoint",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {

	var requestPayload RequestPayload

	err :=  app.readJSON(w,r,&requestPayload)
	if err != nil{
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	switch requestPayload.Action {
	case "auth":
		app.Authenticate(w,requestPayload.Auth)
	default:
		app.errorJSON(w,errors.New("invalid action"),http.StatusBadRequest)
	}
}

func (app *Config) Authenticate(w http.ResponseWriter, a AuthPayload) {

	jsonData, _ := json.MarshalIndent(a, "", "\t")

	request,err := http.NewRequest("POST","http://authentication-service/authenticate", bytes.NewBuffer(jsonData))

	if err != nil{
		app.errorJSON(w,err,http.StatusBadRequest)
		return
	}

	client := &http.Client{}
	resp, err := client.Do(request)
	
	if err != nil{
		app.errorJSON(w,err,http.StatusBadRequest)
		return
	}

	defer resp.Body.Close()

	if(resp.StatusCode == http.StatusUnauthorized){
		app.errorJSON(w,errors.New("invalid email or password"))
		return
	}else if(resp.StatusCode != http.StatusAccepted){
		app.errorJSON(w,errors.New("authentication service error"))
		return
	}

	var jsonFromService jsonResponse

	err = json.NewDecoder(resp.Body).Decode(&jsonFromService)
	if err != nil{
		app.errorJSON(w,err)
		return
	}

	if jsonFromService.Error{
		app.errorJSON(w,err)
		return
	}

	payload := jsonResponse{
		Error: false,
		Message: jsonFromService.Message,
		Data: jsonFromService.Data,
	}

	app.writeJSON(w,http.StatusAccepted,payload)
}
