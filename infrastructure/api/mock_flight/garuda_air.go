package mockflight

import (
	"encoding/json"
	"kevinjuniawan/bookcabin/internal"
	"kevinjuniawan/bookcabin/pkg/helper"
	"log"
)

type GAAirport struct {
	Airport  string `json:"airport"`
	City     string `json:"city"`
	Time     string `json:"time"`
	Terminal string `json:"terminal"`
}

type GAPrice struct {
	Amount   int    `json:"amount"`
	Currency string `json:"currency"`
}

type GABag struct {
	CarryOn int8 `json:"carry_on"`
	Checked int8 `json:"checked"`
}

type GATransit struct {
	FlightNumber    string    `json:"flight_number"`
	Departure       GAAirport `json:"departure"`
	Arrival         GAAirport `json:"arrival"`
	DurationMinutes int16     `json:"duration_minutes"`
	LayoverMinutes  int16     `json:"layover_minutes,omitempty"`
}

type GarudaFlight struct {
	FlightID        string      `json:"flight_id"`
	Airline         string      `json:"airline"`
	AirlineCode     string      `json:"airline_code"`
	Departure       GAAirport   `json:"departure"`
	Arrival         GAAirport   `json:"arrival"`
	DurationMinutes int         `json:"duration_minutes"`
	Stops           int8        `json:"stops"`
	Aircraft        string      `json:"aircraft"`
	Price           GAPrice     `json:"price"`
	AvailableSeats  int16       `json:"available_seats"`
	FareClass       string      `json:"fare_class"`
	Baggage         GABag       `json:"baggage"`
	Amenities       []string    `json:"amenities"`
	Segments        []GATransit `json:"segments,omitempty"`
}

func (g GarudaFlight) isValid() bool {
	if g.FlightID == "" || g.Departure.Airport == "" || g.Arrival.Airport == "" {
		return false
	}

	if g.Price.Amount == 0 {
		return false
	}

	return true
}

type GarudaAirResponse struct {
	Status  string         `json:"status"`
	Flights []GarudaFlight `json:"flights"`
}

func (g GarudaAirResponse) Normalize() []internal.Flight {
	flights := []internal.Flight{}
	for _, flight := range g.Flights {
		if !flight.isValid() {
			log.Printf("[%s - %s] invalid flight data\n", GarudaAirline, flight.FlightID)
			continue
		}

		arrivalCodeData := flight.Arrival.Airport
		arrivalTimeData := flight.Arrival.Time
		if len(flight.Segments) > 0 {
			arrivalCodeData = flight.Segments[len(flight.Segments)-1].Arrival.Airport
			arrivalTimeData = flight.Segments[len(flight.Segments)-1].Arrival.Time
		}

		departureCity, arrivalCity, isValid := GetCityNameFromAirportCode(flight.Departure.Airport, arrivalCodeData)
		if !isValid {
			log.Printf("[%s - %s] fail map origin/destination airport code\n", GarudaAirline, flight.FlightID)
			continue
		}

		departureTime, arrivalTime, duration, isValid := ConvertDateTimeToTime("2006-01-02T15:04:05-07:00", flight.Departure.Time, arrivalTimeData)
		if !isValid {
			log.Printf("[%s - %s] fail parse departure/arrival time\n", GarudaAirline, flight.FlightID)
			continue
		}

		class := internal.BusinessClass
		if flight.FareClass == "economy" {
			class = internal.EconomyClass
		}

		layover := 0
		for _, transit := range flight.Segments {
			layover += int(transit.LayoverMinutes)
		}

		flights = append(flights, internal.Flight{
			ID:             flight.FlightID + "_" + mapAirlineCodeToName[GarudaAirline],
			Provider:       mapAirlineCodeToName[GarudaAirline],
			Airline:        internal.Airline{Code: mapAirlineCodeToName[GarudaAirline], Name: string(GarudaAirline)},
			FlightNumber:   flight.FlightID,
			Departure:      internal.Airport{Airport: flight.Departure.Airport, City: departureCity, Datetime: departureTime.Format(), Timestamp: departureTime.Time.Unix()},
			Arrival:        internal.Airport{Airport: flight.Arrival.Airport, City: arrivalCity, Datetime: arrivalTime.Format(), Timestamp: arrivalTime.Time.Unix()},
			Duration:       internal.Duration{TotalMinute: int16(duration.Duration.Minutes()), Formatted: duration.Format()},
			Stops:          flight.Stops,
			Price:          internal.Price{Amount: flight.Price.Amount, Currency: flight.Price.Currency},
			AvailableSeats: flight.AvailableSeats,
			CabinClass:     class,
			Aircraft:       &flight.Aircraft,
			Amenities:      flight.Amenities,
			Baggage:        internal.Bag{CarryOn: mapBaggageTypeToText[flight.Baggage.CarryOn], Checked: mapBaggageTypeToText[flight.Baggage.Checked]},
			Layover:        layover,
		})

	}

	return flights
}

type GarudaAir struct{}

func NewGarudaAir() *GarudaAir {
	return &GarudaAir{}
}

func (b *GarudaAir) GetFlights(origin, destination, departureDate string) (GarudaAirResponse, error) {
	helper.RandomDelay(50, 150)
	if origin != "CGK" || destination != "DPS" || departureDate != "2025-12-15" {
		return GarudaAirResponse{Status: "success", Flights: []GarudaFlight{}}, nil
	}

	var response GarudaAirResponse
	_ = json.Unmarshal([]byte(strGarudaAirMockResponse), &response)

	return response, nil
}
