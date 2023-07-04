package main

import (
	"fmt"
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

	// Get a greeting message and print it.
	message := fmt.Sprintf("Battery %s loaded %2f", bat.GetName(), bat.GetSoc())
	fmt.Printf(message)
}
