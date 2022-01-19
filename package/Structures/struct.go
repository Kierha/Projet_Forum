package Structures

type USER struct {
	Id       int
	Username string
	Password string
	Mail     string
}

type ONEPOST struct {
	ID          int
	Title       string
	Description string
	Content     string
	Category    string
	Image       string
}
