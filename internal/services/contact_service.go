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

func (c contactService) Paginate(ctx context.Context, page int, limit int) ([]domain.Contact, int64, error) {
	return c.repository.Paginate(ctx, page, limit)
}

func (c contactService) GetById(ctx context.Context, id int) (*domain.Contact, error) {
	return c.repository.GetById(ctx, id)
}

func (c contactService) Store(ctx context.Context, req *domain.CreateContactRequest) (*domain.Contact, error) {
	return c.repository.Store(ctx, &domain.Contact{
		Name:  req.Name,
		Email: req.Email,
		Phone: req.Phone,
	})
}

func (c contactService) Update(ctx context.Context, id int, req *domain.UpdateContactRequest) (*domain.Contact, error) {
	return c.repository.Update(ctx, id, &domain.Contact{
		Id:    id,
		Name:  req.Name,
		Email: req.Email,
		Phone: req.Phone,
	})
}

func (c contactService) Delete(ctx context.Context, id int) error {
	if err := c.repository.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
