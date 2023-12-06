// The package describes the structures and configuration methods for the service
package configs

// Main configuration struct
type Config struct {
	EnvConfig EnvConfig
	AppConfig AppConfig
}


// Consructor of Config object
func NewConfig() (*Config, error) {

	envConfig := newEnvConfig()
	appConfig, err := newAppConfig(envConfig.appConfigFilePath)


	return &Config{
		EnvConfig: *envConfig,
		AppConfig: *appConfig,
	}, err
}

