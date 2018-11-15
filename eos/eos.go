package eos

import (
	"encoding/json"
	"eosclient/logger"
	"github.com/eoscanada/eos-go"
)

type EosClient struct {
	cli *eos.API
}

func (e *EosClient) DeployContract(privateKeyHex string) error {
	return nil
}

func (e *EosClient) InvokeContract(richPrivKeyHex, contractAddrHex string) error {
	return nil
}

func (e *EosClient) GetBalance(account, symbol, code string) error {
	assets, err := e.cli.GetCurrencyBalance(eos.AN(account), symbol, eos.AN(code))
	if err != nil {
		logger.Error(err)
		return err
	}

	data, _ := json.Marshal(assets)
	logger.Info(string(data))

	return nil
}

func (e *EosClient) GetNewWallet() error {
	return nil
}

func (e *EosClient) BlockByNumber(num int64) error {
	return nil
}

func (e *EosClient) QueryTransaction(txHex string) error {
	return nil
}

func (e *EosClient) GetCode(account string) error {
	code, err := e.cli.GetCode(eos.AN(account))
	if err != nil {
		logger.Error(err)
		return err
	}

	data, _ := json.Marshal(code)
	logger.Info(string(data))

	return nil
}

func (e *EosClient) Transfer(from, to, num string) error {
	tx, err := e.cli.PushTransaction(nil)
	if err != nil {
		logger.Error(err)
		return err
	}

	data, _ := json.Marshal(tx)
	logger.Info(string(data))

	return nil
}

func (e *EosClient) Transfer2(from, to, num, fromPriv string) error {
	return nil
}

func (e *EosClient) GetInfo() error {
	info, err := e.cli.GetInfo()
	if err != nil {
		logger.Error(err)
		return err
	}
	data, _ := json.Marshal(info)
	logger.Info(string(data))

	return nil
}

func (e *EosClient) GetAccount(name string) error {
	account, err := e.cli.GetAccount(eos.AN(name))
	if err != nil {
		logger.Error(err)
		return err
	}

	data, _ := json.Marshal(account)

	logger.Info(string(data))
	return nil
}

func (e *EosClient) GetBlockByNum(blockNumber uint32) error {
	block, err := e.cli.GetBlockByNum(blockNumber)
	if err != nil {
		logger.Error(err)
		return err
	}

	data, _ := json.Marshal(block)

	logger.Info(string(data))
	return nil
}

func (e *EosClient) GetProducers() error {
	producers, err := e.cli.GetProducers()
	if err != nil {
		logger.Error(err)
		return err
	}

	data, _ := json.Marshal(producers)

	logger.Info(string(data))
	return nil
}

func NewEosClient(ipport string) (*EosClient, error) {
	e := new(EosClient)
	e.cli = eos.New("http://" + ipport)
	return e, nil
}
