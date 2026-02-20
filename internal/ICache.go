package internal

import (
	"context"
	"time"
)

type ICache interface {
	GetSortedFlightsByParams(ctx context.Context, params GetFlightsParams) ([]Flight, error)
	GetSortedFlightsByPrice(ctx context.Context, origin, destination, departureDate string, isAscending bool) ([]Flight, error)
	GetSortedFlightsByDuration(ctx context.Context, origin, destination, departureDate string, isAscending bool) ([]Flight, error)
	GetSortedFlightsByDepartureTime(ctx context.Context, origin, destination, departureDate string, isAscending bool) ([]Flight, error)
	GetSortedFlightsByArrivalTime(ctx context.Context, origin, destination, departureDate string, isAscending bool) ([]Flight, error)
	GetSortedFlightsByBestValue(ctx context.Context, origin, destination, departureDate string, isAscending bool) ([]Flight, error)
	SetFlights(ctx context.Context, flights []Flight, params GetFlightsParams, ttl time.Duration) error
}
