package parkmap

import (
	"embed"
	"encoding/json"
	"fmt"
	"golang.org/x/text/encoding/charmap"
	"io/ioutil"
	"main/app/parkmap/geo"
)

type ParkingMap struct {
	entities []ParkingEntity
}

// https://data.mos.ru/opendata/7704786030-platnye-parkovki-na-ulichno-dorojnoy-seti
//go:embed data/data-4905-2021-07-16.json
var f embed.FS

func NewParkingMap() *ParkingMap {
	jsonFile, err := f.Open("data/data-4905-2021-07-16.json")
	if err != nil {
		panic(err)
	}
	r := charmap.Windows1251.NewDecoder().Reader(jsonFile)
	byteValue, err := ioutil.ReadAll(r)
	keys := make([]ParkingEntity, 0)
	json.Unmarshal(byteValue, &keys)

	return &ParkingMap{entities: keys}
}

func (p *ParkingMap) Test() {
	for _, e := range p.entities {
		//e.GeoData.Print()
		for _, p := range e.GeoData.Polygons {
			fmt.Printf("%s\n", p.IsClosed())
		}
	}
}

func (p *ParkingMap) Closest(point *geo.Point, paid bool, disabled bool, minDistanceKm float64) (error, *ParkingEntity) {
	minMatch := false
	minEntity := ParkingEntity{}
	min := .0

	for _, e := range p.entities {
		if paid && !e.IsPaid() {
			continue
		}

		if (disabled && !e.IsDisabledAlowed()) || (!disabled && !e.IsNonDisabledAlowed()) {
			continue
		}

		for _, parkingSlot := range e.GeoData.Polygons {
			if parkingSlot.Contains(point) {
				fmt.Printf("Contains: %s\n", e.Id)
				return nil, &e
			}
			if minMatch && (min > parkingSlot.MinimumDistance(point)) {
				min = parkingSlot.MinimumDistance(point)
				minEntity = e
			}
			if !minMatch {
				minMatch = true
				min = parkingSlot.MinimumDistance(point)
				minEntity = e
			}
		}
	}
	if !minMatch {
		return fmt.Errorf("Match not found"), &minEntity
	}
	if min > minDistanceKm {
		return fmt.Errorf("Maximum distance exceeded"), &minEntity
	}
	return nil, &minEntity
}
