package mockflight

import (
	"encoding/json"
	"kevinjuniawan/bookcabin/internal"
	"kevinjuniawan/bookcabin/pkg/helper"
	"log"
	"time"
)

type IDPrice struct {
	BasePrice    int    `json:"basePrice"`
	Taxes        int    `json:"taxes"`
	TotalPrice   int    `json:"totalPrice"`
	CurrencyCode string `json:"currencyCode"`
	Class        string `json:"class"`
}

type IDTransit struct {
	StopAirport  string `json:"stopAirport"`
	StopDuration string `json:"stopDuration"`
}

type BatikFlight struct {
	FlightNumber      string      `json:"flightNumber"`
	AirlineName       int         `json:"airlineName"` // Specified as int in model
	AirlineIATA       string      `json:"airlineIATA"`
	Origin            string      `json:"origin"`
	Destination       string      `json:"destination"`
	DepartureDateTime string      `json:"departureDateTime"`
	ArrivalDateTime   string      `json:"arrivalDateTime"`
	TravelTime        string      `json:"travel_time"`
	NumberOfStops     int8        `json:"numberOfStops"`
	Fare              IDPrice     `json:"fare"`
	SeatAvailable     int16       `json:"seatAvailable"`
	AircraftModel     string      `json:"aircraftModel"`
	BaggageInfo       string      `json:"baggageInfo"`
	OnBoardServices   []string    `json:"onBoardServices"`
	Connections       []IDTransit `json:"connections"` // Nullable
}

type BatikAirResponse struct {
	Code    int           `json:"code"`
	Message string        `json:"message"`
	Results []BatikFlight `json:"results"`
}

func (b BatikFlight) IsValid() bool {
	if b.FlightNumber == "" || b.Origin == "" || b.Destination == "" {
		return false
	}

	if b.Fare.TotalPrice == 0 {
		return false
	}

	return true
}

func (b BatikAirResponse) Normalize() []internal.Flight {
	flights := []internal.Flight{}
	for _, flight := range b.Results {
		if !flight.IsValid() {
			log.Printf("[%s - %s] invalid flight data\n", BatikAirAirline, flight.FlightNumber)
			continue
		}

		departureCity, arrivalCity, isValid := GetCityNameFromAirportCode(flight.Origin, flight.Destination)
		if !isValid {
			log.Printf("[%s - %s]fail map origin/destination airport code\n", BatikAirAirline, flight.FlightNumber)
			continue
		}

		departureTime, arrivalTime, duration, isValid := ConvertDateTimeToTime("2006-01-02T15:04:05-0700", flight.DepartureDateTime, flight.ArrivalDateTime)
		if !isValid {
			log.Printf("[%s - %s]fail parse departure/arrival time\n", BatikAirAirline, flight.FlightNumber)
			continue
		}

		class := internal.BusinessClass
		if flight.Fare.Class == "Y" {
			class = internal.EconomyClass
		}

		baggageData := GetBaggageInfo(flight.BaggageInfo)

		layover := 0
		for _, stop := range flight.Connections {
			stopDuration, err := time.ParseDuration(stop.StopDuration)
			if err != nil {
				log.Printf("[%s - %s]fail convert stop duration, Err : %v\n", BatikAirAirline, flight.FlightNumber, err)
				continue
			}
			layover += int(stopDuration.Minutes())
		}

		flights = append(flights, internal.Flight{
			ID:             flight.FlightNumber + "_" + mapAirlineCodeToName[BatikAirAirline],
			Provider:       mapAirlineCodeToName[BatikAirAirline],
			Airline:        internal.Airline{Code: mapAirlineCodeToName[BatikAirAirline], Name: string(BatikAirAirline)},
			FlightNumber:   flight.FlightNumber,
			Departure:      internal.Airport{Airport: flight.Origin, City: departureCity, Datetime: departureTime.Format(), Timestamp: departureTime.Time.Unix()},
			Arrival:        internal.Airport{Airport: flight.Destination, City: arrivalCity, Datetime: arrivalTime.Format(), Timestamp: arrivalTime.Time.Unix()},
			Duration:       internal.Duration{TotalMinute: int16(duration.Duration.Minutes()), Formatted: duration.Format()},
			Stops:          flight.NumberOfStops,
			Price:          internal.Price{Amount: flight.Fare.TotalPrice, Currency: flight.Fare.CurrencyCode},
			AvailableSeats: flight.SeatAvailable,
			CabinClass:     class,
			Aircraft:       &flight.AircraftModel,
			Amenities:      flight.OnBoardServices,
			Baggage:        baggageData,
			Layover:        layover,
		})
	}
	return flights
}

type BatikAir struct{}

func NewBatikAir() *BatikAir {
	return &BatikAir{}
}

func (b *BatikAir) GetFlights(origin, destination, departureDate string) (BatikAirResponse, error) {
	helper.RandomDelay(200, 400)
	if origin != "CGK" || destination != "DPS" || departureDate != "2025-12-15" {
		return BatikAirResponse{Code: 200, Message: "success", Results: []BatikFlight{}}, nil
	}
	var response BatikAirResponse
	_ = json.Unmarshal([]byte(strBatikAirMockResponse), &response)

	return response, nil
}
