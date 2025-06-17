package grfrf

import (
	"fmt"
	"log"

	"github.com/leavemealonemf/go-route-filter/fkalman"
	glocation "github.com/leavemealonemf/go-route-filter/google-geolocation"
	"github.com/leavemealonemf/go-route-filter/gps"
	"googlemaps.github.io/maps"
)

type GoogleMapsApi struct {
	ApiKey string
}

type KalmanFilterParams struct {
	InitFilternitialEstimate float64
	ProcessNoise             float64
	MeasurementNoise         float64
}

type FilterInitializeData struct {
	GoogleMapsApi      *GoogleMapsApi
	KalmanFilterParams *KalmanFilterParams
}

type Filter struct {
	googleMapsConnection *maps.Client
	kalmanF              *fkalman.KalmanFilter
}

func (f *Filter) DeadReconing(prev *gps.Packet, curr *gps.Packet) *gps.Packet {
	if f.kalmanF != nil {
		if curr.Azimuth == 0 {
			curr.Azimuth = f.kalmanF.Update(f.kalmanF.X)
		} else {
			curr.Azimuth = f.kalmanF.Update(curr.Azimuth)
		}
	}

	lat, lon := gps.DeadReckoning(prev, curr)

	curr.Lat = lat
	curr.Lon = lon

	return curr
}

// return false if distance beetwen points greatest then provided threshold
func (f *Filter) CompareDistanceBetweenPoints(a, b *gps.Point, distanseThresholdM int) bool {
	distance := gps.CalculateDistance(a, b)
	fmt.Println("CALCULATED DISTANCE", distance)
	return distance <= float64(distanseThresholdM)
}

func (f *Filter) FindTower() {

}

func InitFilter(f *FilterInitializeData) *Filter {
	filter := &Filter{}

	if f.KalmanFilterParams != nil {
		filter.kalmanF = UseKalmanFilter(f.KalmanFilterParams)
	}

	if f.GoogleMapsApi != nil {
		c, err := glocation.InitMapsConnection(f.GoogleMapsApi.ApiKey)
		if err != nil {
			log.Println(err)
		} else {
			filter.googleMapsConnection = c
		}
	}
	return filter
}

// this is not require function, privide it in the filter constuructor if you want use it
func UseKalmanFilter(kalmanFilterParams *KalmanFilterParams) *fkalman.KalmanFilter {
	return fkalman.NewKalmanFilter(
		kalmanFilterParams.InitFilternitialEstimate,
		kalmanFilterParams.ProcessNoise,
		kalmanFilterParams.MeasurementNoise,
	)
}
