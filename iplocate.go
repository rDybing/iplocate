/********************************
  iplocate.go

  License: MIT

  Copyright (c) 2019 Roy Dybing

  github   : rDybing
  Linked In: Roy Dybing
  MeWe     : Roy Dybing

  Full license text in README.md
*********************************/

package main

import (
	"bufio"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

const apiFile = "./settings/api.json"
const conFile = "./settings/config.json"
const version = "0.1.1"
const divider = "------------------------------------------------------------------"

var client = &http.Client{
	Timeout: time.Second * 15,
}

type apiT struct {
	IP   string
	Key  string
	Call string
}

type ipDetailsT struct {
	Date    time.Time
	IP      string
	Method  string
	Country string `json:"country_name"`
	Region  string `json:"region_name"`
	City    string `json:"city_name"`
}

type logT struct {
	auth string
	f2b  string
	hist string
}

func main() {
	fmt.Printf("iplocate v%v\n", version)

	var logLoc logT
	api, logLoc, interval := initApp()

	if interval < 1 {
		log.Panic("Polling Interval set too low, exiting...")
	}

	ips := logLoc.importLogs(api)

	fmt.Println(divider)
	for i := range ips {
		date := fmt.Sprint(ips[i].Date.String())
		date = fmt.Sprint(ips[i].Date.Format("2006-01-02 15:04:05"))
		fmt.Printf("%2d: %-20s - %-20s - %s\n", i, date, ips[i].IP, ips[i].Method)
		fmt.Printf("    %-20s - %-20s - %-20s\n", ips[i].Country, ips[i].Region, ips[i].City)
		fmt.Println(divider)
	}
}

func initApp() (apiT, logT, int) {
	var api apiT
	var logFile logT

	if err := api.loadAPI(); err != nil {
		log.Panicf("Could not load or decode API credentials - exiting...\n%v\n", err)
	}
	interval, err := logFile.loadConfig()
	if err != nil {
		log.Panicf("Could not load or decode configuration file - exiting...\n%v\n", err)
	}

	return api, logFile, interval
}

func (l logT) importLogs(api apiT) []ipDetailsT {
	var ips []ipDetailsT

	ipMapOld := l.loadHistoryLog()
	for i := range ipMapOld {
		ipLoc := ipMapOld[i]
		ips = append(ips, ipLoc)
	}
	ipsHashOld := hashIps(ips)

	ipMapF2B := l.loadFail2BanLog()
	for i := range ipMapF2B {
		if _, found := ipMapOld[ipMapF2B[i].IP]; !found {
			ipLoc := api.getIPLocation(ipMapF2B[i].IP)
			ipLoc.Method = ipMapF2B[i].Method
			ipLoc.Date = ipMapF2B[i].Date
			ips = append(ips, ipLoc)
		}
	}
	ipsHashNew := hashIps(ips)

	if ipsHashOld != ipsHashNew {
		fmt.Println("Saving New History File")
		l.saveHistoryLog(ips)
	}

	sort.Slice(ips, func(i, j int) bool {
		return ips[i].Date.After(ips[j].Date)
	})
	return ips
}

func (l *logT) loadConfig() (int, error) {
	type conDataT struct {
		LogDir   string
		AuthLog  string
		F2bLog   string
		History  string
		Interval int
	}

	var conData conDataT

	f, err := os.Open(conFile)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	temp := json.NewDecoder(f)
	if err := temp.Decode(&conData); err != nil {
		return 0, err
	}

	l.auth = conData.LogDir + conData.AuthLog
	l.f2b = conData.LogDir + conData.F2bLog
	l.hist = conData.History
	return conData.Interval, nil
}

func (l logT) loadHistoryLog() map[string]ipDetailsT {
	var outJSON []ipDetailsT
	var out map[string]ipDetailsT
	out = make(map[string]ipDetailsT)

	if _, err := os.Stat(l.hist); err == nil {
		f, err := os.Open(l.hist)
		if err != nil {
			log.Printf("Could not open History File...\n%v\n", err)
			return out
		}
		defer f.Close()

		temp := json.NewDecoder(f)
		if err := temp.Decode(&outJSON); err != nil {
			log.Printf("Could not decode History file, skipping...\n%v\n", err)
			return out
		}

		for i := range outJSON {
			out[outJSON[i].IP] = outJSON[i]
		}
		return out
	}
	log.Printf("History File do not Exist - skipping...")
	return out
}

func (l logT) saveHistoryLog(in []ipDetailsT) error {
	f, err := os.Create(l.hist)
	if err != nil {
		return err
	}
	defer f.Close()
	tempJSON := json.NewEncoder(f)
	tempJSON.SetIndent("", "    ")
	if err := tempJSON.Encode(in); err != nil {
		return err
	}
	return nil
}

func (l logT) loadFail2BanLog() map[string]ipDetailsT {
	var out map[string]ipDetailsT
	out = make(map[string]ipDetailsT)

	f, err := os.Open(l.f2b)
	if err != nil {
		log.Panicf("Could not open Fail2Ban log - exiting...\n%v\n", err)
		return out
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "WARNING") && strings.Contains(line, "Ban") {
			lines = append(lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Could not scan lines of Fail2Ban log...\n%v\n", err)
	}

	var tempOut ipDetailsT
	layout := "2006-01-02 15:04:05"
	for i := range lines {
		line := strings.Split(lines[i], " ")
		subDate := strings.Split(line[1], ",")
		date := line[0] + " " + subDate[0]
		ip := line[6]
		if _, found := out[ip]; !found {
			tempOut.IP = ip
			if t, err := time.Parse(layout, date); err != nil {
				tempOut.Date = time.Now()
			} else {
				tempOut.Date = t
			}
			tempOut.Method = "Ban"
			out[ip] = tempOut
		}
	}
	return out
}

func hashIps(in []ipDetailsT) string {
	ips := ""
	for i := range in {
		ips += in[i].IP
	}
	data := []byte(ips)
	out := fmt.Sprintf("%x", md5.Sum(data))
	return out
}

func (a *apiT) loadAPI() error {
	f, err := os.Open(apiFile)
	if err != nil {
		return err
	}
	defer f.Close()
	temp := json.NewDecoder(f)
	if err := temp.Decode(&a); err != nil {
		return err
	}
	return nil
}

func (a apiT) getIPLocation(ip string) ipDetailsT {
	var out ipDetailsT
	out.IP = ip

	url := fmt.Sprintf("%s?ip=%s&key=%s&package=WS3", a.IP, ip, a.Key)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("ERROR: API Call Failed Entirely!\n%v\n", err)
		return out
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("ERROR: Could not unpack http.Response body!\n%v\n", err)
		return out
	}
	errJSON := json.Unmarshal(body, &out)
	if errJSON != nil {
		log.Printf("ERROR: JSON decode Failed!\n%v\n", err)
		return out
	}

	return out
}
