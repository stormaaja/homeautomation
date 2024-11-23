package store

type MemoryStore struct {
	Data map[string]interface{}
}

func (m *MemoryStore) GetString(key string) (string, bool) {
	if value, ok := m.Data[key]; !ok {
		return "", false
	} else {
		return value.(string), true
	}
}

func (m *MemoryStore) SetString(key string, value string) {
	m.Data[key] = value
}

func (m *MemoryStore) GetInt(key string) (int, bool) {
	if value, ok := m.Data[key]; !ok {
		return 0, false
	} else {
		return value.(int), true
	}
}

func (m *MemoryStore) SetInt(key string, value int) {
	m.Data[key] = value
}

func (m *MemoryStore) GetFloat(key string) (float64, bool) {
	if value, ok := m.Data[key]; !ok {
		return 0, false
	} else {
		return value.(float64), true
	}
}

func (m *MemoryStore) SetFloat(key string, value float64) {
	m.Data[key] = value
}
