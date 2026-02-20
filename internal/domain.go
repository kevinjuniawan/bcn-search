package internal

import "errors"

type Airline struct {
	Code string `json:"code" validate:"required"`
	Name string `json:"name" validate:"required"`
}

type Airport struct {
	Airport   string `json:"airport" validate:"required"`
	City      string `json:"city" validate:"required"`
	Datetime  string `json:"datetime" validate:"required"`
	Timestamp int64  `json:"timestamp" validate:"required"`
}

type Duration struct {
	TotalMinute int16  `json:"total_minute" validate:"required"`
	Formatted   string `json:"formatted" validate:"required"`
}

type Price struct {
	Amount   int    `json:"amount" validate:"required"`
	Currency string `json:"currency" validate:"required,oneof=IDR USD"`
}

func (p Price) AmountInIDR() int {
	if p.Currency == "IDR" {
		return p.Amount
	}
	return p.Amount * 15000 //Infer Other than IDR is USD & fixed rate 1 USD = 15000 IDR
}

type Bag struct {
	CarryOn string `json:"carry_on" validate:"required"`
	Checked string `json:"checked" validate:"required"`
}

type Class string

const (
	EconomyClass  Class = "economy"
	BusinessClass Class = "business"
)

type SearchResponse struct {
	Metadata Metadata `json:"metadata" validate:"required"`
	Flights  []Flight `json:"flights" validate:"required"`
}

type Metadata struct {
	ProviderCount     int16 `json:"provider_count" validate:"required"`
	SucceededProvider int16 `json:"succeeded_provider" validate:"required"`
	SearchTimeMs      int32 `json:"search_time_ms" validate:"required"`
	IsCache           bool
}

type Flight struct {
	ID             string   `json:"id" validate:"required"`
	Provider       string   `json:"provider" validate:"required"`
	Airline        Airline  `json:"airline" validate:"required"`
	FlightNumber   string   `json:"flight_number" validate:"required"`
	Departure      Airport  `json:"departure" validate:"required"`
	Arrival        Airport  `json:"arrival" validate:"required"`
	Duration       Duration `json:"duration" validate:"required"`
	Stops          int8     `json:"stops" validate:"min=0"`
	Price          Price    `json:"price" validate:"required"`
	AvailableSeats int16    `json:"available_seats" validate:"min=0"`
	CabinClass     Class    `json:"cabin_class" validate:"required,oneof=economy business"`
	Aircraft       *string  `json:"aircraft" validate:"omitempty"`
	Amenities      []string `json:"amenities" validate:"omitempty"`
	Baggage        Bag      `json:"baggage" validate:"required"`
	Layover        int      `json:"layover,omitempty" validate:"min=0"`
}

type GetFlightsParams struct {
	Origin        string              `json:"origin" validate:"required"`
	Destination   string              `json:"destination" validate:"required"`
	DepartureDate string              `json:"departure_date" validate:"required"`
	Passenger     int16               `json:"passenger" validate:"required"`
	ReturnDate    *string             `json:"return_date" validate:"omitempty"`
	SortType      SortType            `json:"sort_type" validate:"required"`
	CabinClass    string              `json:"cabin_class" validate:"omitempty"`
	Filter        *FilterFlightParams `json:"filter"`
}

func (p GetFlightsParams) Validate() error {
	if p.Origin == "" {
		return errors.New("origin must be filled")
	}

	if p.Destination == "" {
		return errors.New("destination must be filled")
	}

	if p.DepartureDate == "" {
		return errors.New("departure date must be filled")
	}

	if p.CabinClass == "" || (p.CabinClass != string(BusinessClass) && p.CabinClass != string(EconomyClass)) {
		return errors.New("cabin class is invalid")
	}

	if p.SortType < 0 || p.SortType > 6 {
		return errors.New("sort type is invalid")
	}

	if p.Filter != nil {
		if p.Filter.TimeRange != nil && (p.Filter.TimeRange.Type == "" || p.Filter.TimeRange.From == "" || p.Filter.TimeRange.To == "") {
			return errors.New("time range type, from, and to is invalid")
		}
		if p.Filter.Price != nil && (p.Filter.Price.LowestPrice < 0 || p.Filter.Price.HighestPrice < 0 || p.Filter.Price.LowestPrice > p.Filter.Price.HighestPrice) {
			return errors.New("lowest price and highest price is invalid")
		}
		if p.Filter.Stops != nil && (*p.Filter.Stops < 0) {
			return errors.New("stops is invalid")
		}
	}
	return nil
}

type SortType int

const (
	SortBestValueType SortType = iota
	SortLowestPriceType
	SortHighestPriceType
	SortShortestDurationType
	SortLongestDurationType
	SortDepartureType
	SortArrivalType
)

type FilterFlightTimeType string

const (
	FilterFlightTimeTypeDeparture FilterFlightTimeType = "departure"
	FilterFlightTimeTypeArrival   FilterFlightTimeType = "arrival"
)

type FilterFlightParams struct {
	Airline   []string                 `json:"airline" validate:"omitempty"`
	Price     *FilterFlightPriceParams `json:"price" validate:"omitempty"`
	Stops     *int8                    `json:"stops" validate:"omitempty"`
	TimeRange *FilterFlightTimeParams  `json:"time_range" validate:"omitempty"`
}

type FilterFlightPriceParams struct {
	LowestPrice  int `json:"lowest_price" validate:"omitempty"`
	HighestPrice int `json:"highest_price" validate:"omitempty"`
}

type FilterFlightTimeParams struct {
	Type FilterFlightTimeType `json:"type" validate:"omitempty"`
	From string               `json:"from" validate:"omitempty"`
	To   string               `json:"to" validate:"omitempty"`
}

type FlightData struct {
	CachedData        bool
	ProviderCount     int16
	SucceededProvider int16
	SearchTimeMs      int32
	Flights           []Flight
}
