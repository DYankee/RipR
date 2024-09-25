package Internal

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

var (
	TONAME   string
	FROMNAME string
	EOL      string
)

type TrackInfo struct {
	Name     string  `json:"name"`
	Focused  int     `json:"focused"`
	Selected int     `json:"selected"`
	Kind     string  `json:"kind"`
	Start    int     `json:"start"`
	End      float64 `json:"end"`
	Pan      int     `json:"pan"`
	Gain     int     `json:"gain"`
	Channels int     `json:"channels"`
	Solo     int     `json:"solo"`
	Mute     int     `json:"mute"`
	VZoomMin int     `json:"VZoomMin"`
	VZoomMax int     `json:"VZoomMax"`
}

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

type Connection struct {
	send    *os.File
	recieve *os.File
}

type Audacity struct {
	connection Connection
	Status     bool
}

// Establish a connection to audacity
func (a *Audacity) Open(fileName string) {
	c := exec.Command("audacity", fileName)
	err := c.Start()
	if err != nil {
		fmt.Println(err)
	}
}

func (a *Audacity) Connect() {
	for !a.Status {
		fmt.Println("Write to  \"" + TONAME + "\"")
		if _, err := os.Stat(TONAME); err != nil {
			fmt.Println(" ..does not exist.  Ensure Audacity is running with mod-script-pipe.")
			a.Status = false
		} else {
			a.Status = true
		}
		fmt.Println("Read from \"" + FROMNAME + "\"")
		if _, err := os.Stat(FROMNAME); err != nil {
			fmt.Println(" ..does not exist.  Ensure Audacity is running with mod-script-pipe.")
			a.Status = false
		} else {
			a.Status = true
		}
		time.Sleep(1 * time.Second)
	}
	fmt.Println("-- Both pipes exist.  Good.")
	toFile, err := os.OpenFile(TONAME, os.O_RDWR, os.ModeNamedPipe)
	if err != nil {
		fmt.Println("-- Failed to open file to write to:", err)
		os.Exit(1)
	}
	fmt.Println("-- File to write to has been opened")
	a.connection.send = toFile

	fromFile, err := os.OpenFile(FROMNAME, os.O_RDONLY, os.ModeNamedPipe)
	if err != nil {
		fmt.Println("-- Failed to open file to read from:", err)
		os.Exit(1)
	}
	fmt.Println("-- File to read from has now been opened too\r")
	a.connection.recieve = fromFile
}

func (a Audacity) OpenFile(file string) {
	exec.Command("Audacity", file, "&")
}

func (a Audacity) Close() {
	a.connection.recieve.Close()
	a.connection.send.Close()
}

// send custom command to audacity. reffer to https://manual.audacityteam.org/man/scripting_reference.html for formatting.
func (a Audacity) Do_command(command string) (res string) {
	//send command
	fmt.Println("Send: >>> \n" + command)
	a.connection.send.Write([]byte(command + EOL))

	//get response
	scanner := bufio.NewScanner(a.connection.recieve)
	for scanner.Scan() {
		text := scanner.Text()
		res += text
		if text == "" && len(res) != 0 {
			break
		}
	}
	fmt.Println("Received: <<< \n" + res)
	return res
}

func (a Audacity) SelectRegion(startTime float64, endTime float64) string {
	cmd := fmt.Sprintf("Select: End=\"%f\" RelativeTo=\"ProjectStart\" Start=\"%f\"", endTime, startTime)
	res := a.Do_command(cmd)
	return res
}

func (a Audacity) Split() string {
	res := a.Do_command("SplitNew:")
	return res

}

func (a Audacity) SetLabel(labelId int, labelText string) string {
	cmd := fmt.Sprintf(`SetLabel: Label="%d" Text="%s"`, labelId, labelText)
	res := a.Do_command(cmd)
	return res
}

func (a Audacity) ExportAudio(destination string, fileName string) string {
	cmd := fmt.Sprintf(`Export2: Filename="%s/%s" NumChannels="2"`, destination, fileName)
	res := a.Do_command(cmd)
	return res
}

func (a Audacity) GetInfo() TrackInfo {
	info := TrackInfo{}
	cmd := `GetInfo: Format "JSON" Type="Tracks"`
	res := a.Do_command(cmd)
	substrings := strings.SplitAfter(strings.Split(res, "[")[1], "]")
	res = strings.TrimRight(substrings[0], "]")
	json.Unmarshal([]byte(res), &info)
	fmt.Println(info)
	return info
}
