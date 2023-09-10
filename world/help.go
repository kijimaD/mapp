package world

import "image"

// HelpText lines must be 39 characters or less.
var HelpText = []string{`
Welcome to Mapp!          (1/2)
`, `
This is last page...      (2/2)
`,
}

func HelpButtonAt(x, y int) int {
	point := image.Point{x, y}
	for i, rect := range World.HelpButtonRects {
		if point.In(rect) {
			return i
		}
	}
	return -1
}

func SetHelpPage(page int) {
	World.HelpPage = page
	World.HelpUpdated = true
	World.HUDUpdated = true
}
