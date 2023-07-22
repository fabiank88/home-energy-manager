package evcc

import (
	"encoding/json"
	"io"
	"net/http"
)

type Config struct {
	Ip   string
	Port string
}

type Evcc struct {
	config Config
	state  state
}

type state struct {
	Result struct {
		Loadpoints []struct {
			Mode          string
			ChargeCurrent int
		}
	}
}

func New(config Config) (Evcc, error) {
	return Evcc{config: config}, nil
}

func (e *Evcc) Refresh() error {
	return e.callState()
}

func (e *Evcc) IsPowerCharging() bool {
	for _, loadpoint := range e.state.Result.Loadpoints {
		if loadpoint.Mode == `now` || (loadpoint.Mode == `pv` && loadpoint.ChargeCurrent == 16) {
			return true
		}
	}
	return false
}

func (e *Evcc) callState() error {
	response, err := http.Get(`http://` + e.config.Ip + `:` + e.config.Port + `/api/state`)

	if err != nil {
		return err
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(responseData, &e.state)
	if err != nil {
		return err
	}
	return nil
}
