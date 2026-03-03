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

func (c contactRepository) Store(ctx context.Context, contact *domain.Contact) (*domain.Contact, error) {
	var exists bool
	// use QueryRow to check if email already exists, since we only expect one row (true/false), and it returns a Row object that we can scan directly.
	err := c.db.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM contacts WHERE email = $1)`, contact.Email).Scan(&exists)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, errors.New("email already exists")
	}

	var newId int
	err = c.db.QueryRow(ctx, `INSERT INTO contacts (name, email, phone) VALUES ($1, $2, $3) RETURNING id`, contact.Name, contact.Email, contact.Phone).Scan(&newId)
	if err != nil {
		return nil, err
	}

	return c.GetById(ctx, newId)
}

func (c contactRepository) Update(ctx context.Context, id int, contact *domain.Contact) (*domain.Contact, error) {
	// Start a transaction to ensure data integrity during the update process.
	tx, err := c.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	// using Exec instead of QueryRow since we don't need to return any data, just check affected rows.
	result, err := tx.Exec(ctx, `UPDATE contacts SET name=$1, email=$2, phone=$3 WHERE id=$4`, contact.Name, contact.Email, contact.Phone, id)
	if err != nil {
		return nil, err
	}

	if result.RowsAffected() == 0 {
		return nil, errors.New("Contact not found")
	}

	// Commit the transaction after successful update.
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return c.GetById(ctx, id)
}

func (c contactRepository) Delete(ctx context.Context, id int) error {
	result, err := c.db.Exec(ctx, `DELETE FROM contacts WHERE id = $1`, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return err
	}

	return nil
}

func NewContactRepository(db *pgxpool.Pool) domain.ContactRepository {
	return &contactRepository{
		db: db,
	}
}
