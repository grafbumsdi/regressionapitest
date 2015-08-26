package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

const requestTimeout time.Duration = 300

func main() {
	// command line flags
	serveraddressPtr := flag.String("serveraddress", "", "the server address under test (e.g.: 192.168.11.23)")
	apiCallsPtr := flag.String("apicalls", "api/v1/wikifolios,api/v1/trades,api/v1/import/wikifolios", "comma separated list of api calls to test")
	logFilePtr := flag.String("logfile", "regressionapitest.log", "specify log file or set to 'stdout' to write to standard output")
	logLevelPtr := flag.String("loglevel", "Info", "set loglevel to 'Trace' to log API response")
	flag.Parse()
	
	initLog(*logFilePtr, *logLevelPtr)
	serverUrl := "http://" + getServerUrl(*serveraddressPtr)

	Info.Println("Creating requests for " + serverUrl)

	apiCalls := strings.Split(*apiCallsPtr, ",")

	for _, apiCall := range apiCalls {
		// WHEN: I call the api
		respbody := getResponseBody(serverUrl + "/" + apiCall)
		Info.Println("Checking API call: " + apiCall)
		// THEN: the response should not be null
		if(respbody == "null") {
			Error.Println("API call " + apiCall + " failed: the response was null")
		} else {
			Info.Println("API call " + apiCall + " succeeded: response was not 'null'")
		}
		// THEN: the response is a valid json
		if(!isJson(respbody)) {
			Error.Println("API call " + apiCall + " failed: the response was not a valid Json")
		} else {
			Info.Println("API call " + apiCall + " succeeded: response was a valid JSON")
		}
	}
}

func initLog(fileName string, loglevel string) {
	// open log stream
	var outputstream *os.File
	if strings.ToLower(fileName) == "stdout" {
		outputstream = os.Stdout
	} else {
		var err error
		outputstream, err = os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalln("Failed to open log file: " + fileName, os.Stdout, ":", err)
			outputstream = os.Stdout
		}
	}
	// initialize loggers
	if strings.ToLower(loglevel) == "trace" {
		Trace = log.New(outputstream,
			"TRACE:   ",
			log.Ldate|log.Ltime|log.Lshortfile)
	} else {
		Trace = log.New(ioutil.Discard,
			"TRACE:   ",
			log.Ldate|log.Ltime|log.Lshortfile)
	}

	Info = log.New(outputstream,
		"INFO:    ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(outputstream,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	multi := io.MultiWriter(outputstream, os.Stdout)
	Error = log.New(multi,
		"ERROR:   ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

func isJson(s string) bool {
	return isJsonObject(s) || isJsonString(s)
}

func isJsonString(s string) bool {
	var js string
	return json.Unmarshal([]byte(s), &js) == nil
}

func isJsonObject(s string) bool {
	var js interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}

func getResponseBody(url string) string {
	timeout := time.Duration(requestTimeout * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	Info.Println("Requesting following url: " + url)
	resp, err := client.Get(url)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		Error.Fatalln("Error while getting response for: " + url)
		panic(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		Error.Fatalln("Error while reading response body for: " + url)
		panic(err)
	}
	Trace.Println("Got following response body: \n" + string(body));
	return string(body)
}

func getServerUrl(serverUrl string) string {
	if "" == serverUrl {
		// no server address was given -> so we ask the user
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Please enter the web address you want to test (e.g.: 192.168.11.23): ")
		serverUrl, _ = reader.ReadString('\n')
	}
	return strings.TrimSpace(serverUrl)
}