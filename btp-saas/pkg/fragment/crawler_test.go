package fragment

import (
	"btp-saas/global"
	"btp-saas/internal/config"
	"encoding/base64"
	"flag"
	"fmt"
	"strings"
	"testing"

	"github.com/zeromicro/go-zero/core/conf"
)

func TestSearchPremiumGiftRecipient(t *testing.T) {
	var configFile = flag.String("f", "C:\\Users\\59740\\Desktop\\autobuytgvip-saas\\btp-saas/etc/btp-saas.yaml", "the config file")
	var c config.Config
	conf.MustLoad(*configFile, &c)
	global.Conf = c
	duration := 3

	result1, _ := SearchPremiumGiftRecipient("@minggetg", duration)
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
	orderSN := arr[1]

	commentFormatter := `Telegram Premium for 3 months 

Ref#%s`
	fmt.Printf("address: %s\namount: %d\ncomment:%s\n", receiverAddress, amount, fmt.Sprintf(commentFormatter, orderSN))
}
