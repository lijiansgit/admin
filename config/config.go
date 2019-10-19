package config

import (
	"flag"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

var (
	Conf     = &Config{}
	ConfFile = "config.yaml"
)

func init() {
	flag.StringVar(&ConfFile, "c", "config/config.yaml", " set config file path")
}

type Config struct {
	WEB  `yaml:"web"`
	DB   `yaml:"db"`
	LDAP `yaml:"ldap"`
	Log  `yaml:"log"`
}

type WEB struct {
	Addr string `yaml:"addr"`
}

type DB struct {
	File string `yaml:"file"`
}

type LDAP struct {
	Addr     string `yaml:"addr"`
	RootDN   string `yaml:"rootDN"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Log struct {
	Level int    `yaml:"level"`
	Conf  string `yaml:"conf"`
}

func Init() (err error) {
	bf, err := ioutil.ReadFile(ConfFile)
	if err != nil {
		return err
	}

	if err = yaml.Unmarshal(bf, Conf); err != nil {
		return err
	}

	return nil
}
