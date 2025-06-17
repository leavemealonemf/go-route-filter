package gps

import (
	"math"
	"time"

	"github.com/leavemealonemf/go-route-filter/utils"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geo"
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
func DeadReckoning(prev *Packet, current *Packet) (float64, float64) {
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

func CalculateDistance(point1, point2 *Point) float64 {

	// if point1.Lat == point2.Lat && point1.Lon == point2.Lon {
	// 	return 0
	// }

	// lat1Rad := utils.DegToRad(point1.Lat)
	// lon1Rad := utils.DegToRad(point1.Lon)
	// lat2Rad := utils.DegToRad(point2.Lat)
	// lon2Rad := utils.DegToRad(point2.Lon)

	// deltaLat := lat2Rad - lat1Rad
	// deltaLon := lon2Rad - lon1Rad

	// haversineA := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
	// 	math.Cos(lat1Rad)*math.Cos(lat2Rad)*
	// 		math.Sin(deltaLon/2)*math.Sin(deltaLon/2)

	// c := 2 * math.Atan2(math.Sqrt(haversineA), math.Sqrt(1-haversineA))

	// distance := earthRadiusM * c

	p1 := orb.Point{point1.Lat, point1.Lon}
	p2 := orb.Point{point2.Lat, point2.Lon}
	return geo.Distance(p1, p2)
}
