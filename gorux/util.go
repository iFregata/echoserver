package gorux

import (
	"encoding/json"
	"log"
	"os"
)

func LoadConfigFile(file string, conf interface{}) interface{} {
	f, err := os.Open(file)
	if err != nil {
		log.Fatal("Can't open config file: ", err)
	}
	defer f.Close()
	jsonParser := json.NewDecoder(f)
	err = jsonParser.Decode(conf)
	if err != nil {
		log.Fatal("Parse config failed: ", err)
	}
	return conf
}
