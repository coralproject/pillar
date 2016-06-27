package model

// Dimension will hold all the variables on which aggregations/statistics can be get for the specific type
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
