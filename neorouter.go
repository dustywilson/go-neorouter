package neorouter


import (
	"os"
	"fmt"
	"exec"
	"bufio"
	"regexp"
)


type List struct {
	Computers []Computer
}

type Computer struct {
	Group    string
	Name     string
	Ip       string
	IsOnline bool
}


func GetList(domain, username, password string) (list List, err os.Error) {
	cmd := exec.Command("/usr/bin/nrclientcmd", "-p", password)
	cmdStdin, _ := cmd.StdinPipe()   // in  = into process from me
	cmdStdout, _ := cmd.StdoutPipe() // out = out from process to me
	reader := bufio.NewReader(cmdStdout)
	reDomain := regexp.MustCompile("Domain:")
	reUsername := regexp.MustCompile("Username:")
	reGroupname := regexp.MustCompile("> +([0-9A-Za-z]+)")
	reComputerOffline := regexp.MustCompile("\\(offline\\) +([0-9A-Za-z ]+)")
	reComputerOnline := regexp.MustCompile("([0-9\\.]+) +([0-9A-Za-z ]+)")
	reBadLogin := regexp.MustCompile("The system could not sign you in")
	cmd.Start()
	for {
		line, _ := reader.ReadString(':')
		if reDomain.Match([]byte(line)) {
			cmdStdin.Write([]byte(fmt.Sprintf("%s\n", domain)))
		} else if reUsername.Match([]byte(line)) {
			cmdStdin.Write([]byte(fmt.Sprintf("%s\n", username)))
			break
		} else {
			break
		}
	}
	groupname := "UNASSIGNED"
	for {
		line, _, err := reader.ReadLine()
		if err != nil { // probably EOF which is okay
			break
		}
		if reBadLogin.Match(line) {
			return list, os.NewError("Login Failed")
			break
		} else if match := reGroupname.FindStringSubmatch(string(line)); match != nil {
			groupname = match[1]
		} else if match := reComputerOffline.FindStringSubmatch(string(line)); match != nil {
			computername := match[1]
			computer := Computer{Group: groupname, Name: computername, IsOnline: false}
			list.Computers = append(list.Computers, computer)
		} else if match := reComputerOnline.FindStringSubmatch(string(line)); match != nil {
			ip := match[1]
			computername := match[2]
			computer := Computer{Group: groupname, Name: computername, Ip: ip, IsOnline: true}
			list.Computers = append(list.Computers, computer)
		}
		cmdStdin.Write([]byte("quit\n"))
	}
	cmd.Process.Kill()
	return list, nil
}
