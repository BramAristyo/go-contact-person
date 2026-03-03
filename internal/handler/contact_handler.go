package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/BramAristyo/rest-api-contact-person/internal/domain"
	"github.com/BramAristyo/rest-api-contact-person/pkg/response"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ContactHandler struct {
	db       *pgxpool.Pool
	validate *validator.Validate
	service  domain.ContactService
}

func NewContactHandler(db *pgxpool.Pool, validate *validator.Validate, service domain.ContactService) *ContactHandler {
	return &ContactHandler{
		db:       db,
		validate: validate,
		service:  service,
	}
}

func (h *ContactHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	contacts, err := h.service.GetAll(ctx)

	if err != nil {
		response.WriteError(w, "Error iterating contacts", http.StatusInternalServerError)
		return
	}

	response.WriteSuccess(w, contacts, "Contacts retrieved successfully", http.StatusOK)
}

func (h *ContactHandler) Paginate(w http.ResponseWriter, r *http.Request) {
	// Get params and make it int, if error or less than 1, set default value.
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit < 1 {
		limit = 10
	}

	ctx := r.Context()

	contacts, total, err := h.service.Paginate(ctx, page, limit)

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

	contact, err := h.service.GetById(ctx, id)
	if err != nil {
		response.WriteError(w, "Error get contact", http.StatusInternalServerError)
		return
	}

	response.WriteSuccess(w, contact, "Contact retrieved successfully", http.StatusOK)
}

func (h *ContactHandler) Store(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateContactRequest

	// Stream request body and decode into struct, more efficient than read.All and then unmarshal.
	// data, _ := io.ReadAll(r.Body) is not recommended for large payloads as it loads everything into memory at once, while Decoder can handle it in chunks.
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		response.WriteValidationErrors(w, response.FormatValidationError(err), http.StatusBadRequest)
		return
	}

	contact, err := h.service.Store(r.Context(), &req)
	if err != nil {
		response.WriteError(w, "Error while create contact", http.StatusInternalServerError)
	}

	response.WriteSuccess(w, map[string]int{"id": contact.Id}, "Contact created successfully", http.StatusCreated)
}

func (h *ContactHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var req domain.UpdateContactRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		response.WriteValidationErrors(w, response.FormatValidationError(err), http.StatusBadRequest)
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		response.WriteError(w, "Invalid contact ID", http.StatusBadRequest)
		return
	}

	contact, err := h.service.Update(r.Context(), idInt, &req)

	response.WriteSuccess(w, contact, "Contact updated successfully", http.StatusOK)
}

func (h *ContactHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		response.WriteError(w, "Error while parse string to Int", http.StatusBadRequest)
		return
	}

	err = h.service.Delete(r.Context(), idInt)

	if err != nil {
		response.WriteError(w, "Error while delete contact", http.StatusInternalServerError)
		return
	}

	response.WriteSuccess(w, nil, "Contact deleted successfully", http.StatusOK)
}
