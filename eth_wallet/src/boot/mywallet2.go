package main

import (
	"flag"
	"log"
	"strconv"

	"fmt"
	"net/http"

	//"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/console"
	"mywallet"
)

var (
	gethAddr string
	keyDir   string
	keyFile  string
)

const ethToWei = 1 << 17

func init() {
	flag.StringVar(&gethAddr, "g", "http://localhost:8545", "geth server address")
	flag.StringVar(&keyDir, "d", "", "key dir to generate key")
	flag.StringVar(&keyFile, "f", "UTC--2019-01-13T12-36-52.775318000Z--c47fe9c4cf84f583f83a92d7f87e261d9ada7357", "key file path")
}

//成员变量
var client *mywallet.EthClient
var account *mywallet.Account

//客户端用于转账等，查询余额不需要
func getClient() *mywallet.EthClient {
	var err error
	if client == nil {
		client, err = mywallet.NewEthClient(gethAddr)
		if err != nil {
			log.Fatal("conn to geth failed:%s", err.Error())
		}
	}
	return client
}

//查黁余额的方法
func getBalance(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() //解析参数，默认是不会解析的
	//遍历请求中的参数，暂时用不到
	// for k, v := range r.Form {
	//     fmt.Println("key:", k)
	//     fmt.Println("val:", strings.Join(v, ""))
	//  }
	//var account *mywallet.Account

	//调用geth获取账户余额
	balance, err := getClient().GetBalance(account.Address())

	if err != nil {
		fmt.Printf("error: get balance failed:%s\n", err.Error())
	} else {
		fmt.Printf("%.8f eth\n", mywallet.WeiToEth(balance))
	}

	str := strconv.FormatFloat(mywallet.WeiToEth(balance), 'f', 8, 64) //小数点后保留8位，将float64转出string
	fmt.Fprintf(w, str)                                                //这个写入到w的是输出到客户端的

}

func main() {
	flag.Parse()
	log.SetFlags(0)
	fmt.Println("Hello,world")
	//启动时初始化账户
	var err error
	if keyDir != "" {
		account, err = mywallet.GenerateKey(keyDir)

		if err != nil {
			log.Fatal("generate key failed:%s", err.Error())
		}
	} else if keyFile != "" {
		//如果账户还没有根据keyFile生成
		if account == nil {
			//则输入密码生成account
			password, err := console.Stdin.PromptPassword("password:")
			if err != nil {
				log.Fatal("import key failed:%s", err.Error())
			}
			account, err = mywallet.NewAccount(keyFile, password)
			if err != nil {
				log.Fatal("key file validation failed:%s", err.Error())
			}
		}
	}

	//开启服务端
	http.HandleFunc("/getBalance", getBalance) //设置访问的路由
	err = http.ListenAndServe(":9090", nil)    //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

	// for {
	// 	cmd, err := console.Stdin.PromptInput("mywallect> ")
	// 	if err != nil {
	// 		log.Fatal("get cmd failed:%s", err.Error())
	// 	} else if cmd == "get_balance" {
	// 		balance, err := getClient().GetBalance(account.Address())
	// 		if err != nil {
	// 			log.Printf("error: get balance failed:%s\n", err.Error())
	// 		} else {
	// 			log.Printf("%.8f eth\n", mywallet.WeiToEth(balance))
	// 		}
	// 	} else if cmd == "address" {
	// 		log.Printf("%s\n", account.Address().Hex())
	// 	} else if cmd == "transfer" {
	// 		targetAddress, err := console.Stdin.PromptInput("target_account:")
	// 		if err != nil {
	// 			log.Fatal("get target account failed:%s", err.Error())
	// 		}

	// 		if common.IsHexAddress(targetAddress) == false {
	// 			log.Printf("error: %s isn't valid address\n", targetAddress)
	// 			continue
	// 		}

	// 		valueStr, err := console.Stdin.PromptInput("value(ether):")
	// 		if err != nil {
	// 			log.Fatal("get value failed:%s", err.Error())
	// 		}

	// 		value, err := strconv.ParseFloat(valueStr, 64)
	// 		if err != nil {
	// 			log.Println("error: value should be number")
	// 		continue
	// 	}

	// 	err = getClient().Transfer(account, common.HexToAddress(targetAddress), mywallet.EthToWei(value))
	// 	if err != nil {
	// 		log.Printf("error: transfer failed:%s\n", err.Error())
	// 	} else {
	// 		log.Println("transfer succeed")
	// 	}
	// } else if cmd == "exit" {
	// 	break
	// } else {
	// 	log.Println("error: cmd is unknown")
	// }
	//}
}
