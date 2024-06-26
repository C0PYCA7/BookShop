package postgres

import (
	"BookShop/book_service/internal/model"
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"log"
)

func CreateTables(db *sql.DB) error {

	query := `
	CREATE TABLE IF NOT EXISTS author(
	    id SERIAL PRIMARY KEY ,
	    name VARCHAR(55),
	    surname VARCHAR(55) UNIQUE,
	    patronymic VARCHAR(100),
	    birthday DATE);

	CREATE INDEX ON author(surname);

	CREATE TABLE IF NOT EXISTS book(
	    id SERIAL PRIMARY KEY,
	    name VARCHAR(255),
	    genre VARCHAR(55),
	    id_author INT REFERENCES author(id) ON DELETE CASCADE ,
	    date TIMESTAMP,
	    price DECIMAL);

	CREATE INDEX ON book(name);
`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	return nil
}

func (d *Database) AddBook(books *model.AddBook) (int, error) {

	tx, err := d.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("failed to start tx: %w", err)
	}

	var authorId int
	var bookId int
	err = tx.QueryRow("SELECT id FROM author WHERE surname = $1", books.Author).Scan(&authorId)
	if err != nil {
		tx.Rollback()
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("author not found: %w", ErrAuthorNotFound)
		}
		return 0, fmt.Errorf("failed to get author id: %w", ErrInternalServer)
	}

	query := `INSERT INTO book(name, genre, id_author, date, price) 
	VALUES ($1,$2,$3,$4,$5) RETURNING id`

	err = tx.QueryRow(query, books.Name, books.Genre, authorId, books.Date, books.Price).Scan(&bookId)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed to insert book data: %w", ErrInternalServer)
	}

	tx.Commit()

	return bookId, nil
}

func (d *Database) AddAuthor(author *model.AddAuthor) (int, error) {
	query := `INSERT INTO author(name, surname, patronymic, birthday) 
	VALUES ($1, $2, $3, $4) RETURNING id`

	var authorId int

	err := d.db.QueryRow(query, author.Name, author.Surname, author.Patronymic, author.Birthday).Scan(&authorId)
	if err != nil {
		if err.(*pq.Error).Code == "23505" {
			return 0, fmt.Errorf("%w", ErrAuthorExists)
		}
		return 0, fmt.Errorf("%w: %w", ErrInternalServer, err)
	}
	return authorId, nil
}

func (d *Database) GetAuthor(id int) (*model.AuthorInfo, error) {
	tx, err := d.db.Begin()
	if err != nil {
		log.Print(err)
		return nil, err
	}

	var author model.AuthorInfo

	err = tx.QueryRow("SELECT name, surname, patronymic, birthday FROM author WHERE id = $1", id).Scan(&author.Name, &author.Surname, &author.Patronymic, &author.Birthday)
	if err != nil {
		tx.Rollback()
		log.Print(err)
		return nil, fmt.Errorf("failed to get author info: %w", ErrInternalServer)
	}

	rows, err := tx.Query("SELECT name FROM book WHERE id_author = $1", id)
	if err != nil {
		tx.Rollback()
		log.Print(err)
		return nil, fmt.Errorf("failed to get book: %w", ErrInternalServer)
	}
	var arr []string
	for rows.Next() {
		var name string
		err := rows.Scan(&name)
		if err != nil {
			tx.Rollback()
			log.Print(err)
			return nil, fmt.Errorf("failed to get data from rows: %w", ErrInternalServer)
		}
		arr = append(arr, name)
	}
	author.BookList = arr

	tx.Commit()
	return &author, nil
}

func (d *Database) GetBookInfo(id int) (*model.BookInfo, error) {

	query := `SELECT book.name, book.genre, EXTRACT(YEAR FROM book.date) AS year, book.price,
       	author.name AS author_name, author.surname AS author_surname
		FROM book
		JOIN author ON author.id = book.id_author
		WHERE book.id = $1
	`

	var info model.BookInfo

	err := d.db.QueryRow(query, id).Scan(&info.Name, &info.Genre, &info.Year, &info.Price, &info.AuthorName, &info.AuthorSurname)
	if err != nil {
		log.Print(err)
		return nil, fmt.Errorf("failed to get book info: %w", ErrInternalServer)
	}
	return &info, nil
}

func (d *Database) GelAllBooks() ([]model.Book, error) {
	query := `SELECT 
    			book.name AS book_name, 
    			book.genre, 
 				book.price, 
    			author.name AS author_name, 
    			author.surname AS author_surname 
					FROM 
    					book
					JOIN 
    					author 
					ON 
    					author.id = book.id_author;
`

	var books []model.Book

	rows, err := d.db.Query(query)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("failed to get books: %w", ErrInternalServer)
	}

	for rows.Next() {
		var book model.Book
		if err := rows.Scan(&book.Name, &book.Genre, &book.Price, &book.AuthorName, &book.AuthorSurname); err != nil {
			return nil, fmt.Errorf("failed to get all books from rows: %w", ErrInternalServer)
		}
		books = append(books, book)
	}

	return books, nil
}
