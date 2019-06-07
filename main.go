package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
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
			var tempcglabels [][]string
			var tempprecglabels, temppostcglabels []string
			for x := 0; x < 4; x++ {

				match := re.FindStringSubmatch(row[x+1])
				//fmt.Println(match)
				if len(match) > 0 {
					tempcglabels = append(tempcglabels, strings.Split(match[3], ","))
					tempprecglabels = append(tempprecglabels, match[1])
					temppostcglabels = append(temppostcglabels, match[5])
				}
			}
			r.precgtxt = tempprecglabels
			r.postcgtxt = temppostcglabels
			r.regex = row[0]
			r.cglist = tempcglabels

			regex = append(regex, r)
		}
		reg.regexdata = regex
	}
	//fmt.Println(r)
}

func main() {
	var parsefile, hostfile string

	parse := flag.String("parse", "parse-table.csv", "Enter the parse csv")
	names := flag.String("names", "hostname.csv", "Enter the hostname csv")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage:  --parse <parser file> --names <hostname file>\r\n")
		fmt.Fprintf(os.Stderr, "default:  parse: 'parse-table.csv' \r\n")
		fmt.Fprintf(os.Stderr, "default:  names: 'hostname.csv' \r\n")
	}
	flag.Parse()

	//fmt.Println(*parse)
	if *parse != "" {
		parsefile = *parse
	}
	if *names != "" {
		hostfile = *names
	}
	parserec := ReadCSV(parsefile)
	hostrec := ReadCSV(hostfile)

	var data regex
	data.load(parserec)
	for _, x := range hostrec {
		fmt.Println(data.ParseHostname(x[0]))
	}
}
