package test

import (
	"eosclient/eos"
	"eosclient/logger"
	"eosclient/util"
	"sync"
	"time"
)

type EosClientTest struct {
	cli                  *eos.EosClient
	tokenContractAccount string
	richAccount          string
	mainAccount          string
	tokenSymbol          string
}

func (e *EosClientTest) testTx() error {
	var err error
	var txid string

	//create token
	e.tokenSymbol = util.GetRandomUpperString(4)
	txid, err = e.cli.CreateToken(e.tokenContractAccount, e.mainAccount, "99999999.000000", e.tokenSymbol)
	if err != nil {
		logger.Error(err)
		return nil
	}
	time.Sleep(time.Second * 1)

	//issue token
	txid, err = e.cli.IssueToken(e.tokenContractAccount, e.mainAccount, e.richAccount, "4500.000000", e.tokenSymbol)
	if err != nil {
		logger.Error(err)
		return nil
	}
	time.Sleep(time.Second * 1)

	e.cli.GetBalance(e.richAccount, e.tokenSymbol, e.tokenContractAccount)

	//new account
	newAccount := "account" + util.GetRandomLowerString(4)
	_, pubKey, _ := e.cli.GetNewKey()
	txid, err = e.cli.GetNewAccount(e.richAccount, newAccount, pubKey)
	if err != nil {
		logger.Error(err)
		return nil
	}
	time.Sleep(time.Second * 1)

	e.cli.GetBalance(newAccount, e.tokenSymbol, e.tokenContractAccount)

	//transfer
	txid, err = e.cli.Transfer(e.tokenContractAccount, "transfer", e.richAccount, newAccount, "2.000000 "+e.tokenSymbol)
	if err != nil {
		logger.Error(err)
		return nil
	}
	time.Sleep(time.Second * 2)
	e.cli.GetTransaction(txid)

	e.cli.GetBalance(e.richAccount, e.tokenSymbol, e.tokenContractAccount)
	e.cli.GetBalance(newAccount, e.tokenSymbol, e.tokenContractAccount)

	return nil
}

func (e *EosClientTest) testApi() error {
	var err error

	/*	txid, err := e.cli.Transfer(e.tokenContractAccount, "transfer", e.richAccount, "bob", "2.0000 "+"SYS")
		if err != nil {
			logger.Error(err)
			return nil
		}
		e.cli.GetTransaction(txid)

		e.cli.GetBalance(e.richAccount, "SYS", e.tokenContractAccount)
		e.cli.GetBalance("bob", "SYS", e.tokenContractAccount)

		return nil*/

	err = e.testTx()
	if err != nil {
		logger.Error(err)
		return nil
	}
	return nil

	err = e.cli.GetInfo()
	if err != nil {
		logger.Error(err)
		return nil
	}

	headBlockNumber, err := e.cli.GetHeadBlockNumber()
	if err != nil {
		logger.Error(err)
		return nil
	}

	err = e.cli.GetBlockByNum(headBlockNumber)
	if err != nil {
		logger.Error(err)
		return nil
	}

	err = e.cli.GetAccount(e.richAccount)
	if err != nil {
		logger.Error(err)
		return nil
	}

	err = e.cli.GetABI(e.tokenContractAccount)
	if err != nil {
		logger.Error(err)
		return nil
	}

	err = e.cli.GetCode(e.tokenContractAccount)
	if err != nil {
		logger.Error(err)
		return nil
	}

	err = e.cli.GetActions(e.richAccount)
	if err != nil {
		logger.Error(err)
		return nil
	}

	/*	err = e.cli.GetProducers()
		if err != nil {
			logger.Error(err)
			return nil
		}*/

	/*	err = e.cli.GetDBSize()
		if err != nil {
			logger.Error(err)
			return nil
		}*/

	return nil
}

func (e *EosClientTest) testing(wg *sync.WaitGroup) {
	defer wg.Done()
	var err error

	err = e.testApi()
	if err != nil {
		logger.Error(err)
		return
	}
}

func NewEosClientTest(ipPort, keosdIpPort, walletPassword, tokenContractAccount, richAccount, mainAccount, tokenSymbol string, wg *sync.WaitGroup) (*EosClientTest, error) {
	test := new(EosClientTest)

	cli, err := eos.NewEosClient(ipPort, keosdIpPort, walletPassword)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	test.cli = cli
	test.tokenContractAccount = tokenContractAccount
	test.richAccount = richAccount
	test.mainAccount = mainAccount
	//test.tokenSymbol = tokenSymbol

	go test.testing(wg)

	return test, nil
}
