package statistics

type CommentStatistics struct {
	Count int `json:"count" bson:"count"`

	// Replied concerns the comments of this group.
	RepliedCount      int      `json:"replied_count" bson:"replied_count"`
	RepliedToComments []string `json:"replied_comments" bson:"replied_comments"`
	RepliedToUsers    []string `json:"replied_users" bson:"replied_users"`
	RepliedRatio      float64  `json:"replied_ratio" bson:"replied_ratio"`

	// Reply concerns replies to the comments of this group.
	ReplyCount    int      `json:"reply_count" bson:"reply_count"`
	ReplyComments []string `json:"reply_comments" bson:"reply_comments"`
	ReplyUsers    []string `json:"reply_users" bson:"reply_users"`
	ReplyRatio    float64  `json:"reply_ratio" bson:"reply_ratio"`
}

type CommentDimensions struct {
	All   *CommentStatistics
	Types map[string]*CommentStatistics `json:"types" bson:",inline"`
}
