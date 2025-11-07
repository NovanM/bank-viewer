// bank-statement-viewer/handler/http/transaction_handler.go
package http

import (
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/novanm/bank-viewer/backend/domain"
)

const maxUploadSize = 20 * 1024 * 1024 // 20 MB

type TransactionHandler struct {
	service domain.TransactionService
}

func NewTransactionHandler(s domain.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		service: s,
	}
}

func (h *TransactionHandler) RegisterRoutes(mux *http.ServeMux) {

	uploadHandler := http.HandlerFunc(h.Upload)
	mux.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
		uploadHandler.ServeHTTP(w, r)
	})

	mux.HandleFunc("/balance", h.GetBalance)
	mux.HandleFunc("/issues", h.GetIssues)
}

func (h *TransactionHandler) Upload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	reader, err := r.MultipartReader()

	if err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			RespondWithError(w, http.StatusRequestEntityTooLarge, "File exceeds 20MB limit")
		} else {
			RespondWithError(w, http.StatusBadRequest, "Invalid multipart request")
		}
		return
	}

	var filePart io.Reader

	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, "Failed to read multipart part")
			return
		}
		if part.FormName() == "file" {
			filePart = part
			break
		} else {
			io.Copy(io.Discard, part)
		}
	}

	if filePart == nil {
		RespondWithError(w, http.StatusBadRequest, "No 'file' part found in request")
		return
	}

	ctx := r.Context()
	if err := h.service.ProcessUpload(ctx, filePart); err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	RespondWithJSON(w, http.StatusOK, "File uploaded successfully", nil)
}

func (h *TransactionHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	ctx := r.Context()

	balance, err := h.service.GetBalance(ctx)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if balance == nil {
		RespondWithError(w, http.StatusNotFound, "Balance not found")
		return
	}

	RespondWithJSON(w, http.StatusOK, "Balance retrieved successfully", balance)
}

func (h *TransactionHandler) GetIssues(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	q := r.URL.Query()

	page, _ := strconv.Atoi(q.Get("page"))
	if page <= 0 {
		page = 1
	}

	limit, _ := strconv.Atoi(q.Get("limit"))
	if limit <= 0 {
		limit = 10
	}

	if limit > 100 {
		limit = 100
	}

	sortBy := strings.ToLower(q.Get("sort_by"))
	if sortBy == "" {
		sortBy = "timestamp"
	}

	sortDir := strings.ToLower(q.Get("sort_dir"))
	if sortDir != "asc" {
		sortDir = "desc"
	}

	params := domain.PaginationParams{
		Page:    page,
		Limit:   limit,
		SortBy:  sortBy,
		SortDir: sortDir,
	}

	ctx := r.Context()

	issues, err := h.service.GetIssues(ctx, params)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondWithJSON(w, http.StatusOK, "Issues retrieved successfully", issues)
}
