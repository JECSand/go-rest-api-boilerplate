package models

import "time"

// Blacklist is a root struct that is used to store the json encoded data for/from a mongodb blacklist doc.
type Blacklist struct {
	Id           string    `json:"id,omitempty"`
	AuthToken    string    `json:"auth_token,omitempty"`
	LastModified time.Time `json:"last_modified,omitempty"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
}
