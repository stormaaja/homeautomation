package store

type Measurement struct {
	DeviceId        string
	MeasurementType string
	Field           string
	Value           any
}

type DataStore interface {
	GetMeasurement(
		measurementType string,
		key string,
	) (Measurement, bool)

	SetMeasurement(
		measurementType string,
		key string,
		measurement Measurement,
	)
}
