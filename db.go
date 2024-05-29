package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
)

type MyPGSQLStorage struct {
	db *sql.DB
}

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "postgres"
)

func NewMyPGSQLStorage(cfg mysql.Config) *MyPGSQLStorage {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MySQL ")

	return &MyPGSQLStorage{db: db}

}

func (s *MyPGSQLStorage) Init() (*sql.DB, error) {
	//initialize the tables
	if err := s.createProjectsTable(); err != nil {
		return nil, err
	}
	if err := s.createUsersTable(); err != nil {
		return nil, err
	}
	if err := s.createTasksTable(); err != nil {
		return nil, err
	}

	return s.db, nil

}

func (s *MyPGSQLStorage) createProjectsTable() error {
	_, err := s.db.Exec(`CREATE TABLE IF NOT EXISTS projects (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);
	`)
	return err
}

func (s *MyPGSQLStorage) createTasksTable() error {
	_, err := s.db.Exec(`
	CREATE TABLE IF NOT EXISTS tasks (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		status VARCHAR(20) NOT NULL DEFAULT 'TODO',
		project_id INT NOT NULL,
		assigned_to_id INT NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (assigned_to_id) REFERENCES users(id),
		FOREIGN KEY (project_id) REFERENCES projects(id)
	);
		
	`)
	return err
}

func (s *MyPGSQLStorage) createUsersTable() error {
	_, err := s.db.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		email VARCHAR(255) NOT NULL UNIQUE,
		first_name VARCHAR(255) NOT NULL,
		last_name VARCHAR(255) NOT NULL,
		password VARCHAR(255) NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);
		
	`)

	return err
}
