package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	config, err := os.Open("/home/xlei/code/golang/src/github.com/JPublish/conf/cfg.json")
	defer config.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParse := json.NewDecoder(config)
	if jsonParse == nil {
		return
	}
	var marathon MarathonObj
	if err = jsonParse.Decode(&marathon); err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(marathon.Marathoninfo)
}
