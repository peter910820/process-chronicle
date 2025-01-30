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

// the channel to send subprocess(timer) data when mainprocess close
type ChannelForJson struct {
	Alias   string
	Counter int
}

// create json record file when the json record file is not exist
func CreateForJson(path string) []byte {
	data := DataForJson{}

	file, err := os.Create(path)
	if err != nil {
		log.Fatalf("create record file error: %s", err)
	}
	defer file.Close()

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	_, err = file.Write(jsonData)
	if err != nil {
		log.Fatalf("write record file error: %s", err)
	}

	returnfile, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return returnfile
}

// register process into json record file
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
		TotalTime:  "0",
		LastOpened: now.Format("2006-01-02 15:04:05"),
	})

	err = WriteForJson(data)
	if err != nil {
		return err
	}

	return nil
}

// read the json record file
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
			file = CreateForJson(filePath)
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

// write the json record file
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
