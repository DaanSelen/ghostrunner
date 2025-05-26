package utilities

import (
	"log"

	"gopkg.in/ini.v1"
)

const (
	configSection = "ghostserver"
)

func ReadConf(configPath string) ConfigStruct {
	inidata, err := ini.Load(configPath)
	if err != nil {
		log.Println(ErrTag, err)
	}

	section := inidata.Section(configSection)

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

	return config
}
