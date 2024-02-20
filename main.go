package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
)

var (
	TONAME   string
	FROMNAME string
	EOL      string
)

func init() {
	if runtime.GOOS == "windows" {
		fmt.Println("pipe-test.go, running on windows")
		TONAME = `\\.\pipe\ToSrvPipe`
		FROMNAME = `\\.\pipe\FromSrvPipe`
		EOL = "\r\n"
	} else {
		fmt.Println("pipe-test.go, running on linux or mac")
		TONAME = fmt.Sprintf("/tmp/audacity_script_pipe.to.%d", os.Getuid())
		FROMNAME = fmt.Sprintf("/tmp/audacity_script_pipe.from.%d", os.Getuid())
		EOL = "\n"
	}
}

type AudPipe struct {
	send    *os.File
	recieve *os.File
}

func openPipe() (pipe AudPipe) {
	pipe = AudPipe{}
	pipe.send, pipe.recieve = openFiles()
	return pipe
}

func checkConnection() {
	fmt.Println("Write to  \"" + TONAME + "\"")
	if _, err := os.Stat(TONAME); err != nil {
		fmt.Println(" ..does not exist.  Ensure Audacity is running with mod-script-pipe.")
		os.Exit(1)
	}

	fmt.Println("Read from \"" + FROMNAME + "\"")
	if _, err := os.Stat(FROMNAME); err != nil {
		fmt.Println(" ..does not exist.  Ensure Audacity is running with mod-script-pipe.")
		os.Exit(1)
	}
	fmt.Println("-- Both pipes exist.  Good.")
}

func openFiles() (*os.File, *os.File) {
	toFile, err := os.OpenFile(TONAME, os.O_RDWR, os.ModeNamedPipe)
	if err != nil {
		fmt.Println("-- Failed to open file to write to:", err)
		os.Exit(1)
	}

	fmt.Println("-- File to write to has been opened")

	fromFile, err := os.OpenFile(FROMNAME, os.O_RDONLY, os.ModeNamedPipe)
	if err != nil {
		fmt.Println("-- Failed to open file to read from:", err)
		os.Exit(1)
	}

	fmt.Println("-- File to read from has now been opened too\r")

	return toFile, fromFile
}

func sendCommand(toFile *os.File, command string) {
	fmt.Println("Send: >>> \n" + command)
	toFile.Write([]byte(command + EOL))
}

func getResponse(fromFile *os.File) string {
	scanner := bufio.NewScanner(fromFile)
	res := ""
	for scanner.Scan() {
		text := scanner.Text()
		res += text
		if text == "" && len(res) != 0 {
			break
		}
	}
	return res
}

func getres(fromFile *os.File) {
	reader := bufio.NewReader(fromFile)
	line, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(line)
}

type Connection struct {
	send    *os.File
	recieve *os.File
}

type Audacity struct {
	connection Connection
}

func (Audacity) connect() {
	fmt.Println("Write to  \"" + TONAME + "\"")
	if _, err := os.Stat(TONAME); err != nil {
		fmt.Println(" ..does not exist.  Ensure Audacity is running with mod-script-pipe.")
		os.Exit(1)
	}

	fmt.Println("Read from \"" + FROMNAME + "\"")
	if _, err := os.Stat(FROMNAME); err != nil {
		fmt.Println(" ..does not exist.  Ensure Audacity is running with mod-script-pipe.")
		os.Exit(1)
	}
	fmt.Println("-- Both pipes exist.  Good.")

	Audacity.connection.send, err = os.OpenFile(TONAME, os.O_RDWR, os.ModeNamedPipe)
	if err != nil {
		fmt.Println("-- Failed to open file to write to:", err)
		os.Exit(1)
	}

	fmt.Println("-- File to write to has been opened")

	fromFile, err := os.OpenFile(FROMNAME, os.O_RDONLY, os.ModeNamedPipe)
	if err != nil {
		fmt.Println("-- Failed to open file to read from:", err)
		os.Exit(1)
	}

	fmt.Println("-- File to read from has now been opened too\r")

	return toFile, fromFile
}

func main() {
	audacity := Audacity{}
	audacity.connect()
	pipe := openPipe()

	sendCommand(pipe.toFile, "Help: ")
	res := getResponse(fromFile)
	fmt.Println(res)
	sendCommand(toFile, "Help: ")
	getResponse(fromFile)

	toFile.Close()
	fromFile.Close()
}
