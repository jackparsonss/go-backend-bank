package util

import (
	"os"

	"github.com/spf13/viper"
)

// The Config type contains fields for database driver, database source, and server address.
// @property {string} DBDriver - DBDriver is a string property that represents the database driver to
// be used. It is likely used in a configuration file for a Go application. The `mapstructure` tag
// indicates that this property can be mapped to an environment variable or a configuration file key.
// @property {string} DBSource - DBSource is a property that specifies the connection string or data
// source name for the database. It is used to connect to the database using the specified driver.
// @property {string} ServerAddress - The `ServerAddress` property is a string that represents the
// address of the server. It is used to specify the network address on which the server should listen
// for incoming requests. This property is typically used in web applications to specify the IP address
// and port number on which the server should listen for incoming HTTP
type Config struct {
	DBDriver      string `mapstructure:"DB_DRIVER"`
	DBSource      string `mapstructure:"DB_SOURCE"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
}

func LoadConfig(path string) (config Config, err error) {
	if os.Getenv("G_ACTIONS") == "true" {
		viper.SetConfigFile("")
	} else {
		viper.SetConfigFile(path)
	}

	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)

	return
}
