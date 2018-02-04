package configs

import (
	"bytes"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/spf13/viper"
)

var Settings *viper.Viper

func init() {
	Settings = LoadViperFromToml("", "go-codex.toml", path.Join("configs", "go-codex.toml"))
}

func existsFile(filepath string) bool {
	if _, err := os.Stat(filepath); err == nil {
		return true
	}
	return false
}

func LoadViperFromToml(encrptyKey string, configFilePath ...string) *viper.Viper {

	settings := viper.New()
	settings.SetConfigType("toml")
	var buff []byte
	var err error
	for _, configFile := range configFilePath {
		buff = nil
		err = nil
		if strings.HasPrefix(strings.ToLower(configFile), "http://") || strings.HasPrefix(strings.ToLower(configFile), "https://") {
			url := configFile
			buff, err = requestEncryptTOML(url, nil, nil)
		} else {
			if existsFile(configFile) {
				buff, err = ioutil.ReadFile(configFile)
			}
		}

		if err == nil {
			settings.ReadConfig(bytes.NewBuffer(buff))
			break
		}
	}
	return settings
}

func LoadToml(configFile string, v interface{}) error {
	if _, err := toml.DecodeFile(configFile, v); err != nil {
		return err
	}
	return nil
}
