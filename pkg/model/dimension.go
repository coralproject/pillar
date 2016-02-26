package model

// Note denotes a note by a user in the system.
type Dimension struct {
	Name         string   `json:"name" bson:"name" validate:"required"`
	Constituents []string `json:"constituents" bson:"constituents"`
}
