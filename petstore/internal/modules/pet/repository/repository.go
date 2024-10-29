package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"test/internal/models"

	"github.com/lib/pq"
)

type PetStorage struct {
	logger *zap.Logger
	DB     *sqlx.DB
}

func NewPetStorage(db *sqlx.DB, logger *zap.Logger) IPetStorage {
	return &PetStorage{
		logger: logger,
		DB:     db}
}

func (ps *PetStorage) Create(pet *models.Pet) error {
	categoryQuery := `INSERT INTO categories (name) VALUES ($1) RETURNING id`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := ps.DB.QueryRowContext(ctx, categoryQuery, pet.Category.Name).Scan(&pet.Category.ID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			ps.logger.Error("no rows in categories table", zap.Error(err))
			return ErrEditConflict
		default:
			ps.logger.Error(" error on inserting category", zap.Error(err))
			return err
		}
	}

	query := `INSERT INTO pets (name, category_id, photo_urls, status)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, status, photo_urls`

	args := []any{
		pet.Name,
		pet.Category.ID,
		pq.Array(pet.PhotoUrls),
		pet.Status,
	}

	ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err = ps.DB.QueryRowContext(ctx, query, args...).Scan(&pet.ID, &pet.Name, &pet.Status, pq.Array(&pet.PhotoUrls))
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			ps.logger.Error("no rows in pets", zap.Error(err))
			return ErrEditConflict
		default:
			ps.logger.Error("error on inserting pet", zap.Error(err))
			return err
		}
	}

	for i := 0; i < len(pet.Tags); i++ {
		tagQuery := `INSERT INTO tags (name) VALUES ($1) RETURNING id`
		err = ps.DB.QueryRowContext(ctx, tagQuery, pet.Tags[i].Name).Scan(&pet.Tags[i].ID)
		if err != nil {
			switch {
			case errors.Is(err, sql.ErrNoRows):
				ps.logger.Error("no rows in tags table", zap.Error(err))
				return ErrEditConflict
			default:
				ps.logger.Error("some error on inserting tag", zap.Error(err))
				return err
			}
		}
		tagPetQuery := `INSERT INTO pet_tags (pet_id, tag_id) VALUES ($1, $2) RETURNING pet_id, tag_id`
		err = ps.DB.QueryRowContext(ctx, tagPetQuery, pet.ID, pet.Tags[i].ID).Scan(&pet.ID, &pet.Tags[i].ID)
		if err != nil {
			switch {
			case errors.Is(err, sql.ErrNoRows):
				ps.logger.Error("no rows in pet-tag column", zap.Error(err))
				return ErrEditConflict
			default:
				ps.logger.Error("error on inserting pet-tag", zap.Error(err))
				return err
			}
		}
	}

	return nil
}

func (ps *PetStorage) Update(pet *models.Pet) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := ps.DB.ExecContext(ctx, `
		UPDATE pets
		SET name = $1, status = $2
		WHERE id = $3
	`, pet.Name, pet.Status, pet.ID)
	if err != nil {
		ps.logger.Error("error on partially updating pet", zap.Error(err))
		return ErrEditConflict
	}
	return nil
}

func (ps *PetStorage) Update_put(pet *models.Pet) error {
	categoryQuery := `UPDATE  categories SET name = $1 WHERE id = $2 RETURNING id`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := ps.DB.QueryRowContext(ctx, categoryQuery, pet.Category.Name, pet.Category.ID).Scan(&pet.Category.ID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			ps.logger.Error("no rows in category table", zap.Error(err))
			return ErrEditConflict
		default:
			ps.logger.Error(" error on updating category", zap.Error(err), zap.Error(ErrEditConflict))
			return ErrEditConflict
		}
	}
	query := `UPDATE pets 
	SET name = $1, photo_urls = $2, status = $3 
	WHERE id = $4 
	RETURNING id, name, status, photo_urls`

	args := []any{
		pet.Name,
		pq.Array(pet.PhotoUrls),
		pet.Status,
		pet.ID,
	}

	ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err = ps.DB.QueryRowContext(ctx, query, args...).Scan(&pet.ID, &pet.Name, &pet.Status, pq.Array(&pet.PhotoUrls))
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			ps.logger.Error("no rows in pets updated", zap.Error(err))
			return ErrEditConflict
		default:
			ps.logger.Error("error on updating pet", zap.Error(err))
			return err
		}
	}

	return nil
}

func (ps *PetStorage) Delete(id int64) error {
	_, err := ps.DB.Exec(`
		DELETE FROM pets WHERE id = $1
	`, id)
	if err != nil {
		ps.logger.Error("error on deleting pet", zap.Error(err))
		return err
	}

	_, err = ps.DB.Exec(`
		DELETE FROM pet_tags WHERE pet_id = $1
	`, id)
	if err != nil {
		ps.logger.Error("error on deleting pet-tags", zap.Error(err))
		return err
	}

	return nil
}

func (ps *PetStorage) GetByID(id int64) (*models.Pet, error) {

	query := `
		SELECT pets.*, categories.name
FROM pets
INNER JOIN categories ON pets.category_id = categories.id
WHERE pets.id = $1
	`

	pet := &models.Pet{
		Category: &models.Category{},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	err := ps.DB.QueryRowContext(ctx, query, id).Scan(
		&pet.ID,
		&pet.Category.ID,
		&pet.Name,
		&pet.Status,
		pq.Array(&pet.PhotoUrls),
		&pet.Category.Name,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			ps.logger.Error("no rows in pets", zap.Error(err))
			return nil, ErrPetNotFound
		default:
			ps.logger.Error("some error on getting pet by id", zap.Error(err))
			return nil, err
		}
	}

	tagQuery := `
	SELECT tags.*
FROM pet_tags
INNER JOIN tags ON pet_tags.tag_id = tags.id
INNER JOIN pets ON pet_tags.pet_id = pets.id
WHERE pet_tags.pet_id = $1
	`
	ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	rows, err := ps.DB.QueryContext(ctx, tagQuery, id)
	if err != nil {
		ps.logger.Error("some error on getting rows from pet tags", zap.Error(err))
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var tag models.Tag
		err := rows.Scan(
			&tag.ID,
			&tag.Name,
		)
		if err != nil {
			return nil, err
		}
		pet.Tags = append(pet.Tags, &tag)
	}

	return pet, nil
}

func (ps *PetStorage) GetByStatus(status string) ([]models.Pet, error) {
	var result []models.Pet
	query := `
	SELECT pets.*, categories.name
FROM pets
INNER JOIN categories ON pets.category_id = categories.id
WHERE pets.status = $1
`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	rows, err := ps.DB.QueryContext(ctx, query, status)
	if err != nil {
		ps.logger.Error("some error on getting rows from pets", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		pet := models.Pet{
			Category: &models.Category{},
		}

		err := rows.Scan(
			&pet.ID,
			&pet.Category.ID,
			&pet.Name,
			&pet.Status,
			pq.Array(&pet.PhotoUrls),
			&pet.Category.Name,
		)
		if err != nil {
			ps.logger.Error("some error on getting rows from pets", zap.Error(err))
			return nil, err
		}

		tagQuery := `
	SELECT tags.*
FROM pet_tags
INNER JOIN tags ON pet_tags.tag_id = tags.id
INNER JOIN pets ON pet_tags.pet_id = pets.id
WHERE pet_tags.pet_id = $1
	`
		ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)

		defer cancel()

		rows, err := ps.DB.QueryContext(ctx, tagQuery, pet.ID)
		if err != nil {
			ps.logger.Error("some error on getting rows from pet tags", zap.Error(err))
			return nil, err
		}

		defer rows.Close()

		for rows.Next() {
			var tag models.Tag
			err := rows.Scan(
				&tag.ID,
				&tag.Name,
			)
			if err != nil {
				return nil, err
			}
			pet.Tags = append(pet.Tags, &tag)
		}

		result = append(result, pet)

	}

	return result, nil
}
