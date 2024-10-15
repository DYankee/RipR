package Internal

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type osInfo struct {
	toName   string
	fromName string
	eol      string
}

type ClipInfo struct {
	Track int     `json:"track"`
	Start float64 `json:"start"`
	End   float64 `json:"end"`
	Color int     `json:"color"`
	Name  string  `json:"name"`
}

func (ci *ClipInfo) GetClipLength() float64 {
	return ci.End - ci.Start
}

type Connection struct {
	send    *os.File
	receive *os.File
}

type Audacity struct {
	osInfo     osInfo
	connection Connection
	Status     bool
}

func (a *Audacity) Init() {
	if runtime.GOOS == "windows" {
		fmt.Println("pipe-test.go, running on windows")
		a.osInfo.toName = `\\.\pipe\ToSrvPipe`
		a.osInfo.fromName = `\\.\pipe\FromSrvPipe`
		a.osInfo.eol = "\r\n"
	} else {
		fmt.Println("pipe-test.go, running on linux or mac")
		a.osInfo.toName = fmt.Sprintf("/tmp/audacity_script_pipe.to.%d", os.Getuid())
		a.osInfo.fromName = fmt.Sprintf("/tmp/audacity_script_pipe.from.%d", os.Getuid())
		a.osInfo.eol = "\n"
	}
}

// Establish a connection to audacity
func (a *Audacity) Open(fileName string) {
	c := exec.Command("audacity", fileName)
	err := c.Start()
	if err != nil {
		fmt.Println(err)
	}
}

func (a *Audacity) Connect() error {
	a.Status = false
	toFile, err := os.OpenFile(a.osInfo.toName, os.O_RDWR, os.ModeNamedPipe)
	if err != nil {
		log.Println("-- Failed to open file to write to:", err)
		a.Status = false
	} else {
		a.Status = true
		log.Println("-- File to write to has been opened")
	}
	a.connection.send = toFile

	fromfile, err := os.OpenFile(a.osInfo.fromName, os.O_RDWR, os.ModeNamedPipe)
	if err != nil {
		log.Println("-- Failed to open file to read from:", err)
		a.Status = false
	} else {
		a.Status = true
		log.Println("-- File to read from has been opened")
	}
	a.connection.receive = fromfile

	if !a.Status {
		err = errors.New("cannot connect to audacity")
	} else {
		err = nil
	}
	return err
}

func (a *Audacity) OpenFile(file string) {
	exec.Command("Audacity", file, "&")
}

func (a *Audacity) Close() {
	a.connection.receive.Close()
	a.connection.send.Close()
}

// send custom command to audacity. refer to https://manual.audacityteam.org/man/scripting_reference.html for formatting.
func (a *Audacity) Do_command(command string) (res string) {
	//send command
	fmt.Println("Send: >>> \n" + command)
	a.connection.send.Write([]byte(command + a.osInfo.eol))

	//get response
	scanner := bufio.NewScanner(a.connection.receive)
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

func (a *Audacity) SelectRegion(startTime float64, endTime float64) string {
	cmd := fmt.Sprintf(`Select: End="%f" RelativeTo="ProjectStart" Start="%f"`, endTime, startTime)
	res := a.Do_command(cmd)
	return res
}

func (a *Audacity) Split() string {
	res := a.Do_command("SplitNew:")
	return res
}

func (a *Audacity) SetLabel(labelId int, labelText string) string {
	cmd := fmt.Sprintf(`SetLabel: Label="%d" Text="%s"`, labelId, labelText)
	res := a.Do_command(cmd)
	return res
}

func (a *Audacity) ExportAudio(destination string, fileName string) string {
	cmd := fmt.Sprintf(`Export2: Filename="%s/%s" NumChannels="2"`, destination, fileName)
	res := a.Do_command(cmd)
	return res
}

func (a *Audacity) GetClips() []ClipInfo {
	info := []ClipInfo{}
	cmd := `GetInfo: Format="JSON" Type="Clips"`
	res := a.Do_command(cmd)
	log.Println(res)
	substrings := strings.SplitAfter(strings.Split(res, "[")[1], "]")
	res = strings.TrimRight(substrings[0], "]")
	substrings = strings.SplitAfter(res, "},")

	log.Println(res)
	for k, v := range substrings {
		info = append(info, ClipInfo{})
		v = strings.Trim(v, ",")
		err := json.Unmarshal([]byte(v), &info[k])
		if err != nil {
			println(err)
		}
	}
	log.Println("Get Info result")
	log.Println(info)
	return info
}
