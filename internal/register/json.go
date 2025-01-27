package register

import (
	"encoding/json"
	"os"
	"time"

	"processchronicle/internal"
)

var data internal.DataForJson

func RegisterForJson(path string) error {
	// get filter data
	fileRead, err := os.ReadFile("data.json")
	if err != nil {
		return err
	}
	err = json.Unmarshal(fileRead, &data)
	if err != nil {
		return err
	}

	now := time.Now()
	data.Register = append(data.Register, internal.RegisterList{
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
