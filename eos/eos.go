package eos

import (
	"encoding/json"
	"eosclient/logger"
	"errors"
	"github.com/eoscanada/eos-go"
	"strings"
)

type EosClient struct {
	cli      *eos.API
	keosdCli *eos.API
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
	logger.Debug(string(data))

	return nil
}

func (e *EosClient) GetABI(name string) error {
	abi, err := e.cli.GetABI(eos.AN(name))
	if err != nil {
		logger.Error(err)
		return err
	}

	data, _ := json.Marshal(abi)
	logger.Debug(string(data))

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
	args["memo"] = "txnotes"

	txData, err := e.cli.ABIJSONToBin(eos.AN(code), eos.Name(action), args)
	if err != nil {
		logger.Error(err)
		return err
	}

	logger.Info(txData.String())

	txDatajson, err := e.cli.ABIBinToJSON(eos.AN(code), eos.Name(action), txData)
	if err != nil {
		logger.Error(err)
		return err
	}

	logger.Info(txDatajson)

	info, err := e.cli.GetInfo()
	if err != nil {
		logger.Error(err)
		return err
	}

	data, _ := json.Marshal(info)
	logger.Info(string(data))

	block, err := e.cli.GetBlockByNum(info.HeadBlockNum)
	if err != nil {
		logger.Error(err)
		return err
	}

	data, _ = json.Marshal(block)
	logger.Info(string(data))

	walletList, err := e.keosdCli.ListWallets()
	if err != nil {
		logger.Error(err)
		return err
	}

	logger.Info("wallet list:", walletList) //sunlight *
	if len(walletList) < 1 {
		logger.Error("not exist wallet")
		return errors.New("not exist wallet")
	}

	ss := strings.Split(walletList[0], "")

	if len(ss) > 1 && ss[len(ss)-1] == "*" {
		logger.Info(walletList[0], "unlock already")
	} else {
		err = e.keosdCli.WalletUnlock(walletList[0], "PW5HtE9u7i4Dhwk1WJ8XVeaDNrCN1U3KDCKQ1TfrGe9vAYg4xWsMU")
		if err != nil {
			logger.Error(err)
			return err
		}
	}

	var actions []*eos.Action
	actionData := eos.NewActionData(txData)
	permissionLevel, err := eos.NewPermissionLevel(from)
	if err != nil {
		logger.Error(err)
		return err
	}

	ac := &eos.Action{
		Account:       eos.AN(code),
		Name:          eos.ActN(action),
		Authorization: []eos.PermissionLevel{permissionLevel},
		ActionData:    actionData,
	}
	actions = append(actions, ac)

	opts := &eos.TxOptions{
		HeadBlockID:      block.ID,
		MaxNetUsageWords: 0,
		DelaySecs:        0,
		MaxCPUUsageMS:    0,
	}

	tx := eos.NewTransaction(actions, opts)
	data, _ = json.Marshal(tx)
	logger.Info(string(data))

	e.cli.SetSigner(eos.NewWalletSigner(e.keosdCli, ""))
	requiredKeys, err := e.cli.GetRequiredKeys(tx)
	if err != nil {
		logger.Error(err)
		return err
	}

	if len(requiredKeys.RequiredKeys) < 1 {
		logger.Error("not required key")
		return errors.New("not required key")
	}

	data, _ = json.Marshal(requiredKeys)
	logger.Info(string(data))

	signTx := eos.NewSignedTransaction(tx)
	signedData, err := e.keosdCli.WalletSignTransaction(signTx, info.ChainID, requiredKeys.RequiredKeys[0])
	if err != nil {
		logger.Error(err)
		return err
	}

	signTx.Signatures = signedData.Signatures
	data, _ = json.Marshal(signTx)
	logger.Info(string(data))

	packedTx, err := signTx.Pack(eos.CompressionNone)
	if err != nil {
		logger.Error(err)
		return err
	}

	data, _ = json.Marshal(packedTx)
	logger.Info(string(data))

	txout, err := e.cli.PushTransaction(packedTx)
	if err != nil {
		logger.Error(err)
		return err
	}

	data, _ = json.Marshal(txout)
	logger.Info(string(data))

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

func NewEosClient(ipport, keosdIpport string) (*EosClient, error) {
	e := new(EosClient)
	e.cli = eos.New("http://" + ipport)
	e.keosdCli = eos.New("http://" + keosdIpport)
	return e, nil
}
