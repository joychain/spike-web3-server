package txService

import (
	"github.com/spike-engine/spike-web3-server/config"
	"github.com/spike-engine/spike-web3-server/global"
	"github.com/spike-engine/spike-web3-server/model"
	"github.com/spike-engine/spike-web3-server/request"
	"time"
)

type TxService struct {
}

var TxSrv = new(TxService)

func (t *TxService) RechargeToken(service request.RechargeTokenService) error {
	tx := model.SpikeTx{
		OrderId:         service.OrderId,
		From:            service.FromAddress,
		To:              config.Cfg.Contract.GameVaultAddress,
		Amount:          service.Amount,
		ContractAddress: service.ContractAddress,
		TxHash:          service.TxHash,
		Cb:              service.Cb,
		CreateTime:      time.Now().UnixMilli(),
	}
	return global.DbAccessor.SaveTxCb(tx)
}

func (t *TxService) ImportNft(service request.ImportNftService) error {
	tx := model.SpikeTx{
		OrderId:         service.OrderId,
		From:            service.From,
		To:              config.Cfg.Contract.GameVaultAddress,
		TokenId:         service.TokenId,
		ContractAddress: service.ContractAddress,
		TxHash:          service.TxHash,
		Cb:              service.Cb,
		CreateTime:      time.Now().UnixMilli(),
	}
	return global.DbAccessor.SaveTxCb(tx)
}
