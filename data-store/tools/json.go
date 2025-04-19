package tools

import (
	"encoding/json"
	"os"
)

func ReadJsonFile(filePath string, value any) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, value)
}

func WriteJsonFile(filePath string, value any) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, jsonData, 0644)
}
