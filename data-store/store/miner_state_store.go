package store

import (
	"fmt"
	"stormaaja/go-ha/data-store/tools"
)

type MinerState struct {
	DeviceId string
	IsMining bool
}

type MinerStateStore struct {
	States map[string]MinerState
}

func CreateMinerStateStore() MinerStateStore {
	store := MinerStateStore{
		States: make(map[string]MinerState),
	}
	err := store.Load()
	if err != nil {
		fmt.Printf("failed to read miner state store file: %v\n", err)
	}
	return store
}

func (mss *MinerStateStore) Load() error {
	return tools.ReadJsonFile("miner_state.json", mss)
}

func (mss *MinerStateStore) Save() error {
	err := tools.WriteJsonFile("miner_state.json", mss)
	if err != nil {
		return fmt.Errorf("failed to save miner state store: %w", err)
	}
	return nil
}

func (mss *MinerStateStore) GetValue(key string) (MinerState, error) {
	value, exists := mss.States[key]
	if !exists {
		return MinerState{}, fmt.Errorf("value not found")
	}
	return value, nil
}

func (mss *MinerStateStore) SetValue(key string, value MinerState) {
	mss.States[key] = value
	err := mss.Save()
	if err != nil {
		fmt.Printf("failed to save miner state store after setting value: %v\n", err)
	}
}

func (mss *MinerStateStore) ContainsValue(key string) bool {
	_, exists := mss.States[key]
	return exists
}

func (mss *MinerStateStore) DeleteValue(key string) {
	delete(mss.States, key)
	err := mss.Save()
	if err != nil {
		fmt.Printf("failed to save miner state store after deleting value: %v\n", err)
	}
}

func (mss *MinerStateStore) Clear() {
	mss.States = make(map[string]MinerState)
	err := mss.Save()
	if err != nil {
		fmt.Printf("failed to save miner state store after clearing: %v\n", err)
	}
}

func (mss *MinerStateStore) GetIds() []string {
	ids := make([]string, 0, len(mss.States))
	for id := range mss.States {
		ids = append(ids, id)
	}
	return ids
}
