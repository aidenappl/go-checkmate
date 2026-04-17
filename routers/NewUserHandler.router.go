package routers

import (
	"encoding/json"
	"net/http"

	"github.com/aidenappl/go-checkmate/responder"
	"github.com/aidenappl/go-checkmate/structs"
)



func NewUserHandler(w http.ResponseWriter, r *http.Request) {
	var req NewUserHandlerRequest

	// Check if body has been parsed correctly
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responder.ParsingError(w, err)
		return
	}

	// Validate fields
	if req.Name == "" {
		responder.ParamError(w, "name")
		return
	}
	if req.Email == "" {
		responder.ParamError(w, "email")
		return
	}
	if req.Password == "" {
		responder.ParamError(w, "password")
		return
	}

	// Set default role if empty
	if req.Role == "" {
		req.Role = structs.UserRoleMember
	}

	// Send to create user
	user, err := handler.CreateUser(db.DB, structs.User{
		Email:    req.Email,
		Name:     &req.Name,
		Role:     req.Role,
		Password: req.Password,
	})

	w.WriteHeader(http.StatusCreated)
}
