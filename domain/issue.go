package domain

type Issue struct {
	ID        string `json:id`
	Title     string `json:title`
	UserCount int    `json:userCount`
	Permalink string `json:permalink`
}
