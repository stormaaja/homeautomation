package store

import (
	"fmt"
	"stormaaja/go-ha/data-store/tools"
)

type GenericStore struct {
	FilePath string
	Values   map[string]any
}

func CreateGenericStore(filePath string) GenericStore {
	store := GenericStore{
		FilePath: filePath,
		Values:   make(map[string]any),
	}
	if filePath != "" {
		err := store.Load()
		if err != nil {
			fmt.Printf("failed to read store file: %v\n", err)
		}
	}
	return store
}

func (gs *GenericStore) Load() error {
	return tools.ReadJsonFile(gs.FilePath, &gs.Values)
}

func (gs *GenericStore) Save() error {
	err := tools.WriteJsonFile(gs.FilePath, gs.Values)
	if err != nil {
		return fmt.Errorf("failed to save store: %w", err)
	}
	return nil
}

func (gs *GenericStore) GetValue(key string) (any, error) {
	value, exists := gs.Values[key]
	if !exists {
		return nil, fmt.Errorf("value not found")
	}
	return value, nil
}

func (gs *GenericStore) SetValue(key string, value any) {
	gs.Values[key] = value
	err := gs.Save()
	if err != nil {
		fmt.Printf("failed to save store after setting value: %v\n", err)
	}
}

func (gs *GenericStore) ContainsValue(key string) bool {
	_, exists := gs.Values[key]
	return exists
}

func (gs *GenericStore) DeleteValue(key string) {
	delete(gs.Values, key)
	err := gs.Save()
	if err != nil {
		fmt.Printf("failed to save store after deleting value: %v\n", err)
	}
}

func (gs *GenericStore) GetIds() []string {
	ids := make([]string, 0, len(gs.Values))
	for id := range gs.Values {
		ids = append(ids, id)
	}
	return ids
}

func (gs *GenericStore) Clear() {
	gs.Values = make(map[string]any)
	err := gs.Save()
	if err != nil {
		fmt.Printf("failed to save store after clearing: %v\n", err)
	}
}
