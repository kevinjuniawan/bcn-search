package api

import (
	"kevinjuniawan/bookcabin/config"
	mockflight "kevinjuniawan/bookcabin/infrastructure/api/mock_flight"
	"kevinjuniawan/bookcabin/internal"
	"log"
	"sync"
	"time"
)

type FetcherService struct {
	AirAsiaAPI   *mockflight.AirAsia
	GarudaAirAPI *mockflight.GarudaAir
	LionAirAPI   *mockflight.LionAir
	BatikAirAPI  *mockflight.BatikAir
	Cfg          config.Config
}

type FlightCollection struct {
	Flights []internal.Flight
	err     error
}

type FetcherServiceParams struct {
	Cfg config.Config
}

func NewFetcherService(params FetcherServiceParams) *FetcherService {
	return &FetcherService{
		Cfg: params.Cfg,
	}
}

func (f *FetcherService) GetFlights(params internal.GetFlightsParams) (internal.FlightDataResponse, error) {

	var wg sync.WaitGroup
	flightsChan := make(chan FlightCollection, 4)

	wg.Add(4)
	go func() {
		defer wg.Done()
		res, err := f.AirAsiaAPI.GetFlights(params.Origin, params.Destination, params.DepartureDate)
		flightsChan <- FlightCollection{Flights: res.Normalize(), err: err}
	}()
	go func() {
		defer wg.Done()
		res, err := f.GarudaAirAPI.GetFlights(params.Origin, params.Destination, params.DepartureDate)
		flightsChan <- FlightCollection{Flights: res.Normalize(), err: err}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < f.Cfg.MaxRetryCount; i++ {
			res, err := f.LionAirAPI.GetFlights(params.Origin, params.Destination, params.DepartureDate)
			if err == nil {
				flightsChan <- FlightCollection{Flights: res.Normalize(), err: err}
				break
			}
			time.Sleep(time.Duration(f.Cfg.RetryBackOff) * time.Millisecond) // Usually Resti will handle this
			log.Println("Retrying lion air fetching...")
		}
	}()
	go func() {
		defer wg.Done()
		res, err := f.BatikAirAPI.GetFlights(params.Origin, params.Destination, params.DepartureDate)
		flightsChan <- FlightCollection{Flights: res.Normalize(), err: err}
	}()

	wg.Wait()
	close(flightsChan)

	var flights []internal.Flight
	failed := 0
	for flight := range flightsChan {
		if flight.err != nil {
			failed++
		}
		flights = append(flights, flight.Flights...)
	}

	return internal.FlightDataResponse{
		ProviderCount:  4,
		FailedProvider: int16(failed),
		Flights:        flights,
	}, nil
}
