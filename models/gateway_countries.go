package models

type GatewayCountry struct {
	GatewayID int `db:"gateway_id"`
	CountryID int `db:"country_id"`
	Priority  int `db:"priority"`
}
