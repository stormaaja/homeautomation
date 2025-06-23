package store

type Measurement struct {
	DeviceId        string
	MeasurementType string
	Field           string
	Value           any
}
