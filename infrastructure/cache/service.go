package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"kevinjuniawan/bookcabin/config"
	"kevinjuniawan/bookcabin/internal"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

type CacheService struct {
	Client *redis.Client
	Cfg    *config.Config
}

type ServiceParams struct {
	Address  string
	Password string
	DB       int
	Cfg      *config.Config
}

func NewCacheService(ctx context.Context, params ServiceParams) *CacheService {
	client := redis.NewClient(&redis.Options{
		Addr:     params.Address,
		Password: params.Password,
		DB:       params.DB,
	})
	if err := client.Ping(ctx).Err(); err != nil {
		panic(err)
	}
	return &CacheService{
		Client: client,
		Cfg:    params.Cfg,
	}
}

func (c *CacheService) GetSortedFlightsByParams(ctx context.Context, params internal.GetFlightsParams) (flightList []internal.Flight, err error) {
	key, isAscending := makeKey(params)
	var sortedFlightIDs []redis.Z
	if isAscending {
		sortedFlightIDs, err = c.Client.ZRangeByScoreWithScores(ctx, key, &redis.ZRangeBy{
			Min: "-inf",
			Max: "+inf",
		}).Result()
	} else {
		sortedFlightIDs, err = c.Client.ZRevRangeByScoreWithScores(ctx, key, &redis.ZRangeBy{
			Min: "-inf",
			Max: "+inf",
		}).Result()
	}

	if err != nil {
		return nil, err
	}
	if len(sortedFlightIDs) == 0 {
		return nil, redis.Nil
	}
	return c.ContructFlightByZSetMember(ctx, sortedFlightIDs)
}

func (c *CacheService) GetSortedFlightsByPrice(ctx context.Context, origin, destination, departureDate string, isAscending bool) (flightList []internal.Flight, err error) {
	sortType := internal.SortLowestPriceType
	if !isAscending {
		sortType = internal.SortHighestPriceType
	}
	return c.GetSortedFlightsByParams(ctx, internal.GetFlightsParams{
		Origin:        origin,
		Destination:   destination,
		DepartureDate: departureDate,
		SortType:      sortType,
	})
}

func (c *CacheService) GetSortedFlightsByDuration(ctx context.Context, origin, destination, departureDate string, isAscending bool) (flightList []internal.Flight, err error) {
	sortType := internal.SortShortestDurationType
	if !isAscending {
		sortType = internal.SortLongestDurationType
	}
	return c.GetSortedFlightsByParams(ctx, internal.GetFlightsParams{
		Origin:        origin,
		Destination:   destination,
		DepartureDate: departureDate,
		SortType:      sortType,
	})
}

func (c *CacheService) GetSortedFlightsByDepartureTime(ctx context.Context, origin, destination, departureDate string, isAscending bool) (flightList []internal.Flight, err error) {
	return c.GetSortedFlightsByParams(ctx, internal.GetFlightsParams{
		Origin:        origin,
		Destination:   destination,
		DepartureDate: departureDate,
		SortType:      internal.SortDepartureType,
	})
}

func (c *CacheService) GetSortedFlightsByArrivalTime(ctx context.Context, origin, destination, departureDate string, isAscending bool) (flightList []internal.Flight, err error) {
	return c.GetSortedFlightsByParams(ctx, internal.GetFlightsParams{
		Origin:        origin,
		Destination:   destination,
		DepartureDate: departureDate,
		SortType:      internal.SortArrivalType,
	})
}

func (c *CacheService) GetSortedFlightsByBestValue(ctx context.Context, origin, destination, departureDate string, isAscending bool) (flightList []internal.Flight, err error) {
	return c.GetSortedFlightsByParams(ctx, internal.GetFlightsParams{
		Origin:        origin,
		Destination:   destination,
		DepartureDate: departureDate,
		SortType:      internal.SortBestValueType,
	})
}

func (c *CacheService) SetFlights(ctx context.Context, flights []internal.Flight, params internal.GetFlightsParams, ttl time.Duration) error {
	log.Printf("set %d flights to cache\n", len(flights))
	for _, flight := range flights {
		flightJSON, err := json.Marshal(flight)
		if err != nil {
			log.Printf("[%s]fail to marshal flight, Err : %v\n", flight.ID, err)
			return err
		}

		err = c.Client.Set(ctx, flight.ID, flightJSON, ttl).Err()
		if err != nil {
			log.Printf("[%s]fail to set flight, Err : %v\n", flight.ID, err)
			return err
		}

		err = c.addMemberToZSetPrice(ctx, flight, params)
		if err != nil {
			log.Printf("[%s]fail to add member to zset price, Err : %v\n", flight.ID, err)
			return err
		}

		err = c.addMemberToZSetDepartureTime(ctx, flight, params)
		if err != nil {
			log.Printf("[%s]fail to add member to zset departure time, Err : %v\n", flight.ID, err)
			return err
		}

		err = c.addMemberToZSetArrivalTime(ctx, flight, params)
		if err != nil {
			log.Printf("[%s]fail to add member to zset arrival time, Err : %v\n", flight.ID, err)
			return err
		}

		err = c.addMemberToZSetDuration(ctx, flight, params)
		if err != nil {
			log.Printf("[%s]fail to add member to zset duration, Err : %v\n", flight.ID, err)
			return err
		}

		err = c.addMemberToZSetBestValue(ctx, flight, params)
		if err != nil {
			log.Printf("[%s]fail to add member to zset best value, Err : %v\n", flight.ID, err)
			return err
		}
	}
	return nil
}

func (c *CacheService) addMemberToZSetPrice(ctx context.Context, flight internal.Flight, params internal.GetFlightsParams) error {
	params.SortType = internal.SortLowestPriceType
	key, _ := makeKey(params)
	return c.Client.ZAdd(ctx, key, &redis.Z{
		Score:  float64(flight.Price.AmountInIDR()),
		Member: flight.ID,
	}).Err()
}

func (c *CacheService) addMemberToZSetDepartureTime(ctx context.Context, flight internal.Flight, params internal.GetFlightsParams) error {
	params.SortType = internal.SortDepartureType
	key, _ := makeKey(params)
	return c.Client.ZAdd(ctx, key, &redis.Z{
		Score:  float64(flight.Departure.Timestamp),
		Member: flight.ID,
	}).Err()
}

func (c *CacheService) addMemberToZSetArrivalTime(ctx context.Context, flight internal.Flight, params internal.GetFlightsParams) error {
	params.SortType = internal.SortArrivalType
	key, _ := makeKey(params)
	return c.Client.ZAdd(ctx, key, &redis.Z{
		Score:  float64(flight.Arrival.Timestamp),
		Member: flight.ID,
	}).Err()
}

func (c *CacheService) addMemberToZSetDuration(ctx context.Context, flight internal.Flight, params internal.GetFlightsParams) error {
	params.SortType = internal.SortShortestDurationType
	key, _ := makeKey(params)
	return c.Client.ZAdd(ctx, key, &redis.Z{
		Score:  float64(flight.Duration.TotalMinute),
		Member: flight.ID,
	}).Err()
}

func (c *CacheService) addMemberToZSetBestValue(ctx context.Context, flight internal.Flight, params internal.GetFlightsParams) error {
	params.SortType = internal.SortBestValueType
	key, _ := makeKey(params)
	return c.Client.ZAdd(ctx, key, &redis.Z{
		Score:  float64(internal.CalculateBestValue(flight)),
		Member: flight.ID,
	}).Err()
}

func (c *CacheService) ContructFlightByZSetMember(ctx context.Context, setMember []redis.Z) (flightList []internal.Flight, err error) {
	flightIDs := make([]string, len(setMember))
	for i, member := range setMember {
		flightIDs[i] = member.Member.(string)
	}
	return c.ConstructFlightByFlightID(ctx, flightIDs)
}

func (c *CacheService) ConstructFlightByFlightID(ctx context.Context, flightIDs []string) (flightList []internal.Flight, err error) {
	flights, err := c.Client.MGet(ctx, flightIDs...).Result()
	if err != nil {
		return nil, err
	}
	for _, flight := range flights {
		var flightObj internal.Flight
		err = json.Unmarshal([]byte(flight.(string)), &flightObj)
		if err != nil {
			return nil, err
		}
		flightList = append(flightList, flightObj)
	}
	return flightList, nil
}

func makeKey(params internal.GetFlightsParams) (string, bool) {
	baseKey := fmt.Sprintf("flights:%s:%s:%s", params.Origin, params.Destination, params.DepartureDate)
	switch params.SortType {
	case internal.SortLowestPriceType:
		return baseKey + ":" + MapSortTypeToKey[params.SortType], true
	case internal.SortHighestPriceType:
		return baseKey + ":" + MapSortTypeToKey[params.SortType], false
	case internal.SortShortestDurationType:
		return baseKey + ":" + MapSortTypeToKey[params.SortType], true
	case internal.SortLongestDurationType:
		return baseKey + ":" + MapSortTypeToKey[params.SortType], false
	case internal.SortDepartureType:
		return baseKey + ":" + MapSortTypeToKey[params.SortType], true
	case internal.SortArrivalType:
		return baseKey + ":" + MapSortTypeToKey[params.SortType], true
	case internal.SortBestValueType:
		return baseKey + ":" + MapSortTypeToKey[params.SortType], true
	}
	return "", false
}

func (c *CacheService) IsRequestLimiterExceeded(ctx context.Context, URI string) bool {
	key := fmt.Sprintf("request_limiter:%s", URI)
	counter := c.Client.Incr(ctx, key).Val()
	if counter == 1 {
		c.Client.Expire(ctx, key, c.Cfg.RequestLimiterTTL)
	}
	if counter > c.Cfg.RequestLimiterMax+1 {
		return true
	}
	return false
}
