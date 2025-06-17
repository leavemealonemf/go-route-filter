package glocation

import (
	"fmt"

	"googlemaps.github.io/maps"
)

func InitMapsConnection(apiKey string) (*maps.Client, error) {
	c, err := maps.NewClient(maps.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("maps api not connected. started working without it\n%v", err)
	}

	return c, nil
}
