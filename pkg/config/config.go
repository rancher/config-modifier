package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"gopkg.in/yaml.v2"
)

func fileExists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func PlaceConfigFile(hostPath, configPath string, nodeLabels []string) error {
	var (
		token     string
		datastore string
	)
	if fileExists(hostPath) {
		fmt.Println("config file exists on the host")
		filetxt, err := ioutil.ReadFile(hostPath)
		if err != nil {
			return err
		}
		config := make(map[string]interface{})
		err = yaml.Unmarshal(filetxt, config)
		if err != nil {
			return err
		}
		if config["token"] != "" {
			token = fmt.Sprintf("%v", config["token"])
		}
		if config["datastore-endpoint"] != "" {
			datastore = fmt.Sprintf("%v", config["datastore-endpoint"])
		}
	}
	configs, err := ioutil.ReadDir(configPath)
	if err != nil {
		return err
	}
	for _, label := range nodeLabels {
		for _, confFile := range configs {
			fmt.Println("label ,config = ", label, confFile.Name())
			if lbl := strings.Replace(label, "=", "-", -1); confFile.Name() == lbl {
				content, err := ioutil.ReadFile(path.Join(configPath, confFile.Name()))
				if err != nil {
					return err
				}
				contentMap := make(map[string]interface{})
				err = yaml.Unmarshal(content, contentMap)
				if token != "" {
					contentMap["token"] = token
				}
				if datastore != "" {
					contentMap["datastore-endpoint"] = datastore
				}
				if err != nil {
					return err
				}
				content, err = yaml.Marshal(contentMap)
				if err != nil {
					return err
				}
				return ioutil.WriteFile(hostPath, content, 0755)
			}
		}
	}
	return nil
}
