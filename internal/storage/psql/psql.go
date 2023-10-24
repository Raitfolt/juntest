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

type Person struct {
	ID          string `json:"id" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Surname     string `json:"surname" validate:"required"`
	Patronymic  string `json:"patronymic,omitempty"`
	Age         int    `json:"age"`
	Gender      string `json:"gender"`
	Nationality string `json:"nationality"`
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

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS persons(
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
	log.Info("table persons connected")

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

func (s *Storage) ChangePerson(id int64, name, surname, patronymic string, age int, gender, nationality string) (int64, error) {
	var rid int64
	err := s.DB.QueryRow(`UPDATE persons 
				SET name = $1, surname = $2, patronymic = $3, 
				age = $4, gender = $5, nationality = $6
				WHERE id = $7 RETURNING id`,
		name, surname, patronymic, age, gender, nationality, id).Scan(&rid)
	if err != nil {
		return 0, err
	}

	return rid, nil
}

func (s *Storage) DeletePerson(id int64) error {
	_, err := s.DB.Exec("DELETE FROM persons WHERE id = $1", id)
	return err
}

func (s *Storage) GetPerson(id int64) ([]Person, error) {
	rows, err := s.DB.Query("SELECT * FROM persons")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var persons []Person

	for rows.Next() {
		var person Person
		if err := rows.Scan(
			&person.ID,
			&person.Name,
			&person.Surname,
			&person.Patronymic,
			&person.Age,
			&person.Gender,
			&person.Nationality); err != nil {
			return nil, err
		}
		persons = append(persons, person)
	}
	return persons, nil
}

func (s *Storage) ListPersons() ([]Person, error) {
	rows, err := s.DB.Query("SELECT * FROM persons")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var persons []Person

	for rows.Next() {
		var person Person
		if err := rows.Scan(
			&person.ID,
			&person.Name,
			&person.Surname,
			&person.Patronymic,
			&person.Age,
			&person.Gender,
			&person.Nationality); err != nil {
			return nil, err
		}
		persons = append(persons, person)
	}
	return persons, nil
}
