package internal

type FlightDataResponse struct {
	ProviderCount  int16
	FailedProvider int16
	Flights        []Flight
}

type IFetcher interface {
	GetFlights(params GetFlightsParams) (FlightDataResponse, error)
}
