package main

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	// Using UUID instead of integers, so we don't have easily identifiable IDs
	ID        uuid.UUID `json:id`
	FirstName string    `json:first_name`
	LastName  string    `json:last_name`
	Email     string    `json:email`
	ParentId  uuid.UUID `json:parent_id`
	CreatedAt time.Time `json:created_at`
	DeletedAt time.Time `json:deleted_at`
	MergedAt  time.Time `json:merged_at`
}
