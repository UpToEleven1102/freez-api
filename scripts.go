package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func main() {
	scriptName := ""
	if len(os.Args) > 1 {
		scriptName = os.Args[1]
	} else {
		fmt.Println("script accepts 1 argument")
	}

	switch scriptName {
	case "docs":
		cmd := exec.Command("swagger", "-apiPackage=git.nextgencode.io/huyen.vu/freez-app-rest", "-format=swagger", "-output=./docs" )
		err := cmd.Run()
		if err != nil {
			log.Fatalf("Failed to run script: %s \n", err)
		}
	}
}