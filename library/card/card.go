package lib_card

import (
	"errors"
	"utility/helper"
)

type CardApiResponse struct {
	Status   string      `json:"Status"`
	Scheme   string      `json:"Scheme"`
	Type     string      `json:"Type"`
	Issuer   string      `json:"Issuer"`
	CardTier string      `json:"CardTier"`
	Country  CardCountry `json:"Country"`
	Luhn     bool        `json:"Luhn"`
}

type CardCountry struct {
	A2   string `json:"A2"`
	A3   string `json:"A3"`
	N3   string `json:"N3"`
	ISD  string `json:"ISD"`
	Name string `json:"Name"`
	Cont string `json:"Cont"`
}

type DataCardDetail struct {
	Bincard   string `json:"bincard"`
	Brand     string `json:"brand"`
	Type      string `json:"type"`
	Cardtier  string `json:"cardtier"`
	Issuer    string `json:"issuer"`
	Iso_code2 string `json:"country_iso_code2"`
	Iso_code3 string `json:"country_iso_code3"`
	Country   string `json:"country"`
}

func GetCardDetail(card_number string) (DataCardDetail, error) {

	if card_number == "" {
		return DataCardDetail{}, nil
	}

	body, status, err := helper.HttpClientGet(
		// "https://data.handyapi.com/bin/535316",
		"https://data.handyapi.com/bin/"+card_number,
		map[string]string{
			"Authorization": "Bearer token",
		},
	)

	if err != nil {
		return DataCardDetail{}, err
	}

	if status != 200 {
		return DataCardDetail{}, errors.New("status header res not 200 !")
	}

	resp, err := helper.DecodeResponse[CardApiResponse](body)
	if err != nil {
		return DataCardDetail{}, err
	}

	DataCardDetail := DataCardDetail{
		Bincard:   card_number,
		Brand:     resp.Scheme,
		Type:      resp.Type,
		Cardtier:  resp.CardTier,
		Issuer:    resp.Issuer,
		Iso_code2: resp.Country.A2,
		Iso_code3: resp.Country.A3,
		Country:   resp.Country.Name,
	}

	return DataCardDetail, nil

}
