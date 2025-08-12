package utils

import (
	"strings"
	"unicode/utf8"
)

// GetBigFont returns the input string in a "big font" style (mock implementation)
func GetBigFont(s string) string {
	// For demonstration, just return the string surrounded by asterisks
	return "*** " + s + " ***"
}

// GetWeatherIcon returns a weather icon based on the weather string
func GetWeatherIcon(weather string) string {
	switch weather {
	case "Sunny":
		return "☀️"
	case "Overcast":
		return "☁️"
	case "Raining":
		return "🌧️"
	case "Snow":
		return "❄️"
	default:
		return ""
	}
}

func GetWeatherAsciiArt(weather string) string {
	switch strings.ToLower(weather) {
	case "sunny":
		return `
   \   /  
	.-.   
 ― (   ) ―
	'-'   
   /   \  
`
	case "overcast":
		return `
	 .--.  
  .-(    ). 
 (___.__)__) 
`
	case "raining":
		return `
	 .-.    
	(   ).  
   (___(__) 
   ‚ʻ‚ʻ‚ʻ‚ʻ 
   ‚ʻ‚ʻ‚ʻ‚ʻ 
`
	case "snow":
		return `
	 .-.     
	(   ).   
   (___(__)  
   * * * *   
  * * * *    
`
	default:
		return `
	???     
`
	}
}

func AlignLeftRight(left, right string, width int) string {
	leftLen := utf8.RuneCountInString(left)
	rightLen := utf8.RuneCountInString(right)

	if leftLen+rightLen > width {
		if width > rightLen+3 {
			left = left[:width-rightLen-3] + "..."
			leftLen = utf8.RuneCountInString(left)
		} else {
			return right[:width]
		}
	}

	spaces := width - leftLen - rightLen
	if spaces < 1 {
		spaces = 1
	}

	return left + strings.Repeat(" ", spaces) + right
}
