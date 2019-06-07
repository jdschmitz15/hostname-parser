package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/brian1917/illumioapi"
)

type config struct {
	Illumio    illumio    `toml:"illumio"`
	Parser     parser     `toml:"parser"`
	LabelMatch labelMatch `toml:"labelMatch"`
	Logging    logging    `toml:"logging"`
}

type illumio struct {
	FQDN       string `toml:"fqdn"`
	Port       int    `toml:"port"`
	Org        int    `toml:"org"`
	User       string `toml:"user"`
	Key        string `toml:"key"`
	MatchField string `toml:"match_field"`
}

type parser struct {
	Parserfile   string `toml:"parserfile"`
	HostnameFile string `toml:"hostnamefile"`
	OutputFile   string `toml:"outputfile"`
	UseApi       bool   `toml:"useapi"`
	AutoUpdate   bool   `toml:"autoupdate"`
}
type labelMatch struct {
	App         string `toml:"app"`
	Enviornment string `toml:"env"`
	Location    string `toml:"loc"`
	Role        string `toml:"role"`
}
type logging struct {
	LogOnly      bool   `toml:"log_only"`
	LogDirectory string `toml:"log_directory"`
	LogFile      string `toml:"log_file"`
	verbose      bool
}

var configFile, hostFile, outputFile string
var debugLogging, api, auto bool

func init() {
	flag.StringVar(&configFile, "config", "config.toml", "Location of TOML configuration file")
	flag.BoolVar(&debugLogging, "v", false, "Set for verbose logging.")
	flag.BoolVar(&api, "api", false, "Set to pull hostnames from PCE.")
	flag.BoolVar(&auto, "auto", false, "Set to automate label updates on PCE.")
	flag.StringVar(&hostFile, "hostfile", "hostname.csv", "Location of hostnames CSV to parse")
	flag.StringVar(&outputFile, "outputfile", "", "Location of hostnames CSV to parse")

}

func parseConfig() (config, illumioapi.PCE) {
	var config config

	flag.Parse()

	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatal(err)
	}

	_, err = toml.Decode(string(data), &config)
	if err != nil {
		log.Fatal(err)
	}

	// IF A FIELD IS LEFT BLANK WE WANT TO PUT A PLACEHOLDER
	fields := []*string{&config.LabelMatch.App, &config.LabelMatch.Enviornment, &config.LabelMatch.Location, &config.LabelMatch.Role}
	for _, field := range fields {
		if *field == "" {
			*field = "csvPlaceHolderIllumio"
		}

	}

	pce := illumioapi.PCE{
		FQDN: config.Illumio.FQDN,
		Port: config.Illumio.Port,
		Org:  config.Illumio.Org,
		User: config.Illumio.User,
		Key:  config.Illumio.Key}

	// SET THE LOGGING IN THE CONFIG STRUCT
	config.Logging.verbose = debugLogging

	// SET THE Hostname File in Config
	if config.Parser.HostnameFile == "" {
		config.Parser.HostnameFile = hostFile
	}

	// Override toml weather to Pull Hostnames from PCE or not
	if api && !config.Parser.UseApi {
		config.Parser.UseApi = api
	}

	// Pull Hostnames from PCE
	if auto && !config.Parser.AutoUpdate {
		config.Parser.AutoUpdate = auto
	}

	fmt.Printf("%+v", config)
	return config, pce
}
