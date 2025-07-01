package store

import (
	"fmt"
	"stormaaja/go-ha/data-store/tools"
	"sync"
	"time"
)

type MinerState struct {
	DeviceId            string
	IsMining            bool // Deprecated
	LastConfigChanged   time.Time
	SpotPriceLimit      float64
	TemperatureLimit    float64
	TemperatureSensorId string
}

type SaveableMinerStateStore struct {
	States map[string]MinerState
}

type MinerStateStore struct {
	MinerStates *sync.Map
}

func CreateMinerStateStore() MinerStateStore {
	store := MinerStateStore{
		MinerStates: new(sync.Map),
	}
	store.Load()
	return store
}

func (mss *MinerStateStore) Load() error {
	saveableState := SaveableMinerStateStore{
		States: make(map[string]MinerState),
	}
	err := saveableState.Load()
	if err != nil {
		return err
	}
	mss.Clear()
	for key, state := range saveableState.States {
		mss.MinerStates.Store(key, state)
	}
	return nil
}

func (smss *SaveableMinerStateStore) Load() error {
	return tools.ReadJsonFile("miner_state.json", smss)
}

func (mss MinerStateStore) ToSaveable() *SaveableMinerStateStore {
	states := make(map[string]MinerState)
	mss.MinerStates.Range(func(key, value interface{}) bool {
		if minerState, ok := value.(MinerState); ok {
			states[key.(string)] = minerState
		}
		return true
	})
	return &SaveableMinerStateStore{
		States: states,
	}
}

func (mss *MinerStateStore) Save() error {
	saveableState := mss.ToSaveable()
	err := tools.WriteJsonFile("miner_state.json", saveableState)
	if err != nil {
		return fmt.Errorf("failed to save miner state store: %w", err)
	}
	return nil
}

func (mss MinerStateStore) GetValue(key string) (MinerState, error) {
	value, exists := mss.MinerStates.Load(key)
	if !exists {
		return MinerState{}, fmt.Errorf("value not found")
	}
	if minerState, ok := value.(MinerState); ok {
		return minerState, nil
	}
	return MinerState{}, fmt.Errorf("value is not of type MinerState")
}

func (mss *MinerStateStore) SetValue(key string, value MinerState) {
	mss.MinerStates.Store(key, value)
	err := mss.Save()
	if err != nil {
		fmt.Printf("failed to save miner state store after setting value: %v\n", err)
	}
}

func (mss MinerStateStore) ContainsValue(key string) bool {
	_, exists := mss.MinerStates.Load(key)
	return exists
}

func (mss *MinerStateStore) DeleteValue(key string) {
	mss.MinerStates.Delete(key)
	err := mss.Save()
	if err != nil {
		fmt.Printf("failed to save miner state store after deleting value: %v\n", err)
	}
}

func (mss *MinerStateStore) Clear() {
	mss.MinerStates.Clear()
	err := mss.Save()
	if err != nil {
		fmt.Printf("failed to save miner state store after clearing: %v\n", err)
	}
}

func Keys(states *sync.Map) []string {
	keys := make([]string, 0)
	states.Range(func(key, _ interface{}) bool {
		keys = append(keys, key.(string))
		return true
	})
	return keys
}

func (mss MinerStateStore) GetIds() []string {
	return Keys(mss.MinerStates)
}
