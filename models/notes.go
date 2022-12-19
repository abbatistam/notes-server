package models

// Note struct
type Note struct {
	ID   string `json:"id,omitempty" bson:"_id,omitempty"`
	Name string `json:"name"`
	Body string `json:"body"`
}
