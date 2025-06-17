package gps

import (
	"math"
	"time"

	"github.com/leavemealonemf/go-route-filter/utils"
)

type Point struct {
	Lat float64
	Lon float64
}

type Packet struct {
	Lat       float64 `json:"lat"`
	Lon       float64 `json:"lon"`
	SpeedKPH  float64 `json:"speed"`
	SpeedGPS  float64 `json:"speed_gps"`
	Azimuth   float64 `json:"azimuth"`
	Movement  bool    `json:"movement"`
	Timestamp int64   `json:"timestamp"` // Unix microseconds
	Time      time.Time
	MCC       int64 `json:"mcc"`
	MNC       int64 `json:"mnc"`
	LAC       int64 `json:"lac"`
	CellID    int64 `json:"cell_id"`
}

const earthRadiusM = 6371000

// Return correct lat & lon
func DeadReckoning(prev Packet, current Packet) (float64, float64) {
	latRad := utils.DegToRad(prev.Lat)
	lonRad := utils.DegToRad(prev.Lon)
	bearingRad := utils.DegToRad(current.Azimuth)

	deltaTime := current.Time.Sub(prev.Time).Seconds()
	speedMPS := current.SpeedKPH * 1000 / 3600
	distance := speedMPS * deltaTime

	newLatRad := math.Asin(math.Sin(latRad)*math.Cos(distance/earthRadiusM) +
		math.Cos(latRad)*math.Sin(distance/earthRadiusM)*math.Cos(bearingRad))

	newLonRad := lonRad + math.Atan2(
		math.Sin(bearingRad)*math.Sin(distance/earthRadiusM)*math.Cos(latRad),
		math.Cos(distance/earthRadiusM)-math.Sin(latRad)*math.Sin(newLatRad),
	)

	return utils.RadToDeg(newLatRad), utils.RadToDeg(newLonRad)
}

func CalculateDistance(a, b *Point) float64 {
	lat1Rad := utils.DegToRad(a.Lat)
	lon1Rad := utils.DegToRad(a.Lon)
	lat2Rad := utils.DegToRad(b.Lat)
	lon2Rad := utils.DegToRad(b.Lon)

	deltaLat := lat2Rad - lat1Rad
	deltaLon := lon2Rad - lon1Rad

	a_calc := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(deltaLon/2)*math.Sin(deltaLon/2)

	c := 2 * math.Atan2(math.Sqrt(a_calc), math.Sqrt(1-a_calc))

	distance := earthRadiusM * c

	return distance
}
