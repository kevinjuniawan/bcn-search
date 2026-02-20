package http

import (
	"encoding/json"
	"kevinjuniawan/bookcabin/internal"
	"net/http"

	"github.com/gorilla/mux"
)

type Handler struct {
	flightService internal.InternalService
	cacheService  ICache
}

type Params struct {
	FlightService *internal.InternalService
	CacheService  ICache
}

func NewHandler(p Params) *Handler {
	return &Handler{
		flightService: *p.FlightService,
		cacheService:  p.CacheService,
	}
}

func (h *Handler) InitRouter() http.Handler {
	mux := mux.NewRouter()
	mux.HandleFunc("/flights/search", h.SearchFlights).Methods("POST")
	return mux
}

func (h *Handler) SearchFlights(w http.ResponseWriter, r *http.Request) {
	var params internal.GetFlightsParams
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		WriteJSON(w, 400, NewResponse(err.Error(), internal.SearchResponse{}, params))
		return
	}

	if h.cacheService.IsRequestLimiterExceeded(r.Context(), r.URL.String()) {
		WriteJSON(w, 429, NewResponse("Too many requests", internal.SearchResponse{}, params))
		return
	}

	err = params.Validate()
	if err != nil {
		WriteJSON(w, 400, NewResponse(err.Error(), internal.SearchResponse{}, params))
		return
	}

	flights, err := h.flightService.GetFlights(r.Context(), params)
	if err != nil {
		WriteJSON(w, 500, NewResponse(err.Error(), internal.SearchResponse{}, params))
		return
	}

	WriteJSON(w, 200, NewResponse("Flights retrieved successfully", flights, params))
}
