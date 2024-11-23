package store

type DataStore interface {
	GetString(key string) (string, bool)
	SetString(key string, value string)
	GetInt(key string) (int, bool)
	SetInt(key string, value int)
	GetFloat(key string) (float64, bool)
	SetFloat(key string, value float64)
}
