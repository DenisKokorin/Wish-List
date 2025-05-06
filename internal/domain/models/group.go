package models

type Group struct {
	Id    int64
	Title string
}

type GroupWithMembers struct {
	Id      int64
	Title   string
	Members []int64
}
