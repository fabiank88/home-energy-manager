package main

import (
	"flag"
	"hem/battery"
	"hem/evcc"
	"log"
)

func main() {

	configFile := flag.String(`config`, `/etc/hem.yaml`, `reference to the config file yaml`)
	flag.Parse()

	//read yaml config file
	config := GetConf(*configFile)

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
