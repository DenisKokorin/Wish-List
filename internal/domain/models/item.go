package models

type Item struct {
	Title       string
	Description string
}

type ItemResponse struct {
	Id          int64
	Title       string
	Description string
}
