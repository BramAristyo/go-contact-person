package repository

import (
	"context"
	"net/http"

	"github.com/BramAristyo/rest-api-contact-person/internal/domain"
	"github.com/BramAristyo/rest-api-contact-person/pkg/response"
	"github.com/jackc/pgx/v5/pgxpool"
)

type contactRepository struct {
	db *pgxpool.Pool
}

func (c contactRepository) GetAll(ctx context.Context) ([]domain.Contact, error) {
	rows, err := c.db.Query(ctx, `SELECT id, name, email, phone, created_at, updated_at FROM contacts`)
	if err != nil {

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
			response.WriteError(w, "Failed to parse contact data", http.StatusInternalServerError)
			return
		}

		contacts = append(contacts, c)
	}

	if rows.Err() != nil {
		response.WriteError(w, "Error iterating contacts", http.StatusInternalServerError)
		return
	}
}

func (c contactRepository) Paginate(page int, limit int) ([]domain.Contact, int64, error) {
	//TODO implement me
	panic("implement me")
}

func (c contactRepository) GetById(id int) (*domain.Contact, error) {
	//TODO implement me
	panic("implement me")
}

func (c contactRepository) Create(contact *domain.Contact) (*domain.Contact, error) {
	//TODO implement me
	panic("implement me")
}

func (c contactRepository) Update(id int, contact *domain.Contact) (*domain.Contact, error) {
	//TODO implement me
	panic("implement me")
}

func (c contactRepository) Delete(id int) error {
	//TODO implement me
	panic("implement me")
}

func NewContactRepository(db *pgxpool.Pool) domain.ContactRepository {
	return &contactRepository{
		db: db,
	}
}
