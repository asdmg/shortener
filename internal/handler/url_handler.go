package handler

import (
	"encoding/json"
	"net/http"

	"shortener/internal/service"
)

type URLHandler struct {
	service *service.URLService
}

func NewURLHandler(service *service.URLService) *URLHandler {
	return &URLHandler{
		service: service,
	}
}

type createURLRequest struct {
	URL string `json:"url"`
}

type createURLResponse struct {
	Code     string `json:"code"`
	ShortURL string `json:"short_url"`
}

func (h *URLHandler) Create(w http.ResponseWriter, r *http.Request) {

	var request createURLRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(
			w,
			"invalid request",
			http.StatusBadRequest,
		)

		return
	}

	url, err := h.service.Create(
		r.Context(),
		request.URL,
	)

	if err != nil {
		http.Error(
			w,
			"could not create short URL",
			http.StatusInternalServerError,
		)

		return
	}

	response := createURLResponse{
		Code: url.Code,
		ShortURL: "http://localhost:8080/" +
			url.Code,
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(response)
}

func (h *URLHandler) Redirect(
	w http.ResponseWriter,
	r *http.Request,
) {

	code := r.PathValue("code")

	url, err := h.service.FindByCode(
		r.Context(),
		code,
	)

	if err != nil {
		http.NotFound(w, r)
		return
	}

	http.Redirect(
		w,
		r,
		url.OriginalURL,
		http.StatusFound,
	)
}
