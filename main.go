package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/brian1917/illumioapi"
)

//data structure built from the parser.csv
type regex struct {
	regexdata []regexstruct
}

//regex structure with regex and array of replace regex to build the labels
type regexstruct struct {
	regex   string
	labelcg map[string]string
}

func ReadCSV(file string) [][]string {
	csvfile, err := os.Open(file)
	defer csvfile.Close()

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	reader := csv.NewReader(csvfile)
	reader.FieldsPerRecord = -1

	rawCSVdata, err := reader.ReadAll()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	return rawCSVdata
}

//
// Regex method to provide labels for the hostname provided
func (r *regex) LabelFromHostname(hostname string) []string {

	var templabels []string

	found := false
	for _, tmp := range r.regexdata {

		//If you found a match skip to the next hostname
		if !found {
			//Place match regex into regexp data struct
			tmpre := regexp.MustCompile(tmp.regex)
			if tmpre.MatchString(hostname) {
				found = true

				for _, label := range []string{"loc", "env", "app", "role"} {
					//Build Array of based on the specific replace regex for each label
					templabels = append(templabels, tmpre.ReplaceAllString(hostname, tmp.labelcg[label]))

				}
			}
		}
	}
	return templabels
}

//Load the Regex CSV Into the parser struct -
func (reg *regex) load(data [][]string) {

	//Cycle through all the parse data rows in the parse data xls
	for c, row := range data {

		var r regexstruct
		//ignore header
		if c != 0 {

			//Array order 0-LOC,1-ENV,2-APP,3-APP
			tmpmap := make(map[string]string)
			for x, lbl := range []string{"loc", "env", "app", "role"} {
				//place CSV column in map
				tmpmap[lbl] = row[x+1]
			}
			//Put the regex string and capture groups into data structure
			r.regex = row[0]
			r.labelcg = tmpmap

			reg.regexdata = append(reg.regexdata, r)
		}

	}
}

func main() {

	config, _ := parseConfig()

	if len(config.Logging.LogDirectory) > 0 && config.Logging.LogDirectory[len(config.Logging.LogDirectory)-1:] != string(os.PathSeparator) {
		config.Logging.LogDirectory = config.Logging.LogDirectory + string(os.PathSeparator)
	}
	var logfile string
	if config.Logging.LogFile == "" {
		logfile = config.Logging.LogDirectory + "Illumio_Parser_Output_" + time.Now().Format("20060102_150405") + ".log"
	} else {
		logfile = config.Logging.LogFile
	}
	f, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	log.SetOutput(f)

	// LOG THE MODE
	log.Printf("INFO - Log only mode set to %t \r\n", config.Logging.LogOnly)

	labelsAPI, apiResp, err := illumioapi.GetAllLabels(pce)
	fmt.Println(labelsAPI, apiResp, err)
	if config.Logging.verbose == true {
		log.Printf("DEBUG - Get All Labels API HTTP Request: %s %v \r\n", apiResp.Request.Method, apiResp.Request.URL)
		log.Printf("DEBUG - Get All Labels API HTTP Reqest Header: %v \r\n", apiResp.Request.Header)
		log.Printf("DEBUG - Get All Labels API Response Status Code: %d \r\n", apiResp.StatusCode)
		log.Printf("DEBUG - Get All Labels API Response Body: \r\n %s \r\n", apiResp.RespBody)
	}
	if err != nil {
		log.Fatal(err)
	}
	parserec := ReadCSV(config.Parser.Parserfile)

	hostrec := ReadCSV(config.Parser.HostnameFile)

	var data regex
	data.load(parserec)

	for _, x := range hostrec {
		fmt.Println(data.LabelFromHostname(x[0]))
	}

}
