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
