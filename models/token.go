package models

type MarketData struct {
	ATH                   float64 `json:"ath"`
	ATL                   float64 `json:"atl"`
	Change24H             float64 `json:"change_24h"`
	ChangePercentage24H   float64 `json:"change_percentage_24h"`
	ChangePercentage7D    float64 `json:"change_percentage_7d"`
	ChangePercentage14D   float64 `json:"change_percentage_14d"`
	ChangePercentage30D   float64 `json:"change_percentage_30d"`
	InitSupply            float64 `json:"init_supply"`
	MaxSupply             float64 `json:"max_supply"`
	TotalSupply           float64 `json:"total_supply"`
	CurrentSupply         float64 `json:"current_supply"`
	DailyVolume           float64 `json:"daily_volume"`
	MarketCap             float64 `json:"market_cap"`
	FullyDilutedValuation float64 `json:"fully_diluted_valuation"`
	CurrentLiquidity      float64 `json:"current_liquidity"`
	ZilReserve            float64 `json:"zil_reserve"`
	TokenReserve          float64 `json:"token_reserve"`
}

type TokenDetail struct {
	Name                string     `json:"name"`
	Symbol              string     `json:"symbol"`
	AddressBech32       string     `json:"address_bech32"`
	Icon                string     `json:"icon"`
	Decimals            float64    `json:"decimals"`
	Website             string     `json:"website"`
	Whitepaper          string     `json:"whitepaper"`
	Listed              bool       `json:"listed"`
	CurrentSupply       float64    `json:"current_supply"`
	DailyVolume         float64    `json:"daily_volume"`
	SupplySkipAddresses string     `json:"supply_skip_addresses"`
	MarketCap           float64    `json:"market_cap"`
	Rate                float64    `json:"rate"`
	RateUSD             float64    `json:"rate_usd"`
	MarketData          MarketData `json:"market_data"`
}
