package main

import (
	"eosclient/logger"
	"eosclient/test"
	"gopkg.in/ini.v1"
	"log"
	"sync"
)

const (
	EosClientConfFilePath = "conf/my.ini"
)

func initLogger() error {
	//fileHandler := logger.NewFileHandler("test.log")
	//logger.SetHandlers(logger.Console, fileHandler)
	logger.SetHandlers(logger.Console)
	//defer logger.Close()
	logger.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	logger.SetLevel(logger.INFO)

	return nil
}

func main() {
	err := initLogger()
	if err != nil {
		log.Fatalln(err)
	}

	cfg, err := ini.Load(EosClientConfFilePath)
	if err != nil {
		logger.Error(err)
		return
	}

	ipport := cfg.Section("").Key("EosServerIpPort").String()
	keosdIpPort := cfg.Section("").Key("KeosdIpPort").String()
	walletPassword := cfg.Section("").Key("WalletPassword").String()
	tokenContractAccount := cfg.Section("").Key("TokenContractAccount").String()
	richAccount := cfg.Section("").Key("RichAccount").String()
	mainAccount := cfg.Section("").Key("MainAccount").String()
	tokenSymbol := cfg.Section("").Key("TokenSymbol").String()

	var wg sync.WaitGroup
	wg.Add(1)

	_, err = test.NewEosClientTest(ipport, keosdIpPort, walletPassword, tokenContractAccount, richAccount, mainAccount, tokenSymbol, &wg)
	if err != nil {
		logger.Error(err)
		return
	}

	wg.Wait()
	logger.Debug("eos client exit")
}
