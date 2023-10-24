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

	_, err = db.Exec("DROP TABLE persons")
	if err != nil {
		log.Error("clear table", zap.String("error", err.Error()))
		return nil, err
	}
	log.Info("table cleared")

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

func (s *Storage) NewPerson(name, surname, patronymic string, age int, gender, nationality string) (int64, error) {
	var id int64
	err := s.DB.QueryRow(`INSERT INTO persons (name, surname, patronymic, 
		age, gender, nationality)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		name, surname, patronymic, age, gender, nationality).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Storage) DeletePerson(id int) error {
	_, err := s.DB.Exec("DELETE FROM persons WHERE id = $1", id)
	return err
}
