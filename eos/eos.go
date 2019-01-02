package eos

import (
	"encoding/json"
	"eosclient/logger"
	"errors"
	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/ecc"
	"github.com/eoscanada/eos-go/system"
	"strings"
)

type EosClient struct {
	cli      *eos.API
	keosdCli *eos.API

	walletPW string
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

	logger.Debug("Head block number:", info.HeadBlockNum)

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
	logger.Info(account, symbol, "balance:", string(data))

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

func (e *EosClient) GetActions(account string) error {
	req := eos.GetActionsRequest{
		AccountName: eos.AN(account),
		Pos:         -1,
		Offset:      -20,
	}

	actions, err := e.cli.GetActions(req)
	if err != nil {
		logger.Error(err)
		return err
	}

	data, _ := json.Marshal(actions)
	logger.Debug(string(data))

	return nil
}

func (e *EosClient) CreateToken(code, issuer, num, symbol string) (string, error) {
	action := "create"

	args := make(eos.M)
	args["issuer"] = issuer
	args["maximum_supply"] = num + " " + symbol

	txData, err := e.cli.ABIJSONToBin(eos.AN(code), eos.Name(action), args)
	if err != nil {
		logger.Error(err)
		return "", err
	}

	info, err := e.cli.GetInfo()
	if err != nil {
		logger.Error(err)
		return "", err
	}

	block, err := e.cli.GetBlockByNum(info.HeadBlockNum)
	if err != nil {
		logger.Error(err)
		return "", err
	}

	walletList, err := e.keosdCli.ListWallets()
	if err != nil {
		logger.Error(err)
		return "", err
	}

	logger.Debug("wallet list:", walletList) //sunlight *
	if len(walletList) < 1 {
		logger.Error("not exist wallet")
		return "", errors.New("not exist wallet")
	}

	ss := strings.Split(walletList[0], "")
	if len(ss) > 1 && ss[len(ss)-1] == "*" {
		logger.Debug(walletList[0], "unlock already")
	} else {
		err = e.keosdCli.WalletUnlock(walletList[0], e.walletPW)
		if err != nil {
			logger.Error(err)
			return "", err
		}
	}

	var actions []*eos.Action
	actionData := eos.NewActionDataFromHexData(txData)
	permissionLevel, err := eos.NewPermissionLevel(code)
	if err != nil {
		logger.Error(err)
		return "", err
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

	e.cli.SetSigner(eos.NewWalletSigner(e.keosdCli, ""))
	requiredKeys, err := e.cli.GetRequiredKeys(tx)
	if err != nil {
		logger.Error(err)
		return "", err
	}

	if len(requiredKeys.RequiredKeys) < 1 {
		logger.Error("not required key")
		return "", errors.New("not required key")
	}

	signTx := eos.NewSignedTransaction(tx)
	signedData, err := e.keosdCli.WalletSignTransaction(signTx, info.ChainID, requiredKeys.RequiredKeys[0])
	if err != nil {
		logger.Error(err)
		return "", err
	}

	signTx.Signatures = signedData.Signatures
	data, _ := json.Marshal(signTx)
	logger.Debug(string(data))

	packedTx, err := signTx.Pack(eos.CompressionNone)
	if err != nil {
		logger.Error(err)
		return "", err
	}

	data, _ = json.Marshal(packedTx)
	logger.Debug(string(data))

	txRet, err := e.cli.PushTransaction(packedTx)
	if err != nil {
		logger.Error(err)
		return "", err
	}

	logger.Info("CreateToken successfully", issuer, num+" "+symbol)

	return txRet.TransactionID, nil
}

func (e *EosClient) IssueToken(code, from, to, num, symbol string) (string, error) {
	action := "issue"

	args := make(eos.M)
	args["to"] = to
	args["quantity"] = num + " " + symbol
	args["memo"] = "issue token " + num + " " + symbol + " to " + to

	txData, err := e.cli.ABIJSONToBin(eos.AN(code), eos.Name(action), args)
	if err != nil {
		logger.Error(err)
		return "", err
	}

	info, err := e.cli.GetInfo()
	if err != nil {
		logger.Error(err)
		return "", err
	}

	block, err := e.cli.GetBlockByNum(info.HeadBlockNum)
	if err != nil {
		logger.Error(err)
		return "", err
	}

	walletList, err := e.keosdCli.ListWallets()
	if err != nil {
		logger.Error(err)
		return "", err
	}

	logger.Debug("wallet list:", walletList) //sunlight *
	if len(walletList) < 1 {
		logger.Error("not exist wallet")
		return "", errors.New("not exist wallet")
	}

	ss := strings.Split(walletList[0], "")
	if len(ss) > 1 && ss[len(ss)-1] == "*" {
		logger.Debug(walletList[0], "unlock already")
	} else {
		err = e.keosdCli.WalletUnlock(walletList[0], e.walletPW)
		if err != nil {
			logger.Error(err)
			return "", err
		}
	}

	var actions []*eos.Action
	actionData := eos.NewActionDataFromHexData(txData)
	permissionLevel, err := eos.NewPermissionLevel(from)
	if err != nil {
		logger.Error(err)
		return "", err
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

	e.cli.SetSigner(eos.NewWalletSigner(e.keosdCli, ""))
	requiredKeys, err := e.cli.GetRequiredKeys(tx)
	if err != nil {
		logger.Error(err)
		return "", err
	}

	if len(requiredKeys.RequiredKeys) < 1 {
		logger.Error("not required key")
		return "", errors.New("not required key")
	}

	signTx := eos.NewSignedTransaction(tx)
	signedData, err := e.keosdCli.WalletSignTransaction(signTx, info.ChainID, requiredKeys.RequiredKeys[0])
	if err != nil {
		logger.Error(err)
		return "", err
	}

	signTx.Signatures = signedData.Signatures
	data, _ := json.Marshal(signTx)
	logger.Debug(string(data))

	packedTx, err := signTx.Pack(eos.CompressionNone)
	if err != nil {
		logger.Error(err)
		return "", err
	}

	data, _ = json.Marshal(packedTx)
	logger.Debug(string(data))

	txRet, err := e.cli.PushTransaction(packedTx)
	if err != nil {
		logger.Error(err)
		return "", err
	}

	logger.Info("IssueToken successfully", "issue token", num, symbol, "to", to)

	return txRet.TransactionID, nil
}

func (e *EosClient) Transfer(code, action, from, to, num string) (string, error) {
	args := make(eos.M)
	args["from"] = from
	args["to"] = to
	args["quantity"] = num
	args["memo"] = from + " transfer " + to + " " + num

	txData, err := e.cli.ABIJSONToBin(eos.AN(code), eos.Name(action), args)
	if err != nil {
		logger.Error(err)
		return "", err
	}
	logger.Debug(txData.String())

	info, err := e.cli.GetInfo()
	if err != nil {
		logger.Error(err)
		return "", err
	}

	data, _ := json.Marshal(info)
	logger.Debug(string(data))

	block, err := e.cli.GetBlockByNum(info.HeadBlockNum)
	if err != nil {
		logger.Error(err)
		return "", err
	}

	data, _ = json.Marshal(block)
	logger.Debug(string(data))

	walletList, err := e.keosdCli.ListWallets()
	if err != nil {
		logger.Error(err)
		return "", err
	}

	logger.Debug("wallet list:", walletList) //sunlight *
	if len(walletList) < 1 {
		logger.Error("not exist wallet")
		return "", errors.New("not exist wallet")
	}

	ss := strings.Split(walletList[0], "")
	if len(ss) > 1 && ss[len(ss)-1] == "*" {
		logger.Debug(walletList[0], "unlock already")
	} else {
		err = e.keosdCli.WalletUnlock(walletList[0], e.walletPW)
		if err != nil {
			logger.Error(err)
			return "", err
		}
	}

	var actions []*eos.Action
	actionData := eos.NewActionDataFromHexData(txData)
	permissionLevel, err := eos.NewPermissionLevel(from)
	if err != nil {
		logger.Error(err)
		return "", err
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
	logger.Debug(string(data))

	e.cli.SetSigner(eos.NewWalletSigner(e.keosdCli, ""))
	requiredKeys, err := e.cli.GetRequiredKeys(tx)
	if err != nil {
		logger.Error(err)
		return "", err
	}

	if len(requiredKeys.RequiredKeys) < 1 {
		logger.Error("not required key")
		return "", errors.New("not required key")
	}

	data, _ = json.Marshal(requiredKeys)
	logger.Debug(string(data))

	signTx := eos.NewSignedTransaction(tx)
	signedData, err := e.keosdCli.WalletSignTransaction(signTx, info.ChainID, requiredKeys.RequiredKeys[0])
	if err != nil {
		logger.Error(err)
		return "", err
	}

	signTx.Signatures = signedData.Signatures
	data, _ = json.Marshal(signTx)
	logger.Debug(string(data))

	packedTx, err := signTx.Pack(eos.CompressionNone)
	if err != nil {
		logger.Error(err)
		return "", err
	}

	data, _ = json.Marshal(packedTx)
	logger.Debug(string(data))

	txRet, err := e.cli.PushTransaction(packedTx)
	if err != nil {
		logger.Error(err)
		return "", err
	}

	data, _ = json.Marshal(txRet)
	logger.Debug(string(data))

	logger.Info("Transfer successfully")

	return txRet.TransactionID, nil
}

func (e *EosClient) GetNewKey() (string, string, error) {
	priv, err := ecc.NewRandomPrivateKey()
	if err != nil {
		logger.Error(err)
		return "", "", err
	}

	pub := priv.PublicKey()

	logger.Debug(priv.String(), pub.String())

	return priv.String(), pub.String(), nil
}

func (e *EosClient) GetNewAccount(creator, newAccount, pubKey string) (string, error) {
	pub, err := ecc.NewPublicKey(pubKey)
	if err != nil {
		logger.Error(err)
		return "", err
	}

	info, err := e.cli.GetInfo()
	if err != nil {
		logger.Error(err)
		return "", err
	}

	block, err := e.cli.GetBlockByNum(info.HeadBlockNum)
	if err != nil {
		logger.Error(err)
		return "", err
	}

	walletList, err := e.keosdCli.ListWallets()
	if err != nil {
		logger.Error(err)
		return "", err
	}

	logger.Debug("wallet list:", walletList) //sunlight *
	if len(walletList) < 1 {
		logger.Error("not exist wallet")
		return "", errors.New("not exist wallet")
	}

	ss := strings.Split(walletList[0], "")
	if len(ss) > 1 && ss[len(ss)-1] == "*" {
		logger.Debug(walletList[0], "unlock already")
	} else {
		err = e.keosdCli.WalletUnlock(walletList[0], e.walletPW)
		if err != nil {
			logger.Error(err)
			return "", err
		}
	}

	var actions []*eos.Action
	ac := system.NewNewAccount(eos.AN(creator), eos.AN(newAccount), pub)
	actions = append(actions, ac)

	opts := &eos.TxOptions{
		HeadBlockID:      block.ID,
		MaxNetUsageWords: 0,
		DelaySecs:        0,
		MaxCPUUsageMS:    0,
	}

	tx := eos.NewTransaction(actions, opts)

	e.cli.SetSigner(eos.NewWalletSigner(e.keosdCli, ""))
	requiredKeys, err := e.cli.GetRequiredKeys(tx)
	if err != nil {
		logger.Error(err)
		return "", err
	}

	if len(requiredKeys.RequiredKeys) < 1 {
		logger.Error("not required key")
		return "", errors.New("not required key")
	}

	signTx := eos.NewSignedTransaction(tx)
	signedData, err := e.keosdCli.WalletSignTransaction(signTx, info.ChainID, requiredKeys.RequiredKeys[0])
	if err != nil {
		logger.Error(err)
		return "", err
	}

	signTx.Signatures = signedData.Signatures
	data, _ := json.Marshal(signTx)
	logger.Debug(string(data))

	packedTx, err := signTx.Pack(eos.CompressionNone)
	if err != nil {
		logger.Error(err)
		return "", err
	}

	data, _ = json.Marshal(packedTx)
	logger.Debug(string(data))

	txRet, err := e.cli.PushTransaction(packedTx)
	if err != nil {
		logger.Error(err)
		return "", err
	}

	data, _ = json.Marshal(txRet)
	logger.Debug(string(data))

	logger.Info("GetNewAccount successfully", creator, newAccount, pubKey)

	return txRet.TransactionID, nil
}

func NewEosClient(ipport, keosdIpport, walletPW string) (*EosClient, error) {
	e := new(EosClient)
	e.cli = eos.New("http://" + ipport)
	e.keosdCli = eos.New("http://" + keosdIpport)
	e.walletPW = walletPW
	return e, nil
}
