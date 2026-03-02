package services

import (
	"context"

	"github.com/BramAristyo/rest-api-contact-person/internal/domain"
)

type contactService struct {
	repository domain.ContactRepository
}

func NewContactService(repository domain.ContactRepository) domain.ContactService {
	return &contactService{
		repository: repository,
	}
}

func (c contactService) GetAll(ctx context.Context) ([]domain.Contact, error) {
	return c.repository.GetAll(ctx)
}

func (c contactService) Paginate(page int, limit int) ([]domain.Contact, int64, error) {
	//TODO implement me
	panic("implement me")
}

func (c contactService) GetById(id int) (*domain.Contact, error) {
	//TODO implement me
	panic("implement me")
}

func (c contactService) Create(req *domain.CreateContactRequest) (*domain.Contact, error) {
	//TODO implement me
	panic("implement me")
}

func (c contactService) Update(id int, req *domain.UpdateContactRequest) (*domain.Contact, error) {
	//TODO implement me
	panic("implement me")
}

func (c contactService) Delete(id int) error {
	//TODO implement me
	panic("implement me")
}
