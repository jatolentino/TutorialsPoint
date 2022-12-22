package config

import (
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// LoadConfig loads configurations from YAMLs.
// Supposed to be invoked in main.go.
// Takes file name as the argument.
func LoadConfig(v1 *viper.Viper) {
	v1.SetConfigName("application")
	v1.SetConfigType("yaml")
	v1.AddConfigPath(".")
	err := v1.ReadInConfig()
	if err != nil {
		panic("Couldn't load configuration, cannot start. Terminating. Error: " + err.Error())
	}
	log.Println("Config loaded successfully...")
	log.Println("Getting environment variables...")
	for _, k := range v1.AllKeys() {
		value := v1.GetString(k)
		if strings.HasPrefix(value, "${") && strings.HasSuffix(value, "}") {
			values := strings.SplitAfterN(strings.TrimSuffix(strings.TrimPrefix(value, "${"), "}"), ":", 2)
			val, isPresent := os.LookupEnv(values[0])
			if isPresent == false {
				v1.Set(k, values[1])
			} else {
				v1.Set(k, val)
			}
		}
	}
}