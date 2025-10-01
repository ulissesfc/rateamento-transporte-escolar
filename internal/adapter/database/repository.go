package database

import (
	"context"
	"database/sql"

	"github.com/ulissesfc/rateamento-transporte-escolar.git/internal/db"
	"github.com/ulissesfc/rateamento-transporte-escolar.git/internal/domain"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetStudents(limit int) ([]domain.Student, error) {
	queries := db.New(r.db)
	ctx := context.Background()

	users, err := queries.GetUsers(ctx, int32(limit))
	if err != nil {
		return nil, err
	}

	var students []domain.Student
	for _, user := range users {
		students = append(students, domain.Student{
			Id:          int(user.ID),
			Name:        user.Name,
			Institution: user.Institution,
			Address:     user.Address,
			Longitude:   user.Longitude,
			Latitude:    user.Latitude,
		})
	}

	return students, nil
}

func (r *Repository) InsertRoute(param db.InsertRouteParams) (int, error) {
	queries := db.New(r.db)
	ctx := context.Background()

	response, err := queries.InsertRoute(ctx, param)

	return int(response), err
}

func (r *Repository) InsertNode(param db.InsertNodeParams) error {
	queries := db.New(r.db)
	ctx := context.Background()

	_, err := queries.InsertNode(ctx, param)

	return err
}
