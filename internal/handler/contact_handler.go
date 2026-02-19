package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/BramAristyo/rest-api-contact-person/internal/domain"
	"github.com/BramAristyo/rest-api-contact-person/pkg/response"
	"github.com/go-playground/validator/v10"
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
	response.WriteSuccess(w, contacts, "Contacts retrieved successfully", http.StatusOK)
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

	response.WritePaginated(w, contacts, response.PaginationMeta{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}, http.StatusOK)
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

	response.WriteSuccess(w, c, "Contact retrieved successfully", http.StatusOK)
}

func (h *ContactHandler) Store(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateContactRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		response.WriteValidationErrors(w, response.FormatValidationError(err), http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	var exists bool
	err := h.db.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM contacts WHERE email = $1)`, req.Email).Scan(&exists)
	if err != nil {
		response.WriteError(w, "Failed to check existing contact", http.StatusInternalServerError)
		return
	}

	if exists {
		response.WriteError(w, "Email already exists", http.StatusConflict)
		return
	}

	var newId int
	err = h.db.QueryRow(ctx, `INSERT INTO contacts (name, email, phone) VALUES ($1, $2, $3) RETURNING id`, req.Name, req.Email, req.Phone).Scan(&newId)
	if err != nil {
		response.WriteError(w, "Failed to create contact", http.StatusInternalServerError)
		return
	}

	response.WriteSuccess(w, map[string]int{"id": newId}, "Contact created successfully", http.StatusCreated)
}

func (h *ContactHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var req domain.UpdateContactRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		response.WriteValidationErrors(w, response.FormatValidationError(err), http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	result, err := h.db.Exec(ctx, `UPDATE contacts SET name=$1, email=$2, phone=$3 WHERE id=$4`, req.Name, req.Email, req.Phone, id)
	if err != nil {
		response.WriteError(w, "Failed to update contact", http.StatusInternalServerError)
		return
	}

	if result.RowsAffected() == 0 {
		response.WriteError(w, "Contact not found", http.StatusNotFound)
		return
	}

	response.WriteSuccess(w, nil, "Contact updated successfully", http.StatusOK)
}

func (h *ContactHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	ctx := r.Context()

	result, err := h.db.Exec(ctx, `DELETE FROM contacts WHERE id = $1`, id)
	if err != nil {
		response.WriteError(w, "Failed to delete contact", http.StatusInternalServerError)
		return
	}

	if result.RowsAffected() == 0 {
		response.WriteError(w, "Contact not found", http.StatusNotFound)
		return
	}

	response.WriteSuccess(w, nil, "Contact deleted successfully", http.StatusOK)
}
