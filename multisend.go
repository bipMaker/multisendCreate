package main

import (
	// "encoding/hex"
	"fmt"
	"io/ioutil"

	// "math/rand"
	"os"
	// "path/filepath"
	"strings"
	"time"

	"strconv"

	"github.com/BurntSushi/toml"
	m "github.com/ValidatorCenter/minter-go-sdk"
)

var (
	conf     Config
	sdk      m.SDK
	oSndData []m.TxOneSendCoinData
)

type Config struct {
	AccAddress  string          `toml:"address_account"`
	AccKey      string          `toml:"privatkey"`
	ListAddress [][]interface{} `toml:"list_address"`
	Nonce       int             `toml:"nonce"`
	GasCoin     string          `toml:"gas_coin"`
}

func CreateFileTx(nonce int) error {
	msndDt := m.TxMultiSendCoinData{
		List:     oSndData,
		Payload:  "BipMaker+1% multisend",
		GasCoin:  conf.GasCoin,
		GasPrice: 1,
	}
	signstr, err := sdk.GetTxSignMultiSendCoin(&msndDt, nonce)
	if err != nil {
		fmt.Println("Ошибка:", err.Error())
		return err
	}
	TmpFileName := "txNonce" + strconv.Itoa(nonce) + ".bin"
	fmt.Println("FileName:", TmpFileName)
	err = ioutil.WriteFile(TmpFileName, signstr, 0777)
	if err != nil {
		// Если произошла ошибка выводим ее в консоль
		fmt.Println(err)
		return err
	}
	return nil
}

func main() {
	ConfFileName := "multisend.toml"

	if len(os.Args) == 2 {
		ConfFileName = os.Args[1]
	}
	fmt.Printf("TOML=%s\n", ConfFileName)
	if _, err := toml.DecodeFile(ConfFileName, &conf); err != nil {
		fmt.Println("Ошибка при загрузке файла конфигурации:", err.Error())
		return
	} else {
		fmt.Println("Данные файла конфигурации загружены!")
	}
	fNum := 0
	sdk = m.SDK{
		AccAddress:    conf.AccAddress,
		AccPrivateKey: conf.AccKey,
		ChainMainnet:  true,
	}
	for i, d := range conf.ListAddress {
		str0 := ""
		str1 := ""
		coinX := ""
		ok := true

		if str0, ok = d[0].(string); !ok {
			fmt.Println("Ошибка при загрузке файла конфигурации:", d[0], "неверный адрес получателя")
			return
		}
		if str1, ok = d[1].(string); !ok {
			fmt.Println("Ошибка при загрузке файла конфигурации:", d[1], "неверное количество")
			return
		}

		if len(d) == 3 {
			if coinX, ok = d[2].(string); !ok {
				fmt.Println("Ошибка при загрузке файла конфигурации:", d[2], "неверное название монеты")
				return
			}
			coinX = strings.ToUpper(coinX)
		}

		int1, err := strconv.Atoi(str1)
		if err != nil {
			fmt.Println("Ошибка при загрузке файла конфигурации:", str1, "неверное количество")
			return
		}

		if len(str0) == 42 && strings.HasPrefix(str0, "Mx") {
			sndDt := m.TxOneSendCoinData{
				Coin:      coinX,
				ToAddress: str0,
				Value:     float32(int1),
			}
			oSndData = append(oSndData, sndDt)
		} else {
			fmt.Println("Ошибка при загрузке файла конфигурации:", str0, "неверный адрес получателя")
			return
		}
		if (i+1)%100 == 0 {
			CreateFileTx(conf.Nonce + fNum)
			fNum += 1
			oSndData = oSndData[:0]
		}
	}
	if len(oSndData) < 100 {
		CreateFileTx(conf.Nonce + fNum)
	}

	fmt.Println("Пауза 30 сек")
	fmt.Println("Благодарность можно отправлять на адрес Mx9e6ce0f9ebb86c810261f844c5bf9fc089a67fd5")
	time.Sleep(time.Second * 30) // пауза 30 сек

}
