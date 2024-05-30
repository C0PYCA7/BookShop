package model

import "time"

type UserRegistration struct {
	Name       string    `json:"name"`
	Surname    string    `json:"surname"`
	Patronymic string    `json:"patronymic"`
	Mail       string    `json:"email"`
	Birthday   time.Time `json:"birthday"`
	Login      string    `json:"login"`
	Password   string    `json:"password"`
}
