package main

import (
	"fmt"
	"log"
	"strings"

	"example.com/mod/classes"
	"example.com/mod/utils"
	"github.com/jroimartin/gocui"
)

// Store last known terminal size for resize detection
var lastX, lastY int

var favoritesMock = []classes.Weather{
	{City: "New York", Temp: "22°C", Weather: "Sunny"},
	{City: "Los Angeles", Temp: "28°C", Weather: "Overcast"},
	{City: "Chicago", Temp: "16°C", Weather: "Raining"},
	{City: "Denver", Temp: "10°C", Weather: "Snow"},
	{City: "Miami", Temp: "30°C", Weather: "Sunny"},
	{City: "Seattle", Temp: "18°C", Weather: "Overcast"},
	{City: "Boston", Temp: "20°C", Weather: "Raining"},
	{City: "San Francisco", Temp: "19°C", Weather: "Sunny"},
	{City: "Austin", Temp: "25°C", Weather: "Sunny"},
	{City: "Phoenix", Temp: "35°C", Weather: "Sunny"},
}

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Cursor = true
	g.Mouse = true
	g.Highlight = true
	g.SelFgColor = gocui.ColorWhite
	g.SelBgColor = gocui.ColorBlack
	g.InputEsc = true

	// Set layout manager
	g.SetManagerFunc(layout)

	// Quit keybinding
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func layout(g *gocui.Gui) error {
	selectedCity := classes.Weather{
		City:    "New York",
		Temp:    "22°C",
		Weather: "Sunny",
	}

	maxX, maxY := g.Size()

	// Detect resize
	if w, h := g.Size(); w != lastX || h != lastY {
		lastX, lastY = w, h
		onResize(g, w, h)
	}

	// Column widths
	leftW := maxX / 5
	centerW := maxX * 3 / 5

	// Left column split
	searchH := 3

	// Left: Search
	if v, err := g.SetView("search", 0, 0, leftW-1, searchH-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Search"
		v.Wrap = true
		v.Highlight = true
		v.Write([]byte("Enter location to search..."))
	}

	// Left: Favorites
	if v, err := g.SetView("favorites", 0, searchH, leftW-1, maxY-2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Favorites"
		v.Autoscroll = true
		v.Editor = gocui.DefaultEditor
		v.Editable = true
		g.SetKeybinding("favorites", gocui.MouseLeft, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
			_, cy := v.Cursor()
			lines := strings.Split(string(v.Buffer()), "\n")
			if cy >= 0 && cy < len(lines) {
				selected := lines[cy]
				fields := strings.Fields(selected)
				if len(fields) >= 2 {
					// Simulate updating current weather for clicked city
					// updateCurrentWeather(g, city, "N/A", "N/A", leftW+centerW-1)

					selectedCity = favoritesMock[cy]
					updateCurrentWeather(g, selectedCity.City, selectedCity.Temp, selectedCity.Weather, centerW-1)
				}
			}
			return nil
		})
		v.Write(listFavoritesMock(leftW - 2))
		v.Highlight = true
	}

	// Center: Current
	centerH := maxY * 4 / 6
	if v, err := g.SetView("current", leftW, 0, leftW+centerW-1, centerH-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Current"
		updateCurrentWeather(g, selectedCity.City, selectedCity.Temp, selectedCity.Weather, centerW-1)
	}

	// Center: Hour by Hour
	if v, err := g.SetView("Hour by Hour", leftW, centerH, leftW+centerW-1, maxY-2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Hour by Hour"
	}

	// Right: 10 Day Forecast
	if v, err := g.SetView("10day", leftW+centerW, 0, maxX-1, maxY-2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "10 Day Forecast"
	}

	// Help view at the bottom
	if v, err := g.SetView("help", 0, maxY-2, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = false
		v.BgColor = gocui.ColorWhite
		v.FgColor = gocui.ColorBlack
		fmt.Fprint(v, " Ctrl+C: Quit | Tab: Next | Shift+Tab: Previous ")
	}

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func listFavoritesMock(availableWidth int) []byte {

	var lines []string
	for _, fav := range favoritesMock {
		lines = append(lines,
			// fmt.Sprintf("%s\033[1m%-15s %-6s\033 ",
			// 	utils.GetWeatherIcon(fav.Weather),
			// 	fav.City,
			// 	fav.Temp))
			utils.AlignLeftRight(fav.City, utils.GetWeatherIcon(fav.Weather)+" "+fav.Weather, availableWidth),
		)

	}

	return []byte(strings.Join(lines, "\n"))
}

func updateCurrentWeather(g *gocui.Gui, city string, temp string, weather string, availableWidth int) error {
	v, err := g.View("current")
	if err != nil {
		return err
	}
	v.Clear()

	// Example aligned header
	fmt.Fprintln(v, utils.AlignLeftRight(city, weather, availableWidth-1))

	return nil
}

func onResize(g *gocui.Gui, w, h int) {
	centerW := w * 3 / 5
	leftW := w / 5

	// Update the "current" view with new width
	updateCurrentWeather(g, "New York", "22°C", "Sunny", centerW-1)
	listFavoritesMock(leftW - 2)
}
