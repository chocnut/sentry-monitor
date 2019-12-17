package domain

type Issue struct {
	ID        string `json:id`
	Title     string `json:title`
	UserCount int64  `json:userCount`
	Permalink string `json:permalink`
	Count     string `json:count`
	LastSeen  string `json:lastSeen`
}
