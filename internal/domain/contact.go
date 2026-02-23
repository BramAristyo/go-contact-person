package domain

import "time"

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
	GetAll() ([]Contact, error)
	Paginate(page int, limit int) ([]Contact, int64, error)
	GetById(id int) (*Contact, error)
	Create(contact *Contact) (*Contact, error)
	Update(id int, contact *Contact) (*Contact, error)
	Delete(id int) error
}

type ContactService interface {
	GetAll() ([]Contact, error)
	Paginate(page int, limit int) ([]Contact, int64, error)
	GetById(id int) (*Contact, error)
	Create(req *CreateContactRequest) (*Contact, error)
	Update(id int, req *UpdateContactRequest) (*Contact, error)
	Delete(id int) error
}
