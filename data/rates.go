package data

import (
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

type ExchangeRates struct {
	rates map[string]float64
}

func NewRates() (*ExchangeRates, error) {
	er := &ExchangeRates{rates: map[string]float64{}}
	err := er.getRates()
	return er, err
}

func (e *ExchangeRates) GetRate(base, dest string) (float64, error) {
	br, ok := e.rates[base]
	if !ok {
		return 0, fmt.Errorf("Rate not found for currency %s", base)
	}
	dr, ok := e.rates[dest]
	if !ok {
		return 0, fmt.Errorf("Rate not found for currency %s", dest)
	}
	return dr / br, nil
}

func (er ExchangeRates) getRates() error {
	resp, err := http.DefaultClient.Get("https://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml")
	if err != nil {
		fmt.Println("error in sending http to europe bank")
		return errors.New("error in sending http to europe bank")
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Println("error in sending http to europe bank, bad status code", resp.StatusCode)
		return errors.New(fmt.Sprintf("error in sending http to europe bank, bad status code =%d", resp.StatusCode))
	}

	md := &Cubes{}
	xml.NewDecoder(resp.Body).Decode(&md)

	for _, c := range md.CubeData {
		r, err := strconv.ParseFloat(c.Rate, 64)
		if err != nil {
			return err
		}
		er.rates[c.Currency] = r
	}
	er.rates["EUR"] = 1

	return nil
}

type Cubes struct {
	CubeData []Cube `xml:"Cube>Cube>Cube"`
}

type Cube struct {
	Currency string `xml:"currency,attr"`
	Rate     string `xml:"rate,attr"`
}
