package internal

type ProcessList struct {
	Pid  int32
	Name string
	Path string
}

type GuiComponent struct {
	Name string
	Item interface{}
}

type RegisterList struct {
	Path       string
	TotalTime  string
	LastOpened string
}

type DataForJson struct {
	Filter   []string       `json:"filter"`
	Register []RegisterList `json:"register"`
}
