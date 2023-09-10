package world

import "image"

// HelpText lines must be 39 characters or less.
var HelpText = []string{`
Welcome to City Limits!          (1/2)
As the new mayor, it's time to run
things YOUR way. For better or worse...
Will you lead the clean energy front,
or will you put profits before people?
`, `
Moving Via Mouse                 (2/2)
To move around, press and hold your
middle mouse button while moving your
mouse, or press right click to center
the camera on an area immediately.
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
