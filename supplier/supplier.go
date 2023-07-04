package supplier

import "errors"

type Config struct {
	Type  string
	Token string
}

type Supplier interface {
	GetHourlyPrices() Prices
}

func New(config Config) (Supplier, error) {
	if config.Type == `tibber` {
		return &tibber{
			config: config,
		}, nil
	} else {
		return nil, errors.New(`no such type found ` + config.Type)
	}
}
