package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gabrielleite03/kenjix_core/internal/dto"
	"github.com/gabrielleite03/kenjix_core/internal/helpers"
	"github.com/gabrielleite03/kenjix_core/internal/service"
)

type ExpenseHandler struct {
	service service.ExpenseService
}

func NewExpenseHandler(service service.ExpenseService) *ExpenseHandler {
	return &ExpenseHandler{service: service}
}

func (h *ExpenseHandler) FindAll(w http.ResponseWriter, r *http.Request) {

	data, err := h.service.FindAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, data)
}

func (h *ExpenseHandler) FindByID(w http.ResponseWriter, r *http.Request) {
	id := getID(r)

	data, err := h.service.FindByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, data)
}

func (h *ExpenseHandler) Create(w http.ResponseWriter, r *http.Request) {

	var req dto.ExpenseCreateUpdateDTO
	r.ParseMultipartForm(10 << 20)

	req.Description = r.FormValue("description")
	req.Status = dto.ExpenseStatus(r.FormValue("status"))

	categoryID, _ := strconv.ParseInt(r.FormValue("category_id"), 10, 64)
	req.CategoryID = categoryID

	amount, _ := strconv.ParseFloat(r.FormValue("amount"), 64)
	req.Amount = helpers.DecimalFromFloat(amount)

	req.Date = helpers.ParseDate(r.FormValue("date"))

	if r.MultipartForm != nil {
		req.Files = r.MultipartForm.File["files"]
	}

	err := h.service.Create(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *ExpenseHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := getID(r)

	var req dto.ExpenseCreateUpdateDTO
	r.ParseMultipartForm(10 << 20)

	req.Description = r.FormValue("description")
	req.Status = dto.ExpenseStatus(r.FormValue("status"))

	categoryID, _ := strconv.ParseInt(r.FormValue("category_id"), 10, 64)
	req.CategoryID = categoryID

	amount, _ := strconv.ParseFloat(r.FormValue("amount"), 64)
	req.Amount = helpers.DecimalFromFloat(amount)

	req.Date = helpers.ParseDate(r.FormValue("date"))

	if r.MultipartForm != nil {
		req.Files = r.MultipartForm.File["files"]
	}

	data, err := h.service.Update(id, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, data)
}

func (h *ExpenseHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := getID(r)

	err := h.service.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func getID(r *http.Request) int64 {
	idStr := r.URL.Path[strings.LastIndex(r.URL.Path, "/")+1:]
	id, _ := strconv.ParseInt(idStr, 10, 64)
	return id
}

func (h *ExpenseHandler) FindAllExpenseCategories(w http.ResponseWriter, r *http.Request) {

	data, err := h.service.FindAllExpenseCategories()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, data)
}
