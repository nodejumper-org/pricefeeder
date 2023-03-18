package sources

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/NibiruChain/nibiru/x/common/set"
	"github.com/NibiruChain/pricefeeder/types"
)

const (
	Okex = "okex"
)

var _ types.FetchPricesFunc = OkexPriceUpdate

type OkexTicker struct {
	Symbol string `json:"instId"`
	Price  string `json:"last"`
}

type Response struct {
	Data []OkexTicker `json:"data"`
}

// OkexPriceUpdate returns the prices for given symbols or an error.
// Uses OKEX API at https://www.okx.com/docs-v5/en/#rest-api-market-data.
func OkexPriceUpdate(symbols set.Set[types.Symbol]) (rawPrices map[types.Symbol]float64, err error) {

	url := "https://www.okx.com/api/v5/market/tickers?instType=SPOT"

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response Response
	err = json.Unmarshal(b, &response)
	if err != nil {
		return nil, err
	}

	rawPrices = make(map[types.Symbol]float64)
	for _, ticker := range response.Data {

		symbol := types.Symbol(strings.Replace(ticker.Symbol, "-", "", -1))
		price, err := strconv.ParseFloat(ticker.Price, 64)
		if err != nil {
			return rawPrices, err
		}

		if _, ok := symbols[symbol]; ok {
			rawPrices[symbol] = price
		}

	}
	return rawPrices, nil
}