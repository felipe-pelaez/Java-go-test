package config

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"io"
	"log"
	"net/http"
)

// ? =========================== Structs =========================== ?

// springCloudConfig is a struct for parsing spring cloud config
type springCloudConfig struct {
	Name            string           `json:"name"`
	Profiles        []string         `json:"profiles"`
	Label           string           `json:"label"`
	Version         string           `json:"version"`
	PropertySources []propertySource `json:"propertySources"`
}

// * ============

// propertySource is a struct for parsing spring cloud config
type propertySource struct {
	Name   string                 `json:"name"`
	Source map[string]interface{} `json:"source"`
}

// ? =========================== Functions =========================== ?

// LoadConfigurationFromBranch loads configuration from config server
func LoadConfigurationFromBranch(configServerUrl string, appName string, profile string, branch string) {

	url := fmt.Sprintf("%s/%s/%s/%s", configServerUrl, appName, profile, branch)
	log.Printf("loading config from %s\n", url)

	body, err := fetchConfiguration(url)
	if err != nil {
		log.Fatalln("couldn't load configuration, cannot start. terminating. error: " + err.Error())
	}

	parseConfiguration(body)

}

// * ============

// FetchConfiguration fetches configuration from config server
func fetchConfiguration(url string) ([]byte, error) {

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln("couldn't load configuration, cannot start. terminating. error: " + err.Error())
	}
	body, err := io.ReadAll(resp.Body)
	return body, err
}

// * ============

// ParseConfiguration parses configuration from config server
func parseConfiguration(body []byte) {

	var cloudConfig springCloudConfig

	err := json.Unmarshal(body, &cloudConfig)
	if err != nil {
		log.Fatalln("cannot parse configuration, message: " + err.Error())
	}

	for key, value := range cloudConfig.PropertySources[0].Source {
		viper.Set(key, value)
	}

	if cloudConfig.Name != "" {
		log.Printf("successfully loaded configuration for service %s\n", cloudConfig.Name)
	}

}
