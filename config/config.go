package config

var Conf = &YamlConfigs{}

func init() {
	Conf = Load("./config.yml")
}
