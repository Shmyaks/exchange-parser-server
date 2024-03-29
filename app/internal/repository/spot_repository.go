package repository

import (
	"context"

	"github.com/Shmyaks/exchange-parser-server/app/internal/data"
	"github.com/Shmyaks/exchange-parser-server/app/internal/models"
	"github.com/Shmyaks/exchange-parser-server/app/internal/models/markets"
	pkg "github.com/Shmyaks/exchange-parser-server/app/pkg/redis"

	jsoniter "github.com/json-iterator/go"
)

// SPOTRepository struct
type SPOTRepository struct {
	datas           []data.SPOT
	redisConnection pkg.Connection
}

// NewSPOTRepository fabric for CurrencyRepository
func NewSPOTRepository(curencyDatas []data.SPOT, redisConnection pkg.Connection) *SPOTRepository {
	mpCurrency := make([]data.SPOT, len(curencyDatas))
	for _, curencyData := range curencyDatas {
		mpCurrency[*curencyData.GetMarketID()-1] = curencyData
	}

	return &SPOTRepository{datas: mpCurrency, redisConnection: redisConnection}
}

// GetData get spot Data of repository
func (r *SPOTRepository) GetData(m markets.SPOTMarket) data.SPOT {
	if len(r.datas) <= int(m)-1 {
		panic("Datas not have this market")
	}
	return r.datas[m-1]
}

// GetPairNamesFromData method for get P2POrders from Data
func (r *SPOTRepository) GetPairNamesFromData(market markets.SPOTMarket) ([]models.BaseCurrencyPair, error) {
	return r.GetData(market).GetAllPairsAPI()
}

// GetAllFromData method for get P2POrders from Data
func (r *SPOTRepository) GetAllFromData(market markets.SPOTMarket) ([]models.CurrencyPair, error) {
	return r.GetData(market).GetDetailPairsAPI()
}

// Set method: SetCurrencyPair pair to Redis
func (r *SPOTRepository) Set(market markets.SPOTMarket, pair models.CurrencyPair) error {
	bs, _ := jsoniter.Marshal(&pair)
	cmd := r.redisConnection.Pool.HSet(context.Background(), market.GetName(), string(pair.CurencyPairName), bs)
	return cmd.Err()
}

// SetMany method: SetCurrencyPair pair to Redis
func (r *SPOTRepository) SetMany(market markets.SPOTMarket, pairs []models.CurrencyPair) error {
	mp := make(map[string]interface{})
	for _, pair := range pairs {
		bs, _ := jsoniter.Marshal(&pair)
		mp[string(pair.CurencyPairName)] = bs
	}
	cmd := r.redisConnection.Pool.HMSet(context.Background(), market.GetName(), mp)
	return cmd.Err()
}

// GetFromCache method: Set P2Ppair to Redis
func (r *SPOTRepository) GetFromCache(market markets.SPOTMarket, currencyPairName models.CurencyPairName) *models.CurrencyPair {
	pair := new(models.CurrencyPair)
	str, _ := r.redisConnection.Pool.HGet(context.Background(), market.GetName(), string(currencyPairName)).Result()
	_ = jsoniter.UnmarshalFromString(str, &pair)
	return pair
}

// GetAllFromCache SPOT pairs from market
func (r *SPOTRepository) GetAllFromCache(market markets.SPOTMarket) []models.CurrencyPair {
	pairs := []models.CurrencyPair{}
	mp, _ := r.redisConnection.Pool.HGetAll(context.Background(), market.GetName()).Result()
	for _, mp := range mp {
		pair := new(models.CurrencyPair)
		_ = jsoniter.UnmarshalFromString(mp, pair)
		pairs = append(pairs, *pair)
	}
	return pairs
}
