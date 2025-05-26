package utilities

import (
	"log"

	"gopkg.in/ini.v1"
)

const (
	serverSection = "ghostserver"
	runnerSection = "ghostrunner"
)

func ReadConf(configPath string) ConfigStruct {
	inidata, err := ini.Load(configPath)
	if err != nil {
		log.Println(ErrTag, err)
	}

	section := inidata.Section(serverSection)

	var config ConfigStruct

	// System
	config.Address = section.Key("address").String()

	// Authentication
	config.AdminToken = section.Key("admin_token").String()
	config.TokenKeyFile = section.Key("token_key_file").String()

	// Protocol
	config.Secure, err = section.Key("secure").Bool()
	if err != nil {
		log.Println(ErrTag, err)
	}

	// API Protocol Certificate
	config.ApiCertFile = section.Key("api_cert_file").String()
	config.ApiKeyFile = section.Key("api_key_file").String()

	// Service Configuration
	config.Interval, err = section.Key("interval").Int()
	if err != nil {
		log.Println(ErrTag, err)
	}

	section = inidata.Section(runnerSection)

	config.MeshHostname = section.Key("hostname").String()
	config.MeshUsername = section.Key("username").String()
	config.MeshPassword = section.Key("password").String()

	config.PyVenvName = section.Key("python_venv_name").String()

	return config
}
