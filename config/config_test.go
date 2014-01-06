package config

import (
	"github.com/ghthor/gospec"
	. "github.com/ghthor/gospec"
)

func DescribeConfigLoading(c gospec.Context) {
	c.Specify("Config can be parsed from json file", func() {
		expectedConfig := ServerConfig{
			"127.0.0.1",
			80,
			443,
			"tls-gen/cert.pem",
			"tls-gen/key.pem",

			DatabaseConfig{
				"dbuser",
				"dbpassword",
				"dbname",
				"filepath/to/filedb",
			},
		}

		config, err := ReadFromFile("config.default.json")
		c.Assume(err, IsNil)
		c.Expect(config, Equals, expectedConfig)
	})
}
