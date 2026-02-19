package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/BramAristyo/rest-api-contact-person/internal/domain"
	"github.com/BramAristyo/rest-api-contact-person/pkg/response"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ContactHandler struct {
	db *pgxpool.Pool
}

func NewContactHandler(db *pgxpool.Pool) *ContactHandler {
	return &ContactHandler{
		db: db,
	}
}

func (h *ContactHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	rows, err := h.db.Query(ctx, `SELECT id, name, email, phone, created_at, updated_at FROM contacts`)
	if err != nil {
		response.WriteError(w, "Failed to fetch contacts", http.StatusInternalServerError)
		return
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
			response.WriteError(w, "Failed to parse contact data", http.StatusInternalServerError)
			return
		}

		contacts = append(contacts, c)
	}

	if rows.Err() != nil {
		response.WriteError(w, "Error iterating contacts", http.StatusInternalServerError)
		return
	}
	response.WriteSucess(w, contacts, "Contacts retrieved successfully", http.StatusOK)
}

func (h *ContactHandler) Paginate(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	ctx := r.Context()
	rows, err := h.db.Query(ctx, `SELECT id, name, email, phone, created_at, updated_at FROM contacts ORDER BY id LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		response.WriteError(w, "Failed to fetch contacts", http.StatusInternalServerError)
		return
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
			response.WriteError(w, "Failed to parse contact data", http.StatusInternalServerError)
			return
		}

		contacts = append(contacts, c)
	}

	if rows.Err() != nil {
		response.WriteError(w, "Error iterating contacts", http.StatusInternalServerError)
		return
	}

	var total int64
	err = h.db.QueryRow(ctx, `SELECT COUNT(*) FROM contacts`).Scan(&total)
	if err != nil {
		response.WriteError(w, "Failed to count contacts", http.StatusInternalServerError)
		return
	}

	totalPages := (total + int64(limit) - 1) / int64(limit)

	response.WriteSucess(w, response.PaginatedResponse{
		Data: contacts,
		Meta: response.PaginationMeta{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}, "Contacts retrieved successfully", http.StatusOK)
}

func (h *ContactHandler) GetById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		response.WriteError(w, "Invalid contact ID", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	var c domain.Contact

	err = h.db.QueryRow(ctx, `SELECT id, name, email, phone, created_at, updated_at FROM contacts WHERE id = $1`, id).Scan(
		&c.Id,
		&c.Name,
		&c.Email,
		&c.Phone,
		&c.CreatedAt,
		&c.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.WriteError(w, "Contact not found", http.StatusNotFound)
			return
		}

		response.WriteError(w, "Failed to fetch contact", http.StatusInternalServerError)
		return
	}

	response.WriteSucess(w, c, "Contact retrieved successfully", http.StatusOK)
}
