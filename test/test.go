package test

import (
	"eosclient/eos"
	"eosclient/logger"
	"sync"
)

type EosClientTest struct {
	eosIpPort string
	cli       *eos.EosClient
}

func (e *EosClientTest) testApi() error {
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

	/*	err = cli.GetInfo()
		if err != nil {
			logger.Error(err)
			return
		}

		err = cli.GetAccount("alice")
		if err != nil {
			logger.Error(err)
			return
		}

		err = cli.GetBlockByNum(1000)
		if err != nil {
			logger.Error(err)
			return
		}

		err = cli.GetProducers()
		if err != nil {
			logger.Error(err)
			return
		}

		err = cli.GetBalance("bob", "SYS", "eosio.token")
		if err != nil {
			logger.Error(err)
			return
		}

		err = cli.GetCode("eosio.token")
		if err != nil {
			logger.Error(err)
			return
		}*/
}

func NewEosClientTest(ipPort string, wg *sync.WaitGroup) (*EosClientTest, error) {
	test := new(EosClientTest)
	test.eosIpPort = ipPort

	cli, err := eos.NewEosClient(test.eosIpPort)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	test.cli = cli

	go test.testing(wg)

	return test, nil
}
