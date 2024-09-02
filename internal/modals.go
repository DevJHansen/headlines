package internal

type Headline struct {
	Title      string `json:"title"`
	Content    string `json:"content"`
	Source     string `json:"source"`
	Link       string `json:"link"`
	Media      string `json:"media"`
	CreatedAt  int64  `json:"createdAt"`
	Posted     bool   `json:"posted"`
	DatePosted int64  `json:"datePosted"`
	Deleted    bool   `json:"deleted"`
}
