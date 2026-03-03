package repository

import (
	"context"
	"errors"

	"github.com/BramAristyo/rest-api-contact-person/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type contactRepository struct {
	db *pgxpool.Pool
}

func (c contactRepository) GetAll(ctx context.Context) ([]domain.Contact, error) {
	rows, err := c.db.Query(ctx, `SELECT id, name, email, phone, created_at, updated_at FROM contacts`)
	if err != nil {
		return nil, err
	}

	// Rows is stream of data from database, we need to close it after we're done to free up resources.
	defer rows.Close()

	var contacts []domain.Contact

	// Iterate over the rows and scan data into Structs.
	// Memory efficient for large datasets since it doesn't load everything into memory at once.
	for rows.Next() {
		var c domain.Contact
		err := rows.Scan(
			&c.Id,
			&c.Name,
			&c.Email,
			&c.Phone,
			&c.CreatedAt,
			&c.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		contacts = append(contacts, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return contacts, nil
}

func (c contactRepository) Paginate(ctx context.Context, page int, limit int) ([]domain.Contact, int64, error) {
	offset := (page - 1) * limit

	// use Query instead of QueryRow since we expect multiple rows, and it returns a Rows object that we can iterate over.
	rows, err := c.db.Query(ctx, `SELECT id, name, email, phone, created_at, updated_at FROM contacts ORDER BY id LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()

	var contacts []domain.Contact
	for rows.Next() {
		var c domain.Contact
		err := rows.Scan(
			&c.Id,
			&c.Name,
			&c.Email,
			&c.Phone,
			&c.CreatedAt,
			&c.UpdatedAt,
		)

		if err != nil {
			return nil, 0, err
		}

		contacts = append(contacts, c)
	}

	if rows.Err() != nil {
		return nil, 0, err
	}

	var total int64
	err = c.db.QueryRow(ctx, `SELECT COUNT(*) FROM contacts`).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return contacts, total, nil
}

func (c contactRepository) GetById(ctx context.Context, id int) (*domain.Contact, error) {
	var contact domain.Contact

	err := c.db.QueryRow(ctx, `SELECT id, name, email, phone, created_at, updated_at FROM contacts WHERE id = $1`, id).Scan(
		&contact.Id,
		&contact.Name,
		&contact.Email,
		&contact.Phone,
		&contact.CreatedAt,
		&contact.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("not found")
		}

		return nil, err
	}

	return &contact, nil
}

func (c contactRepository) Create(ctx context.Context, contact *domain.Contact) (*domain.Contact, error) {
	//TODO implement me
	panic("implement me")
}

func (c contactRepository) Update(ctx context.Context, id int, contact *domain.Contact) (*domain.Contact, error) {
	//TODO implement me
	panic("implement me")
}

func (c contactRepository) Delete(ctx context.Context, id int) error {
	//TODO implement me
	panic("implement me")
}

func NewContactRepository(db *pgxpool.Pool) domain.ContactRepository {
	return &contactRepository{
		db: db,
	}
}
