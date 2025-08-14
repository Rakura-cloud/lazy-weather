package main

import (
	"example.com/mod/classes"
	"example.com/mod/fetchers"
	"example.com/mod/utils"
	"fmt"
	"github.com/jroimartin/gocui"
	"log"
	"strings"
)

// Store last known terminal size for resize detection
var lastX, lastY int

var favoritesMock = []classes.Weather{}

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
	selectedCity := fetchers.GetWeatherFromLatLig(48.1486, 17.1077) // Default city

	maxX, maxY := g.Size()

	// Detect resize
	if w, h := g.Size(); w != lastX || h != lastY {
		lastX, lastY = w, h
		onResize(g, w, h, selectedCity)
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
	if v, err := g.SetView("favorites", 0, searchH, leftW-1, maxY-3); err != nil {
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

					updateCurrentWeather(g, selectedCity.Name, fmt.Sprintf("%.1f", selectedCity.Main.Temp), selectedCity.Weather[len(selectedCity.Weather)-1].Main, centerW-1)
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
		updateCurrentWeather(g, selectedCity.Name, fmt.Sprintf("%.1f", selectedCity.Main.Temp), selectedCity.Weather[len(selectedCity.Weather)-1].Main, centerW-1)
	}

	// Center: Hour by Hour
	if v, err := g.SetView("Hour by Hour", leftW, centerH, leftW+centerW-1, maxY-3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Hour by Hour"
	}

	// Right: 10 Day Forecast
	if v, err := g.SetView("10day", leftW+centerW, 0, maxX-1, maxY-3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "10 Day Forecast"
	}

	// Help view at the bottom
	if v, err := g.SetView("help", 0, maxY-3, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = false
		v.Wrap = true
		//	fmt.Fprint(v, " Ctrl+C: Quit | Tab: Next | Shift+Tab: Previous ")
		v.Write([]byte("Search: s | Refresh: t | Units: u | Quit: q"))
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
			utils.AlignLeftRight(
				fav.Name,
				func() string {
					if len(fav.Weather) > 0 {
						return utils.GetWeatherIcon(fav.Name) + " " + fav.Weather[0].Main
					}
					return ""
				}(),
				availableWidth,
			),
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

func onResize(g *gocui.Gui, w, h int, weather *classes.Weather) {
	centerW := w * 3 / 5
	leftW := w / 5

	// Update the "current" view with new width
	updateCurrentWeather(g, weather.Name, "", "", centerW-1)

	updateCurrentWeather(g, weather.Name, fmt.Sprintf("%.1f", weather.Main.Temp), weather.Weather[len(weather.Weather)-1].Main, centerW-1)

	listFavoritesMock(leftW - 2)
}
