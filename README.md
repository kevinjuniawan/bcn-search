- **bin/**: Compiled binaries and executables.
- **cmd/http/**: Entry point for the HTTP server application.
- **internal/**: Private application logic, not intended for external use.
- **adapter/http/**: HTTP delivery layer (controllers, handlers).
- **adapter/grpc/**: gRPC delivery layer.
- **adapter/event/**: Event-driven adapters (e.g., message brokers).
- **infrastructure/db/**: Database implementations and related infrastructure code.
- **config/**: Configuration files and settings.
- **pkg/logger/**: Logging utilities and abstractions.
- **pkg/validator/**: Input validation utilities.
- **test/**: Test files and test data.
- **scripts/**: Helper scripts for development and deployment.
- **proto/**: Protocol buffer definitions for gRPC or other services.


# API Doc

## SearchFlight API

URI : /v1/search-flight
Method : POST

Request Body :
```json
{
    "origin": "CGK",
    "destination": "DPS",
    "departure_date": "2025-12-15",
    "cabin_class": "economy",
    "sort_type": 3, // 0 Best value, 1 Lowest Price, 2 Highest Price, 3 Shortest duration, 4 Longest duration, 5 Departure time, 6 Arrival time
    "filter": {
        "airlines": ["Garuda", "Citilink"],
        "stopover": 0,
        "price": {
            "lowest_price": 100000,
            "highest_price": 1500000
        },
        "stop": 0,
        "time_range":{
            "type": "departure",
            "from": "2025-12-15T00:00:00Z",
            "to": "2025-12-15T23:59:59Z"
        }
    }
}
```

## Technical Design : [Click Here](https://drive.google.com/file/d/16Ua99vXi1d9bgk3nZL5ZWBohBREQQv6q/view?usp=sharing)
