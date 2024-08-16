package cmd

import (
	"os"
	"path/filepath"
    
    "gopkg.in/yaml.v3"
)

func readConfig() error {

    dir, err := os.UserConfigDir()
    check(err)

    pwd, err := os.Getwd()
    check(err)
    
    source.Config.Database = filepath.Join(pwd, "database.db")

    configPath := filepath.Join(dir, "money-tree") 
    if err = os.Mkdir(configPath, 0750); err != nil && !os.IsExist(err) {
        check(err)
    }

    configPath = filepath.Join(configPath, "config.yml")
    
    file, err := os.Open(configPath)
    if err != nil && !os.IsExist(err) {
        initConfig(pwd, configPath)
        return nil
    }

    defer file.Close()

    decoder := yaml.NewDecoder(file)    
    err = decoder.Decode(&source.Config)
    check(err) 
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
    encoder := yaml.NewEncoder(file)
    err = encoder.Encode(&source.Config)
    check(err)
}
