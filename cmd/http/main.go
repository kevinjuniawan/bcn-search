package main

import (
	"context"
	httpAdapter "kevinjuniawan/bookcabin/adapter/http"
	"kevinjuniawan/bookcabin/config"
	"kevinjuniawan/bookcabin/infrastructure/api"
	"kevinjuniawan/bookcabin/infrastructure/cache"
	"kevinjuniawan/bookcabin/internal"
	"log"
	"net/http"
	"strconv"
)

func main() {
	cfg, _ := config.Load()
	log.Printf("Initializing %s...\n", cfg.AppName)
	ctx := context.Background()
	api := api.NewFetcherService(api.FetcherServiceParams{Cfg: *cfg})
	redis := cache.NewCacheService(ctx, cache.ServiceParams{
		Address:  cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})
	internal := internal.NewInternalService(internal.InternalServiceParams{FetcherService: api, CacheService: redis})
	handler := httpAdapter.NewHandler(httpAdapter.Params{FlightService: internal})

	log.Printf("Starting listening for request on port %d \n", cfg.Port)
	http.ListenAndServe(":"+strconv.Itoa(cfg.Port), handler.InitRouter())
}
