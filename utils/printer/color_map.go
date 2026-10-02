package printer

import "github.com/fatih/color"

type ColorMap = map[string]*color.Color

func defaultColorMap() ColorMap {
	return ColorMap{
		"info":    color.New(color.FgWhite),
		"success": color.New(color.FgGreen),
		"warn":    color.New(color.FgYellow),
		"error":   color.New(color.FgRed),
	}
}

func FillMissingColorMapColors(colorMap ColorMap) {
	dcm := defaultColorMap()

	for k := range dcm {
		_, ok := colorMap[k]
		if !ok {
			colorMap[k] = dcm[k]
		}
	}
}
