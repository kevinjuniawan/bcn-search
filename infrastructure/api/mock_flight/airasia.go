package mockflight

import (
	"encoding/json"
	"kevinjuniawan/bookcabin/internal"
	"kevinjuniawan/bookcabin/pkg/helper"
	"log"
	"time"
)

type QZTransit struct {
	Airport         string `json:"airport"`
	WaitTimeMinutes int    `json:"wait_time_minutes"`
}

type AirAsiaFlight struct {
	FlightCode    string      `json:"flight_code"`
	Airline       string      `json:"airline"`
	FromAirport   string      `json:"from_airport"`
	ToAirport     string      `json:"to_airport"`
	DepartTime    string      `json:"depart_time"`
	ArriveTime    string      `json:"arrive_time"`
	DurationHours float32     `json:"duration_hours"`
	DirectFlight  bool        `json:"direct_flight"`
	PriceIDR      int         `json:"price_idr"`
	Seats         int16       `json:"seats"`
	CabinClass    string      `json:"cabin_class"`
	BaggageNote   string      `json:"baggage_note"`
	Stops         []QZTransit `json:"stops,omitempty"`
}

func (a AirAsiaFlight) IsValid() bool {

	if a.FlightCode == "" || a.FromAirport == "" || a.ToAirport == "" {
		return false
	}

	if a.PriceIDR == 0 {
		return false
	}

	return true
}

type AirAsiaResponse struct {
	Status  string          `json:"status"`
	Flights []AirAsiaFlight `json:"flights"`
}

func (r AirAsiaResponse) Normalize() []internal.Flight {
	flights := []internal.Flight{}
	for _, flight := range r.Flights {

		if !flight.IsValid() {
			log.Printf("[%s - %s] invalid flight data\n", AirAsiaAirline, flight.FlightCode)
			continue
		}

		departureCity, arrivalCity, isValid := GetCityNameFromAirportCode(flight.FromAirport, flight.ToAirport)
		if !isValid {
			log.Printf("[%s - %s] fail map origin/destination airport code\n", AirAsiaAirline, flight.FlightCode)
			continue
		}

		departureTime, arrivalTime, duration, isValid := ConvertDateTimeToTime(time.RFC3339, flight.DepartTime, flight.ArriveTime)
		if !isValid {
			log.Printf("[%s - %s] fail parse departure/arrival time\n", AirAsiaAirline, flight.FlightCode)
			continue
		}

		cabinClass := internal.EconomyClass
		if flight.CabinClass == "business" {
			cabinClass = internal.BusinessClass
		}

		baggageData := GetBaggageInfo(flight.BaggageNote)

		layover := 0
		if len(flight.Stops) > 0 {
			for _, stop := range flight.Stops {
				layover += stop.WaitTimeMinutes
			}
		}

		flights = append(flights, internal.Flight{
			ID:             flight.FlightCode + "_" + mapAirlineCodeToName[AirAsiaAirline],
			Provider:       mapAirlineCodeToName[AirAsiaAirline],
			Airline:        internal.Airline{Code: mapAirlineCodeToName[AirAsiaAirline], Name: string(AirAsiaAirline)},
			FlightNumber:   flight.FlightCode,
			Departure:      internal.Airport{Airport: flight.FromAirport, City: departureCity, Datetime: departureTime.Format(), Timestamp: departureTime.Time.Unix()},
			Arrival:        internal.Airport{Airport: flight.ToAirport, City: arrivalCity, Datetime: arrivalTime.Format(), Timestamp: arrivalTime.Time.Unix()},
			Duration:       internal.Duration{TotalMinute: int16(duration.Duration.Minutes()), Formatted: duration.Format()},
			Stops:          int8(len(flight.Stops)),
			Price:          internal.Price{Amount: flight.PriceIDR, Currency: "IDR"},
			AvailableSeats: flight.Seats,
			CabinClass:     cabinClass,
			Aircraft:       nil,
			Amenities:      nil,
			Baggage:        baggageData,
			Layover:        layover,
		})
	}
	return flights
}

type AirAsia struct{}

func NewAirAsia() *AirAsia {
	return &AirAsia{}
}

func (b *AirAsia) GetFlights(origin, destination, departureDate string) (AirAsiaResponse, error) {
	helper.RandomDelay(50, 150)
	if origin != "CGK" || destination != "DPS" || departureDate != "2025-12-15" {
		return AirAsiaResponse{Status: "success", Flights: []AirAsiaFlight{}}, nil
	}
	err := helper.RandomError(10)
	if err != nil {
		return AirAsiaResponse{}, err
	}
	var response AirAsiaResponse
	_ = json.Unmarshal([]byte(strAirasiaMockResponse), &response)

	return response, nil
}
