package main

import (
	"flag"
	"io/ioutil"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/brian1917/illumioapi"
)

type config struct {
	Illumio      illumio      `toml:"illumio"`
	Parser       parser       `toml:"parser"`
	LabelMapping labelMapping `toml:"labelMapping"`
	Logging      logging      `toml:"logging"`
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
	parserfile  string `toml:"file"`
	samplehosts string `toml:"samplehosts"`
}
type labelMapping struct {
	App         string `toml:"app"`
	Enviornment string `toml:"enviornment"`
	Location    string `toml:"location"`
	Role        string `toml:"role"`
}
type logging struct {
	LogOnly      bool   `toml:"log_only"`
	LogDirectory string `toml:"log_directory"`
	verbose      bool
}

var configFile string
var debugLogging bool

func init() {
	flag.StringVar(&configFile, "config", "config.toml", "Location of TOML configuration file")
	flag.BoolVar(&debugLogging, "v", false, "Set for verbose logging.")
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
	fields := []*string{&config.LabelMapping.App, &config.LabelMapping.Enviornment, &config.LabelMapping.Location, &config.LabelMapping.Role}
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

	return config, pce
}
