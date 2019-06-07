package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/brian1917/illumioapi"
)

type regex struct {
	regexdata []regexstruct
}
type regexstruct struct {
	regex     string
	cglist    [][]string
	precgtxt  []string
	postcgtxt []string
}

type wkld struct {
	hostname string
	labels   []string
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
// func ParseMatch(regex []regexdata,hostname string) int {
// 	return int
// }

func (r *regex) ParseHostname(hostname string) wkld {

	var tempwkld wkld
	var match []string

	found := false
	for _, tmp := range r.regexdata {
		if !found {
			tmpre := regexp.MustCompile(tmp.regex)
			match = tmpre.FindStringSubmatch(hostname)
			//fmt.Println(match)
			if len(match) > 0 {
				found = true
				for label := 0; label < 4; label++ {
					var tmpstring string

					tmpstring = tmp.precgtxt[label]
					for _, x := range tmp.cglist[label] {
						if x != "" {
							idx, _ := strconv.Atoi(x)
							tmpstring = tmpstring + match[idx]
						}

					}
					tmpstring = tmpstring + tmp.postcgtxt[label]
					tempwkld.labels = append(tempwkld.labels, tmpstring)
				}
			}
		}
	}
	tempwkld.hostname = hostname
	return tempwkld
}

//Load the Regex CSV Into the parser struct -
func (reg *regex) load(data [][]string) {

	var regex []regexstruct
	//regex that parses each labels capture groups and any static text you want.
	// eg. Pod[1], Pod[1,2], [1,3,4]backend
	re := regexp.MustCompile(`([A-Za-z0-9\-\_]*)(\[?)([0-9\,]*)(\]?)([A-Za-z0-9\-\_]*)`)

	//Cycle through all the parse data rows in the parse data xls
	for c, row := range data {
		//ignore header
		if c != 0 {
			var r regexstruct
			// var tempcglabels [][]string
			// var tempprecglabels, temppostcglabels []string

			//Array order 0-LOC,1-ENV,2-APP,3-APP
			for x := 0; x < 4; x++ {

				match := re.FindStringSubmatch(row[x+1])
				//fmt.Println(match)

				//Check if the match criteria
				if len(match) > 0 {
					r.cglist = append(r.cglist, strings.Split(match[3], ","))
					r.precgtxt = append(r.precgtxt, match[1])
					r.postcgtxt = append(r.postcgtxt, match[5])
				} else {
					fmt.Printf("parser CSV has incorrect format in row %d and column %d", c, x)
					os.Exit(1)
				}
			}
			r.regex = row[0]

			// r.precgtxt = tempprecglabels
			// r.postcgtxt = temppostcglabels
			// r.cglist = tempcglabels

			regex = append(regex, r)
		}
		reg.regexdata = regex
	}
	//fmt.Printf("%+v", reg)
}

func main() {

	config, pce := parseConfig()

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

	// labelsAPI, apiResp, err := illumioapi.GetAllLabels(pce)
	// fmt.Println(labelsAPI, apiResp, err)
	// if config.Logging.verbose == true {
	// 	log.Printf("DEBUG - Get All Labels API HTTP Request: %s %v \r\n", apiResp.Request.Method, apiResp.Request.URL)
	// 	log.Printf("DEBUG - Get All Labels API HTTP Reqest Header: %v \r\n", apiResp.Request.Header)
	// 	log.Printf("DEBUG - Get All Labels API Response Status Code: %d \r\n", apiResp.StatusCode)
	// 	log.Printf("DEBUG - Get All Labels API Response Body: \r\n %s \r\n", apiResp.RespBody)
	// }
	// if err != nil {
	// 	log.Fatal(err)
	// }

	workloads, apiResp, err := illumioapi.GetAllWorkloads(pce)
	fmt.Println(workloads, apiResp, err)
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
		fmt.Println(data.ParseHostname(x[0]))
	}
	re := regexp.MustCompile(`E1(D)(1)(2)(3)(S)(S)(1)(7)(0)`)
	fmt.Println(re.FindStringSubmatch("E1D123SS170"))
	fmt.Println(re.ReplaceAllString("E1D123SS170", "$"))
}
