package services

import (
	"fmt"
	"os"
)

func HandleErrorWithText(err error, text string) {
	if err != nil {
		fmt.Println(text)
		os.Exit(1)
	}
}

func HandleError(err error) {
	HandleErrorWithText(err, err.Error())
}
