package mockflight

import (
	"fmt"
	"kevinjuniawan/bookcabin/internal"
	"log"
	"strings"
	"time"
)

const strBatikAirMockResponse = `{"code":200,"message":"OK","results":[{"flightNumber":"ID6514","airlineName":"Batik Air","airlineIATA":"ID","origin":"CGK","destination":"DPS","departureDateTime":"2025-12-15T07:15:00+0700","arrivalDateTime":"2025-12-15T10:00:00+0800","travelTime":"1h 45m","numberOfStops":0,"fare":{"basePrice":980000,"taxes":120000,"totalPrice":1100000,"currencyCode":"IDR","class":"Y"},"seatsAvailable":32,"aircraftModel":"Airbus A320","baggageInfo":"7kg cabin, 20kg checked","onboardServices":["Snack","Beverage"]},{"flightNumber":"ID6520","airlineName":"Batik Air","airlineIATA":"ID","origin":"CGK","destination":"DPS","departureDateTime":"2025-12-15T13:30:00+0700","arrivalDateTime":"2025-12-15T16:20:00+0800","travelTime":"1h 50m","numberOfStops":0,"fare":{"basePrice":1050000,"taxes":130000,"totalPrice":1180000,"currencyCode":"IDR","class":"Y"},"seatsAvailable":18,"aircraftModel":"Boeing 737-800","baggageInfo":"7kg cabin, 20kg checked","onboardServices":["Meal","Beverage","Entertainment"]},{"flightNumber":"ID7042","airlineName":"Batik Air","airlineIATA":"ID","origin":"CGK","destination":"DPS","departureDateTime":"2025-12-15T18:45:00+0700","arrivalDateTime":"2025-12-15T23:50:00+0800","travelTime":"3h 5m","numberOfStops":1,"connections":[{"stopAirport":"UPG","stopDuration":"55m"}],"fare":{"basePrice":850000,"taxes":100000,"totalPrice":950000,"currencyCode":"IDR","class":"Y"},"seatsAvailable":41,"aircraftModel":"Airbus A320","baggageInfo":"7kg cabin, 20kg checked","onboardServices":["Snack"]}]}`
const strLionAirMockResponse = `{"success":true,"data":{"available_flights":[{"id":"JT740","carrier":{"name":"Lion Air","iata":"JT"},"route":{"from":{"code":"CGK","name":"Soekarno-Hatta International","city":"Jakarta"},"to":{"code":"DPS","name":"Ngurah Rai International","city":"Denpasar"}},"schedule":{"departure":"2025-12-15T05:30:00","departure_timezone":"Asia/Jakarta","arrival":"2025-12-15T08:15:00","arrival_timezone":"Asia/Makassar"},"flight_time":105,"is_direct":true,"pricing":{"total":950000,"currency":"IDR","fare_type":"ECONOMY"},"seats_left":45,"plane_type":"Boeing 737-900ER","services":{"wifi_available":false,"meals_included":false,"baggage_allowance":{"cabin":"7 kg","hold":"20 kg"}}},{"id":"JT742","carrier":{"name":"Lion Air","iata":"JT"},"route":{"from":{"code":"CGK","name":"Soekarno-Hatta International","city":"Jakarta"},"to":{"code":"DPS","name":"Ngurah Rai International","city":"Denpasar"}},"schedule":{"departure":"2025-12-15T11:45:00","departure_timezone":"Asia/Jakarta","arrival":"2025-12-15T14:35:00","arrival_timezone":"Asia/Makassar"},"flight_time":110,"is_direct":true,"pricing":{"total":890000,"currency":"IDR","fare_type":"ECONOMY"},"seats_left":38,"plane_type":"Boeing 737-800","services":{"wifi_available":false,"meals_included":false,"baggage_allowance":{"cabin":"7 kg","hold":"20 kg"}}},{"id":"JT650","carrier":{"name":"Lion Air","iata":"JT"},"route":{"from":{"code":"CGK","name":"Soekarno-Hatta International","city":"Jakarta"},"to":{"code":"DPS","name":"Ngurah Rai International","city":"Denpasar"}},"schedule":{"departure":"2025-12-15T16:20:00","departure_timezone":"Asia/Jakarta","arrival":"2025-12-15T21:10:00","arrival_timezone":"Asia/Makassar"},"flight_time":230,"is_direct":false,"stop_count":1,"layovers":[{"airport":"SUB","duration_minutes":75}],"pricing":{"total":780000,"currency":"IDR","fare_type":"ECONOMY"},"seats_left":52,"plane_type":"Boeing 737-800","services":{"wifi_available":false,"meals_included":false,"baggage_allowance":{"cabin":"7 kg","hold":"20 kg"}}}]}}`
const strGarudaAirMockResponse = `{"status":"success","flights":[{"flight_id":"GA400","airline":"Garuda Indonesia","airline_code":"GA","departure":{"airport":"CGK","city":"Jakarta","time":"2025-12-15T06:00:00+07:00","terminal":"3"},"arrival":{"airport":"DPS","city":"Denpasar","time":"2025-12-15T08:50:00+08:00","terminal":"I"},"duration_minutes":110,"stops":0,"aircraft":"Boeing 737-800","price":{"amount":1250000,"currency":"IDR"},"available_seats":28,"fare_class":"economy","baggage":{"carry_on":1,"checked":2},"amenities":["wifi","meal","entertainment"]},{"flight_id":"GA410","airline":"Garuda Indonesia","airline_code":"GA","departure":{"airport":"CGK","city":"Jakarta","time":"2025-12-15T09:30:00+07:00","terminal":"3"},"arrival":{"airport":"DPS","city":"Denpasar","time":"2025-12-15T12:25:00+08:00","terminal":"I"},"duration_minutes":115,"stops":0,"aircraft":"Airbus A330-300","price":{"amount":1450000,"currency":"IDR"},"available_seats":15,"fare_class":"economy","baggage":{"carry_on":1,"checked":2},"amenities":["wifi","power_outlet","meal","entertainment"]},{"flight_id":"GA315","airline":"Garuda Indonesia","airline_code":"GA","departure":{"airport":"CGK","city":"Jakarta","time":"2025-12-15T14:00:00+07:00","terminal":"3"},"arrival":{"airport":"SUB","city":"Surabaya","time":"2025-12-15T15:30:00+07:00","terminal":"2"},"duration_minutes":90,"stops":0,"aircraft":"Boeing 737","price":{"amount":1850000,"currency":"IDR"},"segments":[{"flight_number":"GA315","departure":{"airport":"CGK","time":"2025-12-15T14:00:00+07:00"},"arrival":{"airport":"SUB","time":"2025-12-15T15:30:00+07:00"},"duration_minutes":90},{"flight_number":"GA332","departure":{"airport":"SUB","time":"2025-12-15T17:15:00+07:00"},"arrival":{"airport":"DPS","time":"2025-12-15T18:45:00+08:00"},"duration_minutes":90,"layover_minutes":105}],"available_seats":22,"fare_class":"economy","baggage":{"carry_on":1,"checked":2}}]}`
const strAirasiaMockResponse = `{"status":"ok","flights":[{"flight_code":"QZ520","airline":"AirAsia","from_airport":"CGK","to_airport":"DPS","depart_time":"2025-12-15T04:45:00+07:00","arrive_time":"2025-12-15T07:25:00+08:00","duration_hours":1.67,"direct_flight":true,"price_idr":650000,"seats":67,"cabin_class":"economy","baggage_note":"Cabin baggage only, checked bags additional fee"},{"flight_code":"QZ524","airline":"AirAsia","from_airport":"CGK","to_airport":"DPS","depart_time":"2025-12-15T10:00:00+07:00","arrive_time":"2025-12-15T12:45:00+08:00","duration_hours":1.75,"direct_flight":true,"price_idr":720000,"seats":54,"cabin_class":"economy","baggage_note":"Cabin baggage only, checked bags additional fee"},{"flight_code":"QZ532","airline":"AirAsia","from_airport":"CGK","to_airport":"DPS","depart_time":"2025-12-15T19:30:00+07:00","arrive_time":"2025-12-15T22:10:00+08:00","duration_hours":1.67,"direct_flight":true,"price_idr":595000,"seats":72,"cabin_class":"economy","baggage_note":"Cabin baggage only, checked bags additional fee"},{"flight_code":"QZ7250","airline":"AirAsia","from_airport":"CGK","to_airport":"DPS","depart_time":"2025-12-15T15:15:00+07:00","arrive_time":"2025-12-15T20:35:00+08:00","duration_hours":4.33,"direct_flight":false,"stops":[{"airport":"SOC","wait_time_minutes":95}],"price_idr":485000,"seats":88,"cabin_class":"economy","baggage_note":"Cabin baggage only, checked bags additional fee"}]}`

var mapAirportCodeToCity = map[string]string{
	"CGK": "Jakarta",
	"DPS": "Denpasar",
}

var mapBaggageTypeToText = map[int8]string{
	1: "Not Available",
	2: "Allowed",
	3: "Additional Fee",
}

func GetCityNameFromAirportCode(originCode string, destinationCode string) (origin string, destination string, isValid bool) {
	departureCity := ""
	if departureCityCode, exist := mapAirportCodeToCity[originCode]; !exist {
		return "", "", false
	} else {
		departureCity = departureCityCode
	}

	arrivalCity := ""
	if arrivalCityCode, exist := mapAirportCodeToCity[destinationCode]; !exist {
		return "", "", false
	} else {
		arrivalCity = arrivalCityCode
	}

	return departureCity, arrivalCity, true
}

type FlightTime struct {
	Time time.Time
}

type FlightDuration struct {
	Duration time.Duration
}

func (f FlightTime) Format() string {
	return f.Time.Format(time.RFC3339)
}

func (f FlightDuration) Format() string {
	totalMinutes := int(f.Duration.Minutes())
	hours := totalMinutes / 60
	minutes := totalMinutes % 60
	return fmt.Sprintf("%dh %dm", hours, minutes)
}

func ConvertDateTimeToTime(format string, departureDateTime string, arrivalDateTime string) (departureTimestamp FlightTime, arrivalTimestamp FlightTime, duration FlightDuration, isValid bool) {
	var err error
	departureTimestamp.Time, err = time.Parse(format, departureDateTime)
	if err != nil {
		log.Println("fail parsing depature time, Err : ", err)
		return FlightTime{}, FlightTime{}, FlightDuration{}, false
	}
	arrivalTimestamp.Time, err = time.Parse(format, arrivalDateTime)
	if err != nil {
		log.Println("fail parsing arrival time, Err : ", err)
		return FlightTime{}, FlightTime{}, FlightDuration{}, false
	}

	duration = FlightDuration{Duration: arrivalTimestamp.Time.Sub(departureTimestamp.Time)}
	if duration.Duration < 0 {
		log.Println("arrival time is before departure time")
		return FlightTime{}, FlightTime{}, FlightDuration{}, false
	}

	return departureTimestamp, arrivalTimestamp, duration, true
}

func GetBaggageInfo(baggageData string) internal.Bag {
	baggageProviderData := strings.Split(baggageData, ",")
	baggage := internal.Bag{}
	if len(baggageProviderData) > 0 {
		baggage.CarryOn = baggageProviderData[0]
		baggage.Checked = baggageProviderData[1]
	}
	return baggage
}

type AirlineCode string

const (
	AirAsiaAirline  AirlineCode = "QZ"
	GarudaAirline   AirlineCode = "GA"
	LionAirAirline  AirlineCode = "JT"
	BatikAirAirline AirlineCode = "ID"
)

var mapAirlineCodeToName = map[AirlineCode]string{
	AirAsiaAirline:  "AirAsia",
	GarudaAirline:   "GarudaIndonesia",
	LionAirAirline:  "LionAir",
	BatikAirAirline: "BatikAir",
}
