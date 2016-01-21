package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

func parseargs() (config, marathonaction string) {
	flag.StringVar(&config, "j", "", "config file")
	flag.StringVar(&marathonaction, "a", "", "marathon action, eg:app list")
	//flag.StringVar(&format, "f", "", "output format")
	flag.Parse()
	return
}

func JsonGet() MarathonObj {
	configFile, atc := parseargs()
	// todo : test the config
	//
	config, err := os.Open(configFile)
	defer config.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParse := json.NewDecoder(config)
	Check(jsonParse != nil, "json config decode error...")
	var marathonObj MarathonObj
	if err = jsonParse.Decode(&marathonObj); err != nil {
		fmt.Println(err.Error())
	}
	if atc != "" {
		marathonObj.Actioninfo.Act = atc
	}
	return marathonObj
}

func JsonConfig() MarathonObj {
	marathonObj := JsonGet()
	// test info

	return marathonObj
}
