package config

var backendUrl string

type terminal struct {
	width  int
	height int
}

var t terminal

func Init() {
	if backendUrl == "" {
		backendUrl = "http://localhost:3000/api/v1"
	}
}

func BackendUrl() string {
	return backendUrl
}

func TermWidth() int {
	return t.width
}

func TermHeight() int {
	return t.height
}

func TermSize() (width int, height int) {
	width = t.width
	height = t.height
	return
}

func UpdateTermSize(width int, height int) {
	t.width = width
	t.height = height
}
