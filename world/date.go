package world

import "strconv"

var monthNames = []string{
	"January",
	"February",
	"March",
	"April",
	"May",
	"June",
	"July",
	"August",
	"September",
	"October",
	"November",
	"December",
}

func Date() (month string, year string) {
	y, m := World.Ticks/YearTicks, (World.Ticks%YearTicks)/MonthTicks
	return monthNames[m], strconv.Itoa(startingYear + y)
}
