package main

import "github.com/google/uuid"

type User struct {
	ID        uuid.UUID `json:id`
	FirstName string    `json:first_name`
	LastName  string    `json:last_name`
	Email     string    `json:email`
	ParentId  uuid.UUID `json:parent_id`
}
