package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/brian1917/illumioapi"
)

type config struct {
	Illumio illumio `toml:"illumio"`
	Parser  parser  `toml:"parser"`
	Match   match   `toml:"match"`
	Logging logging `toml:"logging"`
}

type illumio struct {
	FQDN string `toml:"fqdn"`
	Port int    `toml:"port"`
	Org  int    `toml:"org"`
	User string `toml:"user"`
	Key  string `toml:"key"`
	//	MatchField string `toml:"match_field"`
	NoPCE bool `toml:"no_pce"`
}

type parser struct {
	Parserfile   string `toml:"parserfile"`
	HostnameFile string `toml:"hostnamefile"`
	OutputFile   string `toml:"outputfile"`
	NoPrompt     bool   `toml:"noprompt"`
}
type match struct {
	AllEmpty    bool   `toml:"allempty"`
	IgnoreMatch bool   `toml:"ignorematch"`
	App         string `toml:"app"`
	Env         string `toml:"env"`
	Loc         string `toml:"loc"`
	Role        string `toml:"role"`
}
type logging struct {
	LogOnly      bool   `toml:"log_only"`
	LogDirectory string `toml:"log_directory"`
	LogFile      string `toml:"log_file"`
	verbose      bool
}

var configFile, hostFile, outputFile string
var debugLogging, api, noprompt bool

func init() {
	flag.StringVar(&configFile, "config", "config.toml", "Location of TOML configuration file")
	flag.BoolVar(&debugLogging, "v", false, "Set for verbose logging.")
	flag.BoolVar(&noprompt, "noprompt", false, "Set to automate label updates on PCE.")
	flag.StringVar(&hostFile, "hostfile", "", "Location of hostnames CSV to parse")
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
	//fields := []*string{&config.LabelMatch.App, &config.LabelMatch.Env, &config.LabelMatch.Loc, &config.LabelMatch.Role}
	// for _, field := range fields {
	// 	if *field == "" {
	// 		*field = "csvPlaceHolderIllumio"
	// 	}

	// }

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

	// Pull Hostnames from PCE
	if noprompt {
		config.Parser.NoPrompt = true

	}

	// if config.Parser.OutputFile != "" {
	// 	config.Parser.NoPrompt = false
	// }

	if config.Illumio.NoPCE && config.Parser.HostnameFile == "" {
		fmt.Printf("\r\nYou must use the CLI -hostfile option or configure HostfileName in the config file when not using PCE Data(no_pce=true)....r\n")
		os.Exit(1)
	}
	//fmt.Printf("%+v", config)
	return config, pce

}
