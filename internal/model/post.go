package model

import "time"

type Post struct {
	Title    string
	Subtitle string
	Slug     string
	Date     time.Time
	Tags     []string
	HTML     string
	ReadTime int
	Draft    bool
}

type Project struct {
	Name        string
	Description string
	URL         string
	State       string // "live" or "wip"
	Stage       string // e.g. "early", "pre-1.0"
	Kind        string // "oss" or "saas"
}
