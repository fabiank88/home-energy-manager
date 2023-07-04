package main

import (
	"fmt"
	"hem/battery"
	"hem/evcc"
	"log"
	"time"
)

type Engine struct {
	battery battery.Battery
	evcc    evcc.Evcc
}

func (e *Engine) run() {
	for {
		err := e.evcc.Refresh()
		if err != nil {
			log.Fatal(err.Error())
		}

		err = e.battery.Refresh()
		if err != nil {
			log.Fatal(err.Error())
		}

		fmt.Println(`------------ Status --------------`)
		fmt.Println(`Battery Mode: `, battery.OperatingModeToString(e.battery.GetOperatingMode()))
		fmt.Println(`Battery SOC: `, e.battery.GetSoc())
		fmt.Println(`Battery Flow: `, e.battery.GetFlow()*-1)
		fmt.Println(`EVCC Power Mode: `, e.evcc.IsPowerCharging())
		fmt.Println(`----------------------------------`)

		if e.evcc.IsPowerCharging() && e.battery.GetOperatingMode() == battery.OPERATING_MODE_AUTOMATIC {
			fmt.Println(`EVCC is power charging, but battery is enabled - disable battery`)
			err := e.battery.SetOperatingModeManual(0)
			if err != nil {
				log.Fatal(err.Error())
			}
		}

		if !e.evcc.IsPowerCharging() && e.battery.GetOperatingMode() == battery.OPERATING_MODE_MANUAL {
			fmt.Println(`EVCC is NOT power charging, but battery disabled - enable battery`)
			err := e.battery.SetOperatingModeAutomatic()
			if err != nil {
				log.Fatal(err.Error())
			}
		}

		time.Sleep(10 * time.Second)
	}
}
