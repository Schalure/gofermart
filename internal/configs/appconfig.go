package configs

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

//	Application configugation struct
type AppConfig struct {
	Env Environment `yaml:"Environment"`
}

//	Values for `yaml:"Environment"`
type Environment string
const (
	Debug Environment = "debug"
	Local Environment = "local"
	Prod Environment = "prod"
)

func newAppConfig(fileName string) (*AppConfig, error) {

	appConfig := AppConfig {
		Env: Debug,
	}

	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		err := appConfig.createAppConfigFile(fileName)
		return &appConfig, errors.Join(fmt.Errorf("using application config with defaul values"), err)
	}

	if err := appConfig.getAppConfigFromFile(fileName); err != nil {
		return &appConfig, errors.Join(fmt.Errorf("using application config with defaul values"), err)
	}

	return &appConfig, nil
}

//	Create application config file with default values
func (c *AppConfig) createAppConfigFile(fileName string) error {

	file, err :=  os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("can't create application config file: \"%s\"", fileName)
	}
	defer file.Close()

	out, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("can't marshal application config")
	}

	if _, err := file.Write(out); err != nil {
		return fmt.Errorf("can't write application config to file")
	}
	return nil
}

//	Get application config from file
func (c *AppConfig) getAppConfigFromFile(fileName string) error {

	data, err := os.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("can't read application config file: \"%s\"", fileName)
	}

	if err := yaml.Unmarshal(data, c); err != nil {
		return fmt.Errorf("can't unmarshal application config file")
	}
	return nil
}