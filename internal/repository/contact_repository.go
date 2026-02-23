package repository

import (
	"github.com/BramAristyo/rest-api-contact-person/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type contactRepository struct {
	db *pgxpool.Pool
}

func (c contactRepository) GetAll() ([]domain.Contact, error) {
	//TODO implement me
	panic("implement me")
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
