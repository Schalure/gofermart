// The package describes the structures and configuration methods for the service
package configs

import (
	"flag"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

// Main configuration struct
type Config struct {
	serviceHost string
	dbHost      string
	accrualHost string
}

// Package constants
const (
	defaultServiceHost = "localhost:8080"
	defaultDBHost      = "localhost:8080"
	defaultAccrualHost = "localhost:8080"

	envNameServiceHost = "RUN_ADDRESS"
	envNameDBHost      = "DATABASE_URI"
	envNameAccrualHost = "ACCRUAL_SYSTEM_ADDRESS"
)

// Consructor of Config object
func NewConfig() *Config {

	config := new(Config)

	config.serviceHost = defaultServiceHost
	config.dbHost = defaultDBHost
	config.accrualHost = defaultAccrualHost

	config.parseFlags()
	config.parseEnvironmental()

	return config
}

// Parse application flags
func (c *Config) parseFlags() {

	pc := "func (c *Config) parseFlags()"

	serviceHost := *flag.String("a", "", "Service host addres. Example: 127.0.0.1:8080")
	dbHost := *flag.String("d", "", "Database host addres. Example: 127.0.0.1:8080")
	accrualHost := *flag.String("r", "", "Accrual system address. Example: 127.0.0.1:8080")
	flag.Parse()

	if isValidHostAddres(serviceHost) {
		c.serviceHost = serviceHost
	} else {
		log.Printf("%s: serviceHost value is not valid: %s", pc, serviceHost)
	}

	if isValidHostAddres(dbHost) {
		c.dbHost = dbHost
	} else {
		log.Printf("%s: dbHost value is not valid: %s", pc, dbHost)
	}

	if isValidHostAddres(accrualHost) {
		c.accrualHost = accrualHost
	} else {
		log.Printf("%s: accrualHost value is not valid: %s", pc, accrualHost)
	}
}

// Parse environmental variables
func (c *Config) parseEnvironmental() {

	pc := "func (c *Config)parseEnvironmental()"

	if serviceHost, ok := os.LookupEnv(envNameServiceHost); ok {
		if isValidHostAddres(serviceHost) {
			c.serviceHost = serviceHost
		} else {
			log.Printf("%s: serviceHost value is not valid: %s", pc, serviceHost)
		}
	}

	if dbHost, ok := os.LookupEnv(envNameDBHost); ok {
		if isValidHostAddres(dbHost) {
			c.dbHost = dbHost
		} else {
			log.Printf("%s: dbHost value is not valid: %s", pc, dbHost)
		}
	}

	if accrualHost, ok := os.LookupEnv(envNameAccrualHost); ok {
		if isValidHostAddres(accrualHost) {
			c.serviceHost = accrualHost
		} else {
			log.Printf("%s: accrualHost value is not valid: %s", pc, accrualHost)
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
