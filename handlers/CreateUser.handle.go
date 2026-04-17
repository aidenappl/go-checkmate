package handlers

import "github.com/aidenappl/go-checkmate/structs"

func CreateUser(db db.Queryable, req structs.User) (*structs.User, error) {
	return &req, nil
}
