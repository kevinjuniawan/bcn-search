package http

import (
	"encoding/json"
	"kevinjuniawan/bookcabin/internal"
	"net/http"
)

type Response struct {
	SearchCriteria SearchCriteriaResponse `json:"search_criteria"`
	Metadata       MetadataResponse       `json:"metadata"`
	Message        string                 `json:"message"`
	Flights        []internal.Flight      `json:"flights"`
}

type SearchCriteriaResponse struct {
	Origin        string                       `json:"origin"`
	Destination   string                       `json:"destination"`
	DepartureDate string                       `json:"departure_date"`
	Passenger     int16                        `json:"passengers"`
	CabinClass    string                       `json:"cabin_class"`
	SortType      internal.SortType            `json:"sort_type"`
	Filter        *internal.FilterFlightParams `json:"filter,omitempty"`
}

type MetadataResponse struct {
	TotalResults       int32 `json:"total_results"`
	ProvidersQueried   int16 `json:"providers_queried"`
	ProvidersSucceeded int16 `json:"providers_succeeded"`
	ProvidersFailed    int16 `json:"providers_failed"`
	SearchTimeMs       int32 `json:"search_time_ms"`
	CacheHit           bool  `json:"cache_hit"`
}

func NewResponse(message string, data internal.SearchResponse, params internal.GetFlightsParams) Response {
	return Response{
		SearchCriteria: SearchCriteriaResponse{
			Origin:        params.Origin,
			Destination:   params.Destination,
			DepartureDate: params.DepartureDate,
			CabinClass:    params.CabinClass,
			Passenger:     params.Passenger,
			SortType:      params.SortType,
			Filter:        params.Filter,
		},
		Metadata: MetadataResponse{
			TotalResults:       int32(len(data.Flights)),
			ProvidersQueried:   data.Metadata.ProviderCount,
			ProvidersSucceeded: data.Metadata.SucceededProvider,
			ProvidersFailed:    data.Metadata.ProviderCount - data.Metadata.SucceededProvider,
			SearchTimeMs:       data.Metadata.SearchTimeMs,
			CacheHit:           data.Metadata.IsCache,
		},
		Message: message,
		Flights: data.Flights,
	}

}

func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
