package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"
	"time"
)

var (
	pingInterval  string
	remoteApp     string
	gorouter      string
	postBodySize  = 10
	sleepDuration time.Duration
	letterRunes   = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

func init() {
	var err error
	pingInterval = os.Getenv("INTERVAL")
	if pingInterval == "" {
		panic("Please set env variable INTERVAL to value supported by https://pkg.go.dev/time#ParseDuration")
	}

	sleepDuration, err = time.ParseDuration(pingInterval)
	if err != nil {
		panic("Failed to parse INTERVAL")
	}

	remoteApp = os.Getenv("REMOTE_APP_HOSTNAME")
	if remoteApp == "" {
		panic("Please set env variable REMOTE_APP_HOSTNAME to the domain name of app you want to ping")
	}

	gorouter = os.Getenv("GOROUTER_ADDRESS")
	if remoteApp == "" {
		panic("Please set env variable GOROUTER_ADDRESS to be an IP of one gorouter")
	}
}

func debug(data []byte, err error) error {
	if err == nil {
		fmt.Printf("%s\n\n", data)
	} else {
		return err
	}
	return nil
}

func doRequest(r *http.Request) error {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Do(r)
	if err != nil {
		return fmt.Errorf("An Error Occured sending POST %v", err)
	}
	defer resp.Body.Close()

	return debug(httputil.DumpResponse(resp, true))
}

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func sendPost(server, hostheader string) error {
	// postBody, _ := json.Marshal(map[string]string{
	// 	"name":  "Toby",
	// 	"email": "Toby@example.com",
	// })

	pb := fmt.Sprintf("{\"data\": \"")
	for i := 1; i < postBodySize; i++ {
		pb += fmt.Sprintf("%s", RandStringRunes(1))
	}
	pb += fmt.Sprintf("\"}")

	postBody, _ := json.Marshal(pb)
	body := bytes.NewBuffer(postBody)

	req, err := http.NewRequest("POST", fmt.Sprintf("https://%s", server), body)
	if err != nil {
		fmt.Printf("Failed to send request using server: %s and host header: %s: %s\n", server, hostheader, err)
	}
	req.Host = hostheader
	req.Header.Add("Host", hostheader)
	return doRequest(req)
}

func ping() {
	for {
		fmt.Printf("Sending POST to %s\n", remoteApp)
		err := sendPost(remoteApp, remoteApp)
		if err != nil {
			fmt.Printf("Failed sending post: %s\n", err)
		}
		fmt.Printf("Sending POST to %s\n", gorouter)
		err = sendPost(gorouter, remoteApp)
		if err != nil {
			fmt.Printf("Failed sending post: %s\n", err)
		}
		time.Sleep(sleepDuration)
	}
}

func configHandler(w http.ResponseWriter, r *http.Request) {
	interval := r.FormValue("interval")
	if interval != "" {
		pingInterval = interval

		testDuration, err := time.ParseDuration(pingInterval)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Printf("Failed to interval duration.. see values supported in https://pkg.go.dev/time#ParseDuration: %s", err)
			return
		}
		sleepDuration = testDuration
	}
	host := r.FormValue("remoteapp")
	if host != "" {
		remoteApp = host
	}
	pbs := r.FormValue("postbodysize")
	if pbs != "" {
		i, err := strconv.Atoi(pbs)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Printf("Failed to parse postbodysize: %s", err)
			return
		}
		postBodySize = i
	}

	g := r.FormValue("gorouter")
	if g != "" {
		gorouter = g
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}
	log.Printf("%s\n", body)
	w.Write([]byte("<html>Hello!</html>"))
}

func main() {
	go ping()
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/config", configHandler)
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
