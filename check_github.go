// check_github polls github and alerts the user when it comes back up.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	//"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"time"
)

const githubStatusEndpoint = "https://status.github.com/api/status.json"

type Status struct {
	Status    string `json="status"`
	Timestamp string `json="timestamp"`
}

var shouldSpeak = true

func speak(text string) (err error) {
	if !shouldSpeak {
		return
	}
	var cmd *exec.Cmd

	cmd = exec.Command("say", text)
	err = cmd.Run()
	return
}

func check() (status *Status) {
	fmt.Println("[+] checking status API")
	resp, err := http.Get(githubStatusEndpoint)
	if err != nil {
		fmt.Println("[!] fatal connect problem: ", err.Error())
		speak("HTTP request failed.")
		return
	}
	defer resp.Body.Close()
	status = new(Status)
	if err = json.NewDecoder(resp.Body).Decode(status); err != nil {
		fmt.Println("[!] error unmarshalling JSON: ", err.Error())
		speak("Could not unmarshal JSON response.")
		status = nil
	}
	return
}

func main() {
	waitStr := flag.String("t", "5m", "time.ParseDuration value")
	goodOnly := flag.Bool("g", false, "keep running until status is 'good'")
	quiet := flag.Bool("q", false, "don't speak status")
	once := flag.Bool("1", false, "only run one check")
	flag.Parse()
	shouldSpeak = !(*quiet)
	wait, err := time.ParseDuration(*waitStr)
	if err != nil {
		fmt.Println("could not parse wait time: ", err.Error())
		speak("Check Github failed to start.")
		os.Exit(1)
	}
	fmt.Println("[+] Github status check")
	speak("Git hub status check starting.")
	for {
		st := check()
		if st == nil {
			continue
		} else if st.Status == "good" {
			fmt.Println("[+] Github is operating normally.")
			speak("Git hub is operating normally.")
			break
		} else if st.Status == "minor" {
			fmt.Println("[+] some performance degradation.")
			speak("Git hub is experiencing degraded performance.")
			if !(*goodOnly) {
				break
			}
		} else {
			fmt.Printf("[+] Git hub is down. Major service interruptions.")
			speak("Git hub is experiencing major service interruptions.")
		}
		if *once {
			break
		}
		<-time.After(wait)
	}
}
