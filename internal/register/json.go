package register

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

type RegisterList struct {
	Alias      string
	Path       string
	TotalTime  string
	LastOpened string
}

type DataForJson struct {
	Filter   []string       `json:"filter"`
	Register []RegisterList `json:"register"`
}

func CreateForJson() []byte {
	file, err := os.Create("data.json")
	if err != nil {
		log.Fatalf("create record file error: %s", err)
	}
	defer file.Close()

	data := DataForJson{}
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	_, err = file.Write(jsonData)
	if err != nil {
		log.Fatalf("write record file error: %s", err)
	}

	returnfile, err := os.ReadFile("data.json")
	if err != nil {
		log.Fatal(err)
	}
	return returnfile
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

	err = WriteForJson(data)
	if err != nil {
		return err
	}

	return nil
}

func ReadForJson(path ...string) (*DataForJson, error) {
	// defined variable
	var data DataForJson
	var file []byte
	var err error

	// get filter data
	if len(path) == 0 {
		file, err = os.ReadFile("data.json")
		if err != nil {
			if os.IsNotExist(err) {
				file = CreateForJson()
			} else {
				return &data, err
			}
		}
	} else {
		file, err = os.ReadFile(path[0])
		if err != nil {
			return &data, err
		}
	}

	err = json.Unmarshal(file, &data)
	if err != nil {
		return &data, err
	}

	return &data, nil
}

func WriteForJson(data *DataForJson) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
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
