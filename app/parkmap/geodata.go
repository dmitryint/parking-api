package parkmap

import "main/app/parkmap/geo"

type GeoData struct {
	Type     string        `json:"type"`
	Polygons []geo.Polygon `json:"coordinates"`
}
