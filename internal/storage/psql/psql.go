package psql

import (
	"database/sql"
	"fmt"

	"github.com/Raitfolt/juntest/config"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type Storage struct {
	DB *sql.DB
}

func New(log *zap.Logger, cfg *config.Config) (*Storage, error) {
	log.Info("create new storage connect")
	psqlInfo := fmt.Sprintf("host=%s port =%d user=%s password=%s sslmode=disable",
		cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresUser, cfg.PostgresPassword)
	log.Info("postgresql", zap.String("connection string", psqlInfo))

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Error("sql.Open", zap.String("error", err.Error()))
		return nil, err
	}

	_, err = db.Exec("DROP DATABASE " + cfg.PostgresDB)
	if err != nil {
		log.Error("clear database", zap.String("error", err.Error()))
		return nil, err
	}
	log.Info("database cleared")

	_, err = db.Exec("CREATE DATABASE " + cfg.PostgresDB)
	if err != nil {
		log.Error("create database", zap.String("error", err.Error()))
		return nil, err
	}
	log.Info("database created")

	_, err = db.Exec(`CREATE TABLE persons(
		id SERIAL PRIMARY KEY,
		name VARCHAR(255),
		surname VARCHAR(255),
		patronymic VARCHAR(255),
		age INTEGER,
		gender VARCHAR(10),
		nationality VARCHAR(255)
	);`)
	if err != nil {
		log.Error("create table", zap.String("error", err.Error()))
		return nil, err
	}
	log.Info("table persons created")

	return &Storage{DB: db}, nil
}

func (s *Storage) NewPerson(name, surname, patronymic string) (int64, error) {
	var id int64
	err := s.DB.QueryRow("INSERT INTO persons (name, surname, patronymic) VALUES ($1, $2, $3) RETURNING id",
		name, surname, patronymic).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Storage) UpdateGender(id int, gender string) error {
	_, err := s.DB.Exec("UPDATE persons SET gender = $1 WHERE id = $2", gender, id)
	return err
}

func (s *Storage) UpdateAge(id int, age int) error {
	_, err := s.DB.Exec("UPDATE persons SET age = $1 WHERE id = $2", age, id)
	return err
}

func (s *Storage) UpdateNationality(id int, nationality string) error {
	_, err := s.DB.Exec("UPDATE persons SET nationality = $1 WHERE id = $2", nationality, id)
	return err
}

func (s *Storage) DeletePerson(id int) error {
	_, err := s.DB.Exec("DELETE FROM persons WHERE id = $1", id)
	return err
}
