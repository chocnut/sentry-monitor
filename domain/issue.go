package domain

type Issue struct {
	ID        string `json:id`
	Title     string `json:title`
	UserCount int64  `json:userCount`
	Permalink string `json:permalink`
}
