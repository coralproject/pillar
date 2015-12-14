package model

// Taxonomy holds all name-value pairs.
type Taxonomy struct {
	Name  string `json:"name" bson:"name"`
	Value string `json:"value" bson:"value"`
}

// Asset denotes an asset in the system e.g. an article or a blog etc.
type Asset struct {
	AssetID    string     `json:"asset_id" bson:"asset_id"`
	SourceID   string     `json:"src_id" bson:"src_id"`
	URL        string     `json:"url" bson:"url" validate:"url"`
	Taxonomies []Taxonomy `json:"taxonomies" bson:"taxonomies"`
}
