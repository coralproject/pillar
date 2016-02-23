package model

// Author denotes a writer who curated an Asset (or story).
type Author struct {
	ID       string `json:"id" bson:"_id" validate:"required"`
	Name     string `json:"name" bson:"name" validate:"required"`
	URL      string `json:"url,omitempty" bson:"url,omitempty"`
	Twitter  string `json:"twitter,omitempty" bson:"twitter,omitempty"`
	Facebook string `json:"facebook,omitempty" bson:"facebook,omitempty"`
}
