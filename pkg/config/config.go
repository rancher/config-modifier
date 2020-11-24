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

func PlaceConfigFile(hostPath, configPath string, nodeLabels, preservedEntries []string) error {
	preserved := make(map[string]interface{})
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
		for _, entry := range preservedEntries {
			if config[entry] != nil {
				preserved[entry] = config[entry]
			}
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
				if err != nil {
					return err
				}
				for entry := range preserved {
					contentMap[entry] = preserved[entry]
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
