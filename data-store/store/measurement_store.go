package store

type MeasurementStore interface {
	AppendItem(
		measurement string,
		location string,
		field string,
		value float64,
	)
	Flush()
}
