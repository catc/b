package color

import "bytes"

// colors
const green = "\033[92m"
const white = "\033[39m"
const yellow = "\033[33m"
const blue = "\033[94m"
const magenta = "\033[95m"

func getColor(s string) string {
	switch s {
	case "green":
		return green
	case "yellow":
		return yellow
	case "blue":
		return blue
	case "magenta":
		return magenta
	}
	return white
}

// CustomColorFunc is a custom color func to avoid reset (and allow bold styling for entire line)
func CustomColorFunc(style string) func(string) string {
	color := getColor(style)
	return func(s string) string {
		if s == "" {
			return s
		}
		buf := bytes.NewBufferString(color)
		buf.WriteString(s)
		result := buf.String()
		return result
	}
}
