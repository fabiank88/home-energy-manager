package supplier

type HourlyPrice struct {
	Total  float64
	Energy int `json:"hour"`
	Tax
	StartsAt
	EndsAt
}

type Prices struct {
	prices []HourlyPrice
}

func (p Prices) GetCurrentPrice() HourlyPrice {

}
