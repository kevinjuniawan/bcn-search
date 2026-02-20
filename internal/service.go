package internal

import (
	"context"
	"kevinjuniawan/bookcabin/config"
	"kevinjuniawan/bookcabin/pkg/helper"
	"sort"
	"time"

	"github.com/go-redis/redis/v8"
)

type InternalService struct {
	FetcherService IFetcher
	CacheService   ICache
	Cfg            config.Config
}

type InternalServiceParams struct {
	FetcherService IFetcher
	CacheService   ICache
}

func NewInternalService(params InternalServiceParams) *InternalService {
	return &InternalService{
		FetcherService: params.FetcherService,
		CacheService:   params.CacheService,
	}
}

func (s *InternalService) GetFlights(ctx context.Context, params GetFlightsParams) (SearchResponse, error) {
	startSearch := time.Now()
	providerCount := int16(0)
	succeededProvider := int16(0)
	isCache := true
	flightsList, err := s.CacheService.GetSortedFlightsByParams(ctx, params)
	if err != nil {
		if err != redis.Nil {
			return SearchResponse{Metadata: Metadata{IsCache: isCache}}, err
		}
		isCache = false
		flightsData, err := s.FetcherService.GetFlights(params)
		if err != nil {
			return SearchResponse{Metadata: Metadata{IsCache: isCache}}, err
		}
		flightsList = flightsData.Flights
		go func() {
			ctxSet := context.Background()
			s.CacheService.SetFlights(ctxSet, flightsList, params, time.Duration(0))
		}()
		flightsList = s.sortFlight(flightsList, params.SortType)
		providerCount = flightsData.ProviderCount
		succeededProvider = flightsData.ProviderCount - flightsData.FailedProvider
	}

	if params.Filter != nil {
		flightsList = FilterFlight(flightsList, *params.Filter)
	}

	duration := time.Since(startSearch)

	return SearchResponse{
		Metadata: Metadata{
			ProviderCount:     providerCount,
			SucceededProvider: succeededProvider,
			SearchTimeMs:      int32(duration.Milliseconds()),
			IsCache:           isCache,
		},
		Flights: flightsList,
	}, nil
}

func FilterFlight(flights []Flight, params FilterFlightParams) []Flight {
	filteredFlights := []Flight{}
	for _, flight := range flights {
		if len(params.Airline) != 0 {
			if exist := helper.ExistInSliceString(params.Airline, flight.Airline.Code); !exist {
				continue
			}
		}

		if params.Price != nil {
			if flight.Price.AmountInIDR() >= params.Price.LowestPrice && flight.Price.AmountInIDR() <= params.Price.HighestPrice {
				continue
			}
		}

		if params.Stops != nil {
			if flight.Stops != *params.Stops {
				continue
			}
		}

		if params.TimeRange != nil {
			fromTime, _ := time.Parse(time.RFC3339, params.TimeRange.From)
			toTime, _ := time.Parse(time.RFC3339, params.TimeRange.To)
			if params.TimeRange.Type == FilterFlightTimeTypeDeparture {
				if flight.Departure.Timestamp >= fromTime.Unix() && flight.Departure.Timestamp <= toTime.Unix() {
					continue
				}
			} else if params.TimeRange.Type == FilterFlightTimeTypeArrival {
				if flight.Arrival.Timestamp >= fromTime.Unix() && flight.Arrival.Timestamp <= toTime.Unix() {
					continue
				}
			}
		}
		filteredFlights = append(filteredFlights, flight)
	}
	return filteredFlights
}

func (s *InternalService) sortFlight(flights []Flight, sortType SortType) []Flight {
	switch sortType {
	case SortLowestPriceType:
		sort.Slice(flights, func(i, j int) bool {
			return flights[i].Price.AmountInIDR() < flights[j].Price.AmountInIDR()
		})
	case SortHighestPriceType:
		sort.Slice(flights, func(i, j int) bool {
			return flights[i].Price.AmountInIDR() > flights[j].Price.AmountInIDR()
		})
	case SortShortestDurationType:
		sort.Slice(flights, func(i, j int) bool {
			return flights[i].Duration.TotalMinute < flights[j].Duration.TotalMinute
		})
	case SortLongestDurationType:
		sort.Slice(flights, func(i, j int) bool {
			return flights[i].Duration.TotalMinute > flights[j].Duration.TotalMinute
		})
	case SortDepartureType:
		sort.Slice(flights, func(i, j int) bool {
			return flights[i].Departure.Timestamp < flights[j].Departure.Timestamp
		})
	case SortArrivalType:
		sort.Slice(flights, func(i, j int) bool {
			return flights[i].Arrival.Timestamp < flights[j].Arrival.Timestamp
		})
	case SortBestValueType:
		sort.Slice(flights, func(i, j int) bool {
			return CalculateBestValue(flights[i]) < CalculateBestValue(flights[j])
		})
	}
	return flights
}

func CalculateBestValue(flight Flight) int {
	multiplierClass := 1
	if flight.CabinClass == BusinessClass {
		multiplierClass = 2 // Infer only 2 class cabin
	}
	return flight.Price.AmountInIDR() + int(float32(flight.Duration.TotalMinute*10000)/float32(multiplierClass))
}
