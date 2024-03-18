package postgres

import (
	"BookShop/user_service/internal/model"
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

func CreateTables(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS users(
	    id SERIAL PRIMARY KEY,
	    name VARCHAR(55),
	    surname VARCHAR(55),
	    patronymic VARCHAR(55),
	    mail VARCHAR(255) UNIQUE ,
	    login VARCHAR(255) UNIQUE ,
	    password VARCHAR(255),
	    permissions VARCHAR(55)
	);

	CREATE INDEX on users(login);`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create query: %w", err)
	}
	return nil
}

func (d *Database) CheckUser(login string) (int, string, string, error) {
	query := `
	SELECT id, password, permissions
	FROM users
	WHERE login = ($1)
	`

	var id int
	var pass string
	var permissions string

	err := d.db.QueryRow(query, login).Scan(&id, &pass, &permissions)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, "", "", ErrUserNotFound
		}
		return 0, "", "", ErrInternalServer
	}

	return id, pass, permissions, nil
}

func (d *Database) CreateUser(user *model.UserRegistration) (int, error) {

	fmt.Println("req in CreateUser", user)

	query := `INSERT INTO users(name, surname, patronymic, mail, login, password)
	VALUES ($1,$2,$3,$4,$5,$6) RETURNING id`

	var userId int

	err := d.db.QueryRow(query, user.Name, user.Surname, user.Patronymic, user.Mail, user.Login, user.Password).Scan(&userId)
	if err != nil {
		if err.(*pq.Error).Code == "23505" {
			return 0, fmt.Errorf("%w", ErrLoginExists)
		}
		return 0, fmt.Errorf("%w: %w", ErrInternalServer, err)
	}

	return userId, nil
}

func (d *Database) DeleteUser(id int) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := d.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", ErrInternalServer)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return ErrInternalServer
	}
	if rows == 0 {
		return fmt.Errorf("%w", ErrUserNotFound)
	}
	return nil
}
