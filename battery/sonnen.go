package battery

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type sonnen struct {
	battery
	status sonnenStatus
}

type sonnenStatus struct {
	USOC          float32
	OperatingMode int `json:"OperatingMode,string,omitempty"` //1 - manual, 2 - automatic
	Pac_total_W   int //current flow - means loading the battery
}

func (b sonnen) GetSoc() float32 {
	return b.status.USOC
}

func (b *sonnen) GetOperatingMode() int {
	return b.status.OperatingMode
}

func (b *sonnen) GetFlow() int {
	return b.status.Pac_total_W
}

func (b *sonnen) SetOperatingModeAutomatic() error {
	values := url.Values{}
	values.Add(`EM_OperatingMode`, strconv.Itoa(OPERATING_MODE_AUTOMATIC))
	return b.setConfiguration(values)
}

func (b *sonnen) SetOperatingModeManual(watts int) error {
	values := url.Values{}
	values.Add(`EM_OperatingMode`, strconv.Itoa(OPERATING_MODE_MANUAL))

	err := b.setConfiguration(values)
	if err != nil {
		return err
	}

	err = b.setCharge(watts)

	return err
}

func (b *sonnen) Refresh() error {
	response, err := http.Get(`http://` + b.config.Ip + `/api/v2/status`)

	if err != nil {
		return err
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(responseData, &b.status)
	if err != nil {
		return err
	}
	return nil
}

func (b *sonnen) setConfiguration(settings url.Values) error {
	req, err := http.NewRequest(http.MethodPut, `http://`+b.config.Ip+`/api/v2/configurations`, strings.NewReader(settings.Encode()))
	req.Header.Set("Auth-Token", b.config.Token)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		errorMessage := fmt.Sprintf(`error while setting battery configuration (%s)`, resp.Status)
		return errors.New(errorMessage)
	}
	defer resp.Body.Close()

	return nil
}

func (b *sonnen) setCharge(watts int) error {
	var endpoint string
	if watts > 0 {
		endpoint = `charge`
	} else {
		endpoint = `discharge`
		watts = watts * -1
	}

	req, err := http.NewRequest(http.MethodPost, `http://`+b.config.Ip+`/api/v2/setpoint/`+endpoint+`/`+strconv.Itoa(watts), bytes.NewBuffer([]byte(``)))
	req.Header.Set("Auth-Token", b.config.Token)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return err
	}
	if resp.StatusCode != 201 {
		errorMessage := fmt.Sprintf(`error while setting battery charge (%s)`, resp.Status)
		return errors.New(errorMessage)
	}

	return nil
}
