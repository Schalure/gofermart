package configs

import (
	"errors"
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

//	Application configugation struct
type AppConfig struct {
	Env Environment `yaml:"Environment"`
	PassRules string `yaml:"PasswordRules"`
	TokenTTL time.Duration `yaml:"TokenTimeToLife"`
}

//	Values for `yaml:"Environment"`
type Environment string
const (
	Debug Environment = "debug"
	Local Environment = "local"
	Prod Environment = "prod"
)

//	Default application configugation values
const (
	defaultEnv = Debug
	defaultPassRules = `[0-9a-zA-Z]`
	defaultTokenTTL = time.Hour * 1
)

//	Application configuration constructor
func newAppConfig(fileName string) (*AppConfig, error) {

	appConfig := AppConfig {
		Env: defaultEnv,
		PassRules: defaultPassRules,
		TokenTTL: defaultTokenTTL,
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

	file, err :=  os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0644)
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