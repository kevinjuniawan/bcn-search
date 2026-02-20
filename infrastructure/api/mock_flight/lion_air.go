package mockflight

import (
	"encoding/json"
	"kevinjuniawan/bookcabin/internal"
	"kevinjuniawan/bookcabin/pkg/helper"
	"time"
)

type JTAirport struct {
	Code string `json:"code"`
	Name string `json:"name"`
	City string `json:"city"`
}

type JTAirline struct {
	Name string `json:"name"`
	Iata string `json:"iata"`
}

type JTRoute struct {
	From JTAirport `json:"from"`
	To   JTAirport `json:"to"`
}

type JTSchedule struct {
	Departure         string `json:"departure"`
	DepartureTimezone string `json:"departure_timezone"`
	Arrival           string `json:"arrival"`
	ArrivalTimezone   string `json:"arrival_timezone"`
}

type JTLayover struct {
	Airport         string `json:"airport"`
	DurationMinutes int    `json:"duration_minutes"`
}

type JTPrice struct {
	Total    int    `json:"total"`
	Currency string `json:"currency"`
	FareType string `json:"fare_type"`
}

type JTBag struct {
	Cabin string `json:"cabin"`
	Hold  string `json:"hold"`
}

type JTService struct {
	WifiAvailable    bool  `json:"wifi_available"`
	MealsIncluded    bool  `json:"meals_included"`
	BaggageAllowance JTBag `json:"baggage_allowance"`
}

type LionFlight struct {
	ID         string      `json:"id"`
	Carrier    JTAirline   `json:"carrier"`
	Route      JTRoute     `json:"route"`
	Schedule   JTSchedule  `json:"schedule"`
	FlightTime int16       `json:"flight_time"`
	IsDirect   bool        `json:"is_direct"`
	StopCount  int8        `json:"stop_count"`
	Layovers   []JTLayover `json:"layovers,omitempty"`
	Pricing    JTPrice     `json:"pricing"`
	SeatsLeft  int16       `json:"seats_left"`
	PlaneType  string      `json:"plane_type"`
	Services   JTService   `json:"services"`
}

func (l LionFlight) isValid() bool {
	if l.ID == "" || l.Route.From.Code == "" || l.Route.To.Code == "" {
		return false
	}

	if l.Pricing.Total == 0 {
		return false
	}

	return true
}

type LionAirResponse struct {
	Success bool `json:"success"`
	Data    struct {
		AvailableFlights []LionFlight `json:"available_flights"`
	} `json:"data"`
}

func (g LionAirResponse) Normalize() []internal.Flight {
	flights := []internal.Flight{}
	for _, flight := range g.Data.AvailableFlights {

		if !flight.isValid() {
			continue
		}

		departureCity, arrivalCity, isValid := GetCityNameFromAirportCode(flight.Route.From.Code, flight.Route.To.Code)
		if !isValid {
			continue
		}

		locDeparture, err := time.LoadLocation(flight.Schedule.DepartureTimezone)
		if err != nil {
			continue
		}
		departureTime, err := time.ParseInLocation("2006-01-02T15:04:05", flight.Schedule.Departure, locDeparture)
		if err != nil {
			continue
		}

		locArrival, err := time.LoadLocation(flight.Schedule.ArrivalTimezone)
		if err != nil {
			continue
		}
		arrivalTime, err := time.ParseInLocation("2006-01-02T15:04:05", flight.Schedule.Arrival, locArrival)
		if err != nil {
			continue
		}

		duration := FlightDuration{
			Duration: arrivalTime.Sub(departureTime),
		}

		class := internal.BusinessClass
		if flight.Pricing.FareType == "ECONOMY" {
			class = internal.EconomyClass
		}

		layover := 0
		for _, transit := range flight.Layovers {
			layover += int(transit.DurationMinutes)
		}

		amenities := []string{}
		if flight.Services.WifiAvailable {
			amenities = append(amenities, "wifi")
		}
		if flight.Services.MealsIncluded {
			amenities = append(amenities, "meal")
		}

		flights = append(flights, internal.Flight{
			ID:             flight.ID + "_" + mapAirlineCodeToName[LionAirAirline],
			Provider:       mapAirlineCodeToName[LionAirAirline],
			Airline:        internal.Airline{Code: mapAirlineCodeToName[LionAirAirline], Name: string(LionAirAirline)},
			FlightNumber:   flight.ID,
			Departure:      internal.Airport{Airport: flight.Route.From.Code, City: departureCity, Datetime: departureTime.Format(time.RFC3339), Timestamp: departureTime.Unix()},
			Arrival:        internal.Airport{Airport: flight.Route.To.Code, City: arrivalCity, Datetime: arrivalTime.Format(time.RFC3339), Timestamp: arrivalTime.Unix()},
			Duration:       internal.Duration{TotalMinute: int16(duration.Duration.Minutes()), Formatted: duration.Format()},
			Stops:          flight.StopCount,
			Price:          internal.Price{Amount: flight.Pricing.Total, Currency: flight.Pricing.Currency},
			AvailableSeats: flight.SeatsLeft,
			CabinClass:     class,
			Aircraft:       &flight.PlaneType,
			Amenities:      amenities,
			Baggage:        internal.Bag{CarryOn: flight.Services.BaggageAllowance.Cabin, Checked: flight.Services.BaggageAllowance.Hold},
			Layover:        layover,
		})

	}
	return flights
}

type LionAir struct{}

func NewLionAir() *LionAir {
	return &LionAir{}
}

func (b *LionAir) GetFlights(origin, destination, departureDate string) (LionAirResponse, error) {
	helper.RandomDelay(200, 400)
	err := helper.RandomError(10)
	if err != nil {
		return LionAirResponse{}, err
	}
	if origin != "CGK" || destination != "DPS" || departureDate != "2025-12-15" {
		return LionAirResponse{Success: true, Data: struct {
			AvailableFlights []LionFlight "json:\"available_flights\""
		}{}}, nil
	}

	var response LionAirResponse
	_ = json.Unmarshal([]byte(strLionAirMockResponse), &response)

	return response, nil
}
