package eos

import (
	"encoding/json"
	"eosclient/logger"
	"github.com/eoscanada/eos-go"
)

type EosClient struct {
	cli *eos.API
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

func (e *EosClient) GetHeadBlockNumber() (uint32, error) {
	info, err := e.cli.GetInfo()
	if err != nil {
		logger.Error(err)
		return 0, err
	}

	logger.Info("Head block number:", info.HeadBlockNum)

	return info.HeadBlockNum, nil
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

func (e *EosClient) GetABI(name string) error {
	abi, err := e.cli.GetABI(eos.AN(name))
	if err != nil {
		logger.Error(err)
		return err
	}

	data, _ := json.Marshal(abi)
	logger.Info(string(data))

	return nil
}

func (e *EosClient) GetCode(account string) error {
	code, err := e.cli.GetCode(eos.AN(account))
	if err != nil {
		logger.Error(err)
		return err
	}

	data, _ := json.Marshal(code)
	logger.Debug(string(data))

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

func (e *EosClient) GetDBSize() error {
	db, err := e.cli.GetDBSize()
	if err != nil {
		logger.Error(err)
		return err
	}

	data, _ := json.Marshal(db)
	logger.Info(string(data))

	return nil
}

func (e *EosClient) GetTransaction(id string) error {
	tx, err := e.cli.GetTransaction(id)
	if err != nil {
		logger.Error(err)
		return err
	}

	data, _ := json.Marshal(tx)
	logger.Info(string(data))

	return nil
}

func (e *EosClient) Transfer(code, action, from, to, num string) error {
	args := make(eos.M)
	args["from"] = from
	args["to"] = to
	args["quantity"] = num
	args["memo"] = "tx notes"

	txData, err := e.cli.ABIJSONToBin(eos.AN(code), eos.Name(action), args)
	if err != nil {
		logger.Error(err)
		return err
	}

	logger.Info(txData.String())

	info, err := e.cli.GetInfo()
	if err != nil {
		logger.Error(err)
		return err
	}

	logger.Info("Head block number:", info.HeadBlockNum)

	block, err := e.cli.GetBlockByNum(info.HeadBlockNum)
	if err != nil {
		logger.Error(err)
		return err
	}

	logger.Info("Head block:", block.Timestamp, block.BlockNum, block.RefBlockPrefix)

	/*	err = e.cli.WalletUnlock("*", "PW5Jpknq77E6U6boTdSK1BpUiEAm9m1LKnzo1ths9FAq3QobzKTky")
		if err != nil {
			logger.Error(err)
			return err
		}*/

	var actions []*eos.Action

	actionData := eos.ActionData{
		Data: txData,
	}

	ac := &eos.Action{
		Account:       eos.AN(code),
		Name:          eos.ActN(action),
		Authorization: []PermissionLevel{{eos.AN(alice), eos.PN("active")}},
		ActionData:    actionData,
	}

	actions = append(actions, ac)

	opts := &eos.TxOptions{}

	tx := eos.NewTransaction(actions, opts)
	requiredKeys, err := e.cli.GetRequiredKeys(tx)
	if err != nil {
		logger.Error(err)
		return err
	}

	/*	signedTx, err := e.cli.WalletSignTransaction()
		if err != nil {
			logger.Error(err)
			return err
		}

		txout, err := e.cli.PushTransaction(nil)
		if err != nil {
			logger.Error(err)
			return err
		}

		data, _ := json.Marshal(txout)
		logger.Info(string(data))*/

	return nil
}

func (e *EosClient) DeployContract(privateKeyHex string) error {
	return nil
}

func (e *EosClient) InvokeContract(richPrivKeyHex, contractAddrHex string) error {
	return nil
}

func (e *EosClient) GetNewWallet() error {
	return nil
}

func NewEosClient(ipport string) (*EosClient, error) {
	e := new(EosClient)
	e.cli = eos.New("http://" + ipport)
	return e, nil
}
