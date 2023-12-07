package configs

import (
	"flag"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

type EnvConfig struct {
	serviceHost       string
	dbHost            string
	accrualHost       string
	appConfigFilePath string
}

// Package constants
const (
	defaultServiceHost       = "localhost:8080"
	defaultDBHost            = "localhost:8080"
	defaultAccrualHost       = "localhost:8080"
	defaultAppConfigFilePath = "AppConfig.yaml"

	envNameServiceHost   = "RUN_ADDRESS"
	envNameDBHost        = "DATABASE_URI"
	envNameAccrualHost   = "ACCRUAL_SYSTEM_ADDRESS"
	envAppConfigFilePath = "GOFERMART_CONFIG_PATH"
)

func newEnvConfig() *EnvConfig {

	envConfig := EnvConfig{
		serviceHost: defaultServiceHost,
		dbHost: defaultDBHost,
		accrualHost: defaultAccrualHost,
		appConfigFilePath: defaultAppConfigFilePath,
	}

	envConfig.parseFlags()
	envConfig.parseEnvironmental()

	return &envConfig
}

// Parse application flags
func (c *EnvConfig) parseFlags() {

	pc := "func (c *Config) parseFlags()"

	serviceHost := *flag.String("a", "", "Service host addres. Example: 127.0.0.1:8080")
	dbHost := *flag.String("d", "", "Database host addres. Example: 127.0.0.1:8080")
	accrualHost := *flag.String("r", "", "Accrual system address. Example: 127.0.0.1:8080")
	appConfigFilePath := *flag.String("c", "", "Application configuration file path")

	flag.Parse()

	if isValidHostAddres(serviceHost) {
		c.serviceHost = serviceHost
	} else {
		log.Printf("Using value by default. %s: serviceHost value is not valid: %s", pc, serviceHost)
	}

	if isValidHostAddres(dbHost) {
		c.dbHost = dbHost
	} else {
		log.Printf("Using value by default. %s: dbHost value is not valid: %s", pc, dbHost)
	}

	if isValidHostAddres(accrualHost) {
		c.accrualHost = accrualHost
	} else {
		log.Printf("Using value by default. %s: accrualHost value is not valid: %s", pc, accrualHost)
	}

	if appConfigFilePath != "" {
		c.appConfigFilePath = appConfigFilePath
	} else {
		log.Printf("Using value by default. %s: appConfigFilePath value is not valid: %s", pc, appConfigFilePath)
	}
}

// Parse environmental variables
func (c *EnvConfig) parseEnvironmental() {

	pc := "func (c *Config)parseEnvironmental()"

	if serviceHost, ok := os.LookupEnv(envNameServiceHost); ok {
		if isValidHostAddres(serviceHost) {
			c.serviceHost = serviceHost
		} else {
			log.Printf("Using value by default. %s: serviceHost value is not valid: %s", pc, serviceHost)
		}
	}

	if dbHost, ok := os.LookupEnv(envNameDBHost); ok {
		if isValidHostAddres(dbHost) {
			c.dbHost = dbHost
		} else {
			log.Printf("Using value by default. %s: dbHost value is not valid: %s", pc, dbHost)
		}
	}

	if accrualHost, ok := os.LookupEnv(envNameAccrualHost); ok {
		if isValidHostAddres(accrualHost) {
			c.serviceHost = accrualHost
		} else {
			log.Printf("Using value by default. %s: accrualHost value is not valid: %s", pc, accrualHost)
		}
	}

	if appConfigFilePath, ok := os.LookupEnv(envAppConfigFilePath); ok {
		if appConfigFilePath != "" {
			c.appConfigFilePath = appConfigFilePath
		} else {
			log.Printf("Using value by default. %s: appConfigFilePath value is not valid: %s", pc, appConfigFilePath)
		}
	}
}

// Checks the address for validity
func isValidHostAddres(addres string) bool {

	args := strings.Split(addres, ":")
	if len(args) != 2 {
		return false
	}

	if args[0] != "localhost" && net.ParseIP(args[0]) == nil {
		return false
	}

	if _, err := strconv.Atoi(args[1]); err != nil {
		return false
	}
	return true
}
