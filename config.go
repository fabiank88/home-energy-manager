package main

import (
	"gopkg.in/yaml.v3"
	"hem/battery"
	"hem/evcc"
	"log"
	"os"
)

type config struct {
	Battery battery.Config
	Evcc    evcc.Config
}

func GetConf(filepath string) config {

	file, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatal("could not read config")
	}
	c := config{}
	err = yaml.Unmarshal(file, &c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}
