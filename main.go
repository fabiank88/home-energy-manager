package main

import (
	"hem/battery"
	"hem/evcc"
	"log"
)

func main() {
	//read yaml config file
	config := GetConf(`config.yaml`)

	bat, err := battery.New(config.Battery)
	if err != nil {
		log.Fatal(err)
	}

	evcc, err := evcc.New(config.Evcc)
	if err != nil {
		log.Fatal(err)
	}

	engine := Engine{battery: bat, evcc: evcc}
	engine.run()
}
