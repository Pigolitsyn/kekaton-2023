package server

import "kekaton/back/internal/storage"

type PointGeo struct {
	Type     string   `json:"type"`
	Geometry Geometry `json:"geometry"`
}

type Geometry struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

func pointsToGeo(points []storage.Point) []PointGeo {
	geos := make([]PointGeo, len(points))

	for i := range points {
		geos[i] = PointGeo{
			Type: "Feature",
			Geometry: Geometry{
				Type:        "Point",
				Coordinates: points[i].Coordinates[:],
			},
		}
	}

	return geos
}
