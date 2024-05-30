package model

import "time"

type AddBook struct {
	Name   string    `json:"name"`
	Genre  string    `json:"genre"`
	Author string    `json:"author"`
	Date   time.Time `json:"date"`
	Price  float64   `json:"price"`
}

type AddAuthor struct {
	Name       string    `json:"name"`
	Surname    string    `json:"surname"`
	Patronymic string    `json:"patronymic"`
	Birthday   time.Time `json:"birthday"`
}

type AuthorInfo struct {
	AddAuthor
	BookList []string
}

type BookInfo struct {
	Name          string
	Genre         string
	Year          int
	Price         float64
	AuthorName    string
	AuthorSurname string
}

type Book struct {
	Name          string
	AuthorName    string
	AuthorSurname string
	Genre         string
	Price         string
}
