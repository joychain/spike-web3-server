package queryService

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	logger "github.com/ipfs/go-log"
	"golang.org/x/xerrors"
	"spike-frame/config"
	"spike-frame/constant"
	"spike-frame/response"
	"spike-frame/util"
)

var log = logger.Logger("service")

var MoralisRateLimit = "{\"message\":\"Rate limit exceeded.\"}"

func QueryWalletNft(cursor, walletAddr, network string, res []response.NftResult) ([]response.NftResult, error) {
	client := resty.New()
	resp, err := client.R().
		SetHeader("Accept", "application/json").
		SetHeader("x-api-key", config.Cfg.Moralis.XApiKey).
		Get(getUrl(config.Cfg.Contract.GameNftAddress, walletAddr, network, cursor))
	if err != nil {
		log.Errorf("query wallet nft , wallet : %s, err : %+v", walletAddr, err)
		return res, err
	}
	if string(resp.Body()) == MoralisRateLimit {
		log.Errorf(MoralisRateLimit)
		return res, xerrors.New(MoralisRateLimit)
	}
	var nrs response.NftResults
	err = json.Unmarshal(resp.Body(), &nrs)
	if err != nil {
		log.Errorf("json unmarshal err : %+v", err)
		return res, err
	}
	res = append(res, nrs.Results...)
	if nrs.Page*nrs.PageSize >= nrs.Total {
		return res, nil
	}
	res, err = QueryWalletNft(nrs.Cursor, walletAddr, network, res)
	return res, nil
}

func getUrl(contractAddr, walletAddr, network, cursor string) string {
	return fmt.Sprintf("%s%s/nft/%s?chain=%s&cursor=%s", constant.MORALIS_API, walletAddr, contractAddr, network, cursor)
}

func (qm *QueryManager) handleNftData(walletAddr string, data []response.NftResult) ([]response.NftResult, error) {
	data = util.ConvertNftResult(data)
	dataList := util.ParseMetadata(data)
	dataMap := util.ParseCacheData(dataList)
	nftType := make([]response.NftType, 0)

	for k, _ := range dataMap {
		nftType = append(nftType, response.NftType{
			Name:   k,
			Amount: len(dataMap[k]),
		})
		cacheByte, err := json.Marshal(dataMap[k])
		if err != nil {
			break
		}
		util.SetFromRedis(walletAddr+constant.NFTLISTSUFFIX+k, string(cacheByte), nftListDuration, qm.redisClient)
	}

	nftTypeByte, err := json.Marshal(nftType)
	if err != nil {
		return data, err
	}
	util.SetFromRedis(walletAddr+constant.NFTTYPESUFFIX, string(nftTypeByte), nftListDuration, qm.redisClient)
	return data, nil
}
