
package controllers

import (
	u "go-contacts/utils"
	"go-contacts/models"
	"encoding/json"
	"net/http"
)

var SignUp = func(w http.ResponseWriter, req *http.Request) {
	user := &models.User{}
	err := json.NewDecoder(req.Body).Decode(user)
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	u.Respond(w, user.Create())
}

var SignIn = func(w http.ResponseWriter, req *http.Request) {
	user := &models.User{}
	err := json.NewDecoder(req.Body).Decode(user)
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	u.Respond(w, models.Login(user.Email, user.Password))
}
