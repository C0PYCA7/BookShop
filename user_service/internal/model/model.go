package model

type UserRegistration struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic"`
	Mail       string `json:"email"`
	Login      string `json:"login"`
	Password   string `json:"password"`
}
