package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	ColorReset = "\033[0m"
	ColorRed   = "\033[31m"
	ColorGreen = "\033[32m"
	ColorBlue  = "\033[34m"
)

var reader = bufio.NewReader(os.Stdin)

func AskFor(message string, validValues ...string) string {
	if len(validValues) == 0 {
		fmt.Printf("%s » %s%s: ", ColorBlue, message, ColorReset)
	} else {
		fmt.Printf("%s » %s %v%s: ", ColorBlue, message, validValues, ColorReset)
	}

	line, err := reader.ReadString('\n')
	HandleError(err, "Cannot read user input")

	response := strings.TrimSuffix(line, "\n")
	if len(validValues) == 0 {
		return response
	}

	for _, value := range validValues {
		if strings.EqualFold(value, response) {
			return value
		}
	}
	HandleError(fmt.Errorf("Invalid response. Valid responses are %v.", validValues), "Invalid response")
	return ""
}
