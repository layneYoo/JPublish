package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

type ImageBuild struct {
	//gitDstPath string
}

func (m ImageBuild) Apply(args []string) {
	Check(len(args) == 5, "four arguments needed")
	if args[0] == "" || args[1] == "" || args[2] == "" {
		log.Fatal("argument null")
		return
	}
	buildPath := args[0]
	registryPath := args[1]
	gitUrl := args[2]
	deployTestTpl := args[3]
	deployProTpl := args[4]
	// testing the path
	_, err := os.Stat(buildPath)
	if err != nil {
		existOr := os.IsExist(err)
		//Check(existOr, "error : ["+buildPath+"] No such directory")
		// not exist, git clone
		if existOr == false {
			out, err := exec.Command("bash", "-c", "/usr/local/bin/git clone "+gitUrl+" "+buildPath).Output()
			Check(err == nil, "git clone error")
			fmt.Println("\n git clone " + string(out))
		}
	} else {
		// exist, git pull
		out, err := exec.Command("bash", "-c", "cd "+buildPath+" && /usr/local/bin/git pull origin master").Output()
		Check(err == nil, "git pull error")
		fmt.Println("\n " + string(out))
	}
	// testing registryPath
	// testing the json
	jsonTestItem := strings.Split(deployTestTpl, ".")
	jsonProItem := strings.Split(deployProTpl, ".")
	Check(len(jsonTestItem) == 2, "the test json syntax error")
	Check(len(jsonProItem) == 2, "the product json syntax error")
	buildCmdHead := "cd "
	buildCmdGitV := ` && /usr/local/bin/git log -1 | head -1 | awk -F" " '{print $2}'`

	// get the commit version
	out, err := exec.Command("bash", "-c", buildCmdHead+buildPath+buildCmdGitV).Output()
	Check(err == nil, "get git version error")
	//gitVersion := string(out[0 : len(out)-2])
	gitVersion := string(out[0:9])

	// build docker image
	buildCmdBuild := ` && docker build -t `
	out, err = exec.Command("bash", "-c", buildCmdHead+buildPath+buildCmdBuild+registryPath+":"+gitVersion+" .").Output()
	Check(err == nil, "build command error:")
	fmt.Println("\n" + string(out))

	// create the marathon's json for deploying
	buildCmdCreateTestjson := " && sed s/VERSION/" + gitVersion + "/g " + deployTestTpl + " > " + jsonTestItem[0] + "_" + gitVersion + ".json"
	buildCmdCreateProductjson := " && sed s/VERSION/" + gitVersion + "/g " + deployProTpl + " > " + jsonProItem[0] + "_" + gitVersion + ".json"
	out, err = exec.Command("bash", "-c", buildCmdHead+buildPath+buildCmdCreateTestjson).Output()
	Check(err == nil, "build test json error:")
	fmt.Println("\n" + string(out))
	out, err = exec.Command("bash", "-c", buildCmdHead+buildPath+buildCmdCreateProductjson).Output()
	Check(err == nil, "build product json error:")
	fmt.Println("\n" + string(out))
}

type ImageUpload struct {
}

func (m ImageUpload) Apply(args []string) {
	Check(len(args) == 2, "two arguments needed")
	if args[0] == "" || args[1] == "" {
		log.Fatal("argument null")
		return
	}
	buildPath := args[0]
	registryPath := args[1]
	buildCmdHead := "cd "
	buildCmdGitV := ` && /usr/local/bin/git log -1 | head -1 | awk -F" " '{print $2}'`

	// get the commit version
	out, err := exec.Command("bash", "-c", buildCmdHead+buildPath+buildCmdGitV).Output()
	Check(err == nil, "get git version error")
	gitVersion := string(out[0:9])

	// push the image[ registryPath:gitVersion ]
	out, err = exec.Command("bash", "-c", "docker push "+registryPath+":"+gitVersion).Output()
	Check(err == nil, "get git version error")
	fmt.Println(string(out))

}
