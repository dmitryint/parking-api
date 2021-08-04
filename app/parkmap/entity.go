package parkmap

type ParkingEntity struct {
	Id                  int      `json:"ID"`
	ParkingName         string   `json:"ParkingName"`
	ParkingZoneNumber   string   `json:"ParkingZoneNumber"`
	AdmArea             string   `json:"AdmArea"`
	District            string   `json:"District"`
	Address             string   `json:"Address"`
	CarCapacity         int      `json:"CarCapacity"`
	CarCapacityDisabled int      `json:"CarCapacityDisabled"`
	Tariffs             []Tariff `json:"Tariffs"`
	GeoData             GeoData  `json:"geoData"`
}

type Tariff struct {
	TariffType     string  `json:"TariffType"`
	TimeRange      string  `json:"TimeRange"`
	FirstHalfHour  float32 `json:"FirstHalfHour,string"`
	FirstHour      float32 `json:"FirstHour"`
	FollowingHours float32 `json:"FollowingHours"`
	HourPrice      float32 `json:"HourPrice"`
}

func (e *ParkingEntity) IsPaid() bool {
	for _, t := range e.Tariffs {
		if t.FirstHalfHour > 0 || t.FirstHour > 0 || t.FollowingHours > 0 || t.HourPrice > 0 {
			return true
		}
	}
	return false
}

func (e *ParkingEntity) IsDisabledAlowed() bool {
	if e.CarCapacityDisabled > 0 {
		return true
	}
	return false
}

func (e *ParkingEntity) IsNonDisabledAlowed() bool {
	if e.CarCapacity-e.CarCapacityDisabled > 0 {
		return true
	}
	return false
}
