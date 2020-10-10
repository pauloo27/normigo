package utils

import "fmt"

func HandleError(err error, msg string) {
	if err != nil {
		fmt.Printf("Error: %s\n", msg)
		panic(err)
	}
}
