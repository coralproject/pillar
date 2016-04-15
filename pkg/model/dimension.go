package model

// Note denotes a note by a user in the system.
type Dimension struct {
	Name         string   `json:"name" bson:"name" validate:"required"`
	Constituents []string `json:"constituents" bson:"constituents"`
}

// Validate validates this Model
func (object Dimension) Validate() error {
	//	errs := validate.Struct(object)
	//	if errs != nil {
	//		return fmt.Errorf("%v", errs)
	//	}

	return nil
}
