package store

import "time"

type Measurement struct {
	DeviceId        string
	MeasurementType string
	Field           string
	Value           any
	UpdatedAt       time.Time
}
