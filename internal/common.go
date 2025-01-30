package internal

type ProcessList struct {
	Pid  int32
	Name string
	Path string
}

// gui component list
type GuiComponent struct {
	Name string
	Item interface{}
}
