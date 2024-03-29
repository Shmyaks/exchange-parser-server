// Package bybit methods file
package bybit

import (
	"strconv"
	"strings"

	"github.com/Shmyaks/exchange-parser-server/app/internal/data/bybit/schemes"
	"github.com/Shmyaks/exchange-parser-server/app/internal/models"
	"github.com/Shmyaks/exchange-parser-server/app/internal/models/filters"
	"github.com/Shmyaks/exchange-parser-server/app/internal/models/markets"
	market "github.com/Shmyaks/exchange-parser-server/app/internal/models/markets"

	"github.com/imroc/req/v3"
)

const (
	p2pURL      string = "https://api2.bybit.com/spot/api/otc/item/list"
	p2pAliasUrl string = "https://api2.bybit.com/spot/api/v1/otc/payment/broker_config_list"
)

// P2PData struct
type P2PData struct {
	client    req.Client
	P2PMarket market.P2PMarket
}

// NewP2PData fabric for Data bybit
func NewP2PData(client req.Client) *P2PData {
	return &P2PData{client: client, P2PMarket: markets.BybitP2P}
}

// GetMarketID get p2p market id
func (d *P2PData) GetMarketID() *market.P2PMarket {
	return &d.P2PMarket
}

// GetOrdersAPI get p2p orders from bybit
func (d *P2PData) GetOrdersAPI(filter filters.P2PFilter) ([]models.P2POrder, error) {
	orders := make([]models.P2POrder, 0, 10)
	var scheme schemes.P2PJSONScheme
	body := map[string]string{
		"tokenId":    string(filter.CryptoCurrency),
		"currencyId": string(filter.Fiat),
		"payment":    mapAliasPayMethod[filter.PayType],
		"side":       mapAliasFilterType[filter.TradeType],
		"size":       "10",
		"page":       "1",
		"amount":     strconv.Itoa(filter.MinAmount),
	}
	resp, err := d.client.NewRequest().SetFormData(body).Post(p2pURL)
	if err != nil {
		return nil, err
	}
	if !resp.IsSuccess() {
		return nil, nil
	}
	err = resp.Into(&scheme)
	if err != nil {
		return nil, nil
	}
	for _, info := range scheme.Result.Items {
		orders = append(
			orders,
			*models.NewP2POrder(info.TokenID, info.Fiat, info.NickName, info.Price, d.P2PMarket, filter.PayType),
		)
	}
	return orders, nil
}

// GetPayMethods get P2P PayMethods for bybit
func (d *P2PData) GetPayMethods() map[models.Fiat][]models.PayMethod {
	mp := make(map[models.Fiat][]models.PayMethod)
	mp[models.RUB] = []models.PayMethod{
		models.Tinkoff, models.Raif, models.Rosbank, models.QIWI, models.MTSBank, models.YOUMoney, models.Payeer,
		models.AdvCash,
	}
	mp[models.TRY] = []models.PayMethod{models.Ziraat}
	mp[models.USD] = []models.PayMethod{models.Payeer, models.AdvCash, models.Wise, models.Revolut, models.Zelle}
	return mp
}

func (d *P2PData) GetPayMethodAlias() error {
	var scheme schemes.AliasJSONScheme
	mp := make(map[models.PayMethod]string)
	resp, err := d.client.NewRequest().Post(p2pAliasUrl)
	if err != nil {
		return err
	}
	if !resp.IsSuccess() {
		return nil
	}
	err = resp.Into(&scheme)
	for _, info := range scheme.Result {
		mp[models.PayMethod(
			strings.ToUpper(
				strings.ReplaceAll(
					info.PaymentName, " ", "",
				),
			),
		)] = strconv.Itoa(info.PaymentType)
	}
	mapAliasPayMethod = mp
	return nil
}
