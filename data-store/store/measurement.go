package store

import "time"

type Measurement struct {
	DeviceId        string
	MeasurementType string
	Field           string
	Value           any
	UpdatedAt       time.Time
}

func (m Measurement) Matches(queryParams map[string]string) bool {
	for paramKey, paramValue := range queryParams {
		switch paramKey {
		case "deviceId":
			if m.DeviceId != paramValue {
				return false
			}
		case "measurementType":
			if m.MeasurementType != paramValue {
				return false
			}
		case "field":
			if m.Field != paramValue {
				return false
			}
		}
	}
	return true
}
