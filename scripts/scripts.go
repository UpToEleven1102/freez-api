package scripts

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
)

func runScript(command string, args... string) {
	var stdOut, stdErr bytes.Buffer

	cmd := exec.Command(command, args...)
	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr

	err := cmd.Run()

	if err != nil {
		log.Fatalf("Failed to run script: %s - %s \n", err, stdErr.String())
	}

	fmt.Printf("%s \n %s \n", stdOut.String(), stdErr.String())
}

func RunScripts(script string) {
	switch script {
	case "docs":
		runScript("swagger", "-apiPackage=git.nextgencode.io/huyen.vu/freez-app-rest", "-format=swagger", "-output=./docs" )
	case "push-code":
		runScript("scp",  "-i", "/home/huyen/.ssh/Freeze.pem","-r", "/home/huyen/gospace/src/git.nextgencode.io/huyen.vu/freez-app-rest/", "ubuntu@35.162.158.187:/home/ubuntu/go/src/git.nextgencode.io/huyen.vu")
	case "ls":
		runScript("ls", "-lA")
	default:
		fmt.Println("No script found!")
	}
}