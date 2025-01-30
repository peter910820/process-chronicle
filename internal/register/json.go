package register

import (
	"encoding/json"
	"log"
	"os"
	"strings"
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

type ChannelForJson struct {
	Alias   string
	Counter int
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
	pathSlice := strings.Split(path, "\\")
	data.Register = append(data.Register, RegisterList{
		Alias:      pathSlice[len(pathSlice)-1],
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
	data := DataForJson{}
	filePath := "data.json"

	if len(path) != 0 {
		filePath = path[0]
	}
	// get filter data
	file, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			file = CreateForJson()
		} else {
			return &data, err
		}
	}

	err = json.Unmarshal(file, &data)
	if err != nil {
		return &data, err
	}

	return &data, nil
}

func WriteForJson(data *DataForJson, path ...string) error {
	filePath := "data.json"
	if len(path) != 0 {
		filePath = path[0]
	}
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	fileWrite, err := os.OpenFile(filePath, os.O_RDWR|os.O_TRUNC, 0666)
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
