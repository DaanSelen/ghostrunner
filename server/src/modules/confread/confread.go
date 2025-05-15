package confread

import (
	"ghostrunner-server/modules/utilities"

	"gopkg.in/ini.v1"
)

const (
	configSection = "ghostserver"
)

func ReadConf(configPath string) ConfigStruct {
	inidata, err := ini.Load(configPath)
	utilities.HandleError(err, "Trying to load the ini config file!")

	section := inidata.Section(configSection)

	var config ConfigStruct

	config.Address = section.Key("address").String()
	utilities.HandleError(err, "Trying to parse apiport field into the struct!")

	config.Secure, err = section.Key("secure").Bool()
	utilities.HandleError(err, "Trying to parse https field into the struct!")

	config.CertFile = section.Key("certfile").String()
	config.KeyFile = section.Key("keyfile").String()

	config.Interval, err = section.Key("interval").Int()
	utilities.HandleError(err, "Trying to parse interval field into the struct!")

	return config
}
