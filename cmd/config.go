package cmd

import (
	"os"
	"path/filepath"
    "log"

	"gopkg.in/yaml.v3"
)

func readConfig() error {

	pwd, err := os.Getwd()
    check(err)

	dir, err := os.UserConfigDir()
    check(err)
	
	configPath := filepath.Join(dir, "money-tree")

    // Create money-tree directory in UserConfigDir
	if err = os.Mkdir(configPath, 0750); err != nil && !os.IsExist(err) {
		log.Fatalln(err)
	}
    
	source.Config.Database = filepath.Join(pwd, "database.db")

    // Create config.yml in config directory
	configPath = filepath.Join(configPath, "config.yml")
	file, err := os.Open(configPath)
	if err != nil && !os.IsExist(err) {
		initConfig(pwd, configPath)
		return nil
	}
	defer file.Close()

    // Read yaml from file
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&source.Config)
	check(err)

	if err = os.Mkdir(source.Config.Attachments, 0750); err != nil && !os.IsExist(err) {
		check(err)
	}
	return nil
}

func initConfig(pwd, configPath string) {

	file, err := os.Create(configPath)
	check(err)

	source.Config.Path = pwd

	defer file.Close()

	attachmentsPath := filepath.Join(pwd, "attachments")
	if err = os.Mkdir(attachmentsPath, 0750); err != nil && !os.IsExist(err) {
		check(err)
	}

	source.Config.Attachments = attachmentsPath
    
    // Write yaml to file
	encoder := yaml.NewEncoder(file)
	err = encoder.Encode(&source.Config)
	check(err)
}
