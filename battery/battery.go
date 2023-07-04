package battery

import (
	"errors"
)

const (
	OPERATING_MODE_MANUAL    = 1
	OPERATING_MODE_AUTOMATIC = 2
)

type Config struct {
	Name  string
	Type  string
	Ip    string
	Token string
}

type Battery interface {
	GetSoc() float32
	GetName() string
	GetOperatingMode() int
	SetOperatingModeAutomatic() error
	SetOperatingModeManual(watts int) error
	GetFlow() int
	Refresh() error
}

type battery struct {
	Battery
	config Config
}

func New(config Config) (Battery, error) {
	if config.Type == `sonnen` {
		return &sonnen{
			battery: battery{
				config: config,
			},
		}, nil
	} else {
		return nil, errors.New(`no such type found ` + config.Type)
	}
}

func (b *battery) GetName() string {
	return b.config.Name
}

func OperatingModeToString(mode int) string {
	switch mode {
	case OPERATING_MODE_MANUAL:
		return `manual`
	case OPERATING_MODE_AUTOMATIC:
		return `automatic`
	}
	return `unknown`
}
