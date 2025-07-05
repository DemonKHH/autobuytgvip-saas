package fragment

import (
	"btp-saas/global"
	"btp-saas/internal/config"
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"regexp"
	"strings"
	"testing"

	"github.com/zeromicro/go-zero/core/conf"
)

func removeInvalidChars4(s string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9\-\_\.\[\]\(\)\{\}]`) //  Keep these characters
	return re.ReplaceAllString(s, "")
}

func TestSearchPremiumGiftRecipient(t *testing.T) {
	var configFile = flag.String("f", "C:\\Users\\59740\\Desktop\\autobuytgvip-saas\\btp-saas/etc/btp-saas.yaml", "the config file")
	var c config.Config
	conf.MustLoad(*configFile, &c)
	global.Conf = c
	log.Printf("global.Conf: %v", global.Conf)
	duration := 3

	result1, _ := SearchPremiumGiftRecipient("@demonkinghaha", duration)
	fmt.Printf("查询Telegram用户信息：%+v\n", result1)

	result2, _ := InitGiftPremium(result1.Found.Recipient, duration)
	fmt.Printf("初始化赠送Telegram会员请求：%+v\n", result2)
	fmt.Printf("请求ID：%s\n", result2.ReqId)
	result3, _ := GetGiftPremiumLink(result2.ReqId)
	fmt.Printf("获取Telegr会员：%+v\n", result3)

	if !result3.Ok {
		fmt.Print("获取Telegr会员失败\n")
		return
	}
	info := result3.Transaction

	receiverAddress := info.Messages[0].Address
	amount := info.Messages[0].Amount
	payload := info.Messages[0].Payload

	decodeBytes, _ := base64.RawStdEncoding.DecodeString(payload)

	arr := strings.Split(string(decodeBytes), "#")
	orderSN := removeInvalidChars4(arr[1])
	fmt.Printf("订单号：%s\n", orderSN)
	commentFormatter := `Telegram Premium for 3 months 

Ref#%s`
	fmt.Printf("address: %s\namount: %s\ncomment:%s\n", receiverAddress, amount, fmt.Sprintf(commentFormatter, orderSN))
}
