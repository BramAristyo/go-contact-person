package domain

import (
	"context"
	"time"
)

type Contact struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateContactRequest struct {
	Name  string `json:"name" validate:"required,min=3"`
	Email string `json:"email" validate:"required,email"`
	Phone string `json:"phone" validate:"required,e164"`
}

type UpdateContactRequest struct {
	Name  string `json:"name" validate:"required,min=3"`
	Email string `json:"email" validate:"required,email"`
	Phone string `json:"phone" validate:"required,e164"`
}

type ContactRepository interface {
	GetAll(ctx context.Context) ([]Contact, error)
	Paginate(ctx context.Context, page int, limit int) ([]Contact, int64, error)
	GetById(ctx context.Context, id int) (*Contact, error)
	Create(ctx context.Context, contact *Contact) (*Contact, error)
	Update(ctx context.Context, id int, contact *Contact) (*Contact, error)
	Delete(ctx context.Context, id int) error
}

type ContactService interface {
	GetAll(ctx context.Context) ([]Contact, error)
	Paginate(ctx context.Context, page int, limit int) ([]Contact, int64, error)
	GetById(ctx context.Context, id int) (*Contact, error)
	Create(ctx context.Context, req *CreateContactRequest) (*Contact, error)
	Update(ctx context.Context, id int, req *UpdateContactRequest) (*Contact, error)
	Delete(ctx context.Context, id int) error
}
