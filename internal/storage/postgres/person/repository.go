package person

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"persons/internal/domain/models"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Save(person models.Person) error {
	const op = "postgres.person.Repository.Save"
	query := `INSERT INTO 
    			persons 
    			(id, email, login, password_hash) 
				VALUES ($1, $2, $3, $4)`

	_, err := r.db.Exec(query, person.Id, person.Email, person.Login, person.PasswordHash)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *Repository) GetById(id uuid.UUID) (*models.Person, error) {
	const op = "postgres.person.Repository.GetById"
	var person models.Person
	query := `SELECT 
    			id, email, login, password_hash, about_me, image 
				FROM
				    persons 
				WHERE id = $1`
	err := r.db.Get(&person, query, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &person, nil
}

func (r *Repository) GetAll() ([]models.Person, error) {
	const op = "postgres.person.Repository.GetAll"
	var persons []models.Person
	query := `SELECT 
				id, email, login, password_hash, about_me, image
				FROM 
				    persons`
	err := r.db.Select(&persons, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return persons, nil
}
