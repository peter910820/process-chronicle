package register

import (
	"encoding/json"
	"os"
	"time"
)

type RegisterList struct {
	Path       string
	TotalTime  string
	LastOpened string
}

type DataForJson struct {
	Filter   []string       `json:"filter"`
	Register []RegisterList `json:"register"`
}

func RegisterForJson(path string) error {
	data, err := ReadForJson()
	if err != nil {
		return err
	}

	now := time.Now()
	data.Register = append(data.Register, RegisterList{
		Path:       path,
		LastOpened: now.Format("2006-01-02 15:04:05"),
	})

	jsonData, err := json.MarshalIndent(&data, "", "  ")
	if err != nil {
		return err
	}

	fileWrite, err := os.OpenFile("data.json", os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer fileWrite.Close()

	_, err = fileWrite.Write(jsonData)
	if err != nil {
		return err
	}

	return nil
}

func ReadForJson(path ...string) (*DataForJson, error) {
	var data DataForJson
	// get filter data
	if len(path) == 0 {
		fileRead, err := os.ReadFile("data.json")
		if err != nil {
			return &data, err
		}
		err = json.Unmarshal(fileRead, &data)
		if err != nil {
			return &data, err
		}

	} else {
		fileRead, err := os.ReadFile(path[0])
		if err != nil {
			return &data, err
		}
		err = json.Unmarshal(fileRead, &data)
		if err != nil {
			return &data, err
		}
	}
	return &data, nil
}
