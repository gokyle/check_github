// check_github polls github and alerts the user when it comes back up.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
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
	var path string
	var cmd *exec.Cmd

	path, err = exec.LookPath("say")
	if err != nil {
		return
	}

	cmd = exec.Command(path, text)
	err = cmd.Run()
	return
}

func speaker(text string) (err error) {
	if shouldSpeak {
		err = speak(text)
	}
	return
}

func check() (status *Status) {
	fmt.Println("[+] checking status API")
	resp, err := http.Get(githubStatusEndpoint)
	if err != nil {
		fmt.Println("[!] fatal connect problem: ", err.Error())
		speaker("HTTP request failed.")
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("[!] bad response from API: ", err.Error())
		speaker("Bad response from status API.")
		return
	}
	status = new(Status)
	err = json.Unmarshal(body, status)
	if err != nil {
		fmt.Println("[!] error unmarshalling JSON: ", err.Error())
		speaker("Could not unmarshal JSON response.")
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
		speaker("Check Github failed to start.")
		os.Exit(1)
	}
	fmt.Println("[+] Github status check")
	speaker("Git hub status check starting.")
	for {
		st := check()
		if st == nil {
			continue
		} else if st.Status == "good" {
			fmt.Println("[+] Github is operating normally.")
			speaker("Git hub is operating normally.")
			break
		} else if st.Status == "minor" && !(*goodOnly) {
                        fmt.Println("[+] some performance degradation.")
                        speaker("Git hub is experiencing degraded performance.")
                        break
                } else {
			fmt.Printf("[+] Git hub is down: %s\n", st.Status)
			speaker("Git hub status: " + st.Status + ". Git hub is still down.")
		}
		if *once {
			break
		}
		<-time.After(wait)
	}
}
