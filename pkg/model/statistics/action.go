package statistics

type ActionStatistics struct {
	Count    int            `json:"count" bson:"count"`
	Users    map[string]int `json:"users" bson:"users"`
	Comments []string       `json:"comments" bson:"comments"`
	Assets   []string       `json:"assets" bson:"assets"`
}

type ActionDimensions struct {
	All   *ActionStatistics
	Types map[string]*ActionStatistics `json:"types" bson:",inline"`
}
