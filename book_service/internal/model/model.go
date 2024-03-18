package model

import "time"

type AddBook struct {
	Name   string    `json:"name"`
	Genre  string    `json:"genre"`
	Author string    `json:"author"`
	Year   time.Time `json:"year"`
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
	Name          string    `db:"name"`
	Genre         string    `db:"genre"`
	Year          time.Time `db:"year"`
	Price         float64   `db:"price"`
	AuthorName    string    `db:"author_name"`
	AuthorSurname string    `db:"author_surname"`
}
