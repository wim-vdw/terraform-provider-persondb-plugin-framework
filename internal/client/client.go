package client

import (
	"database/sql"
	"errors"

	_ "github.com/mattn/go-sqlite3"
)

type Client struct {
	CustomDatabase string
}

func NewClient(databaseFilename string) (*Client, error) {
	c := &Client{
		CustomDatabase: databaseFilename,
	}
	err := c.initDB()
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Client) initDB() error {
	db, err := sql.Open("sqlite3", c.CustomDatabase)
	if err != nil {
		return err
	}
	defer db.Close()
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS persons (
		person_id TEXT NOT NULL PRIMARY KEY,
		last_name TEXT NOT NULL,
		first_name TEXT
	);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) CreatePerson(personID, lastName, firstName string) error {
	db, err := sql.Open("sqlite3", c.CustomDatabase)
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.Exec("INSERT INTO persons (person_id, last_name, first_name) VALUES (?, ?, ?)", personID, lastName, firstName)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) ReadPerson(personID string) (string, string, error) {
	db, err := sql.Open("sqlite3", c.CustomDatabase)
	if err != nil {
		return "", "", err
	}
	defer db.Close()
	var lastName, firstName string
	err = db.QueryRow("SELECT last_name, first_name FROM persons WHERE person_id = ?", personID).Scan(&lastName, &firstName)
	if err != nil {
		return "", "", err
	}
	return lastName, firstName, nil
}

func (c *Client) UpdatePerson(personID, lastName, firstName string) error {
	db, err := sql.Open("sqlite3", c.CustomDatabase)
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.Exec("UPDATE persons SET last_name = ?, first_name = ? WHERE person_id = ?", lastName, firstName, personID)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) DeletePerson(personID string) error {
	db, err := sql.Open("sqlite3", c.CustomDatabase)
	if err != nil {
		return err
	}
	defer db.Close()
	result, err := db.Exec("DELETE FROM persons WHERE person_id = ?", personID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("person not found in the database")
	}
	return nil
}

func (c *Client) CheckPersonExists(personID string) (bool, error) {
	db, err := sql.Open("sqlite3", c.CustomDatabase)
	if err != nil {
		return false, err
	}
	defer db.Close()
	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM persons WHERE person_id = ?)", personID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
