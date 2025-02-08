package fragment

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"btp-saas/global"

	"github.com/zeromicro/go-zero/rest/httpc"
)

type SearchPremiumGiftRecipientRequest struct {
	Cookie string `header:"cookie"`
	Query  string `form:"query"`
	Months int    `form:"months"`
	Method string `form:"method"`
}

type SearchPremiumGiftRecipientResponse struct {
	Ok    bool                                   `json:"ok"`
	Error string                                 `json:"error"`
	Found SearchPremiumGiftRecipientResponseBody `json:"found"`
}

type SearchPremiumGiftRecipientResponseBody struct {
	Myself    bool   `json:"myself"`
	Recipient string `json:"recipient"`
	Photo     string `json:"photo"`
	Name      string `json:"name"`
}

// SearchPremiumGiftRecipient telegram用户名查询
func SearchPremiumGiftRecipient(username string, duration int) (result SearchPremiumGiftRecipientResponse, err error) {
	fragmentUrl := fmt.Sprintf("https://fragment.com/api?hash=%s", global.Conf.AppConf.Hash)
	req := SearchPremiumGiftRecipientRequest{
		Cookie: global.Conf.AppConf.Cookie,
		Query:  username,
		Months: duration,
		Method: "searchPremiumGiftRecipient",
	}
	resp, err := httpc.Do(context.Background(), http.MethodPost, fragmentUrl, req)
	if err != nil {
		return
	}
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	result = SearchPremiumGiftRecipientResponse{}
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return
	}
	return result, nil

}

type InitGiftPremiumRequest struct {
	Cookie    string `header:"cookie"`
	Recipient string `form:"recipient"`
	Months    int    `form:"months"`
	Method    string `form:"method"`
}

type InitGiftPremiumResponse struct {
	ReqId     string `json:"req_id"`
	Myself    bool   `json:"myself"`
	Amount    string `json:"amount"`
	ItemTitle string `json:"item_title"`
	Content   string `json:"content"`
	Button    string `json:"button"`
}

func InitGiftPremium(recipient string, duration int) (result InitGiftPremiumResponse, err error) {
	fragmentUrl := fmt.Sprintf("https://fragment.com/api?hash=%s", global.Conf.AppConf.Hash)
	req := InitGiftPremiumRequest{
		Cookie:    global.Conf.AppConf.Cookie,
		Recipient: recipient,
		Months:    duration,
		Method:    "initGiftPremiumRequest",
	}
	resp, err := httpc.Do(context.Background(), http.MethodPost, fragmentUrl, req)
	if err != nil {
		return
	}
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	result = InitGiftPremiumResponse{}
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return
	}
	return result, nil
}

type GetGiftPremiumLinkRequest struct {
	Cookie      string `header:"cookie"`
	Id          string `form:"id"`
	ShowSender  int    `form:"show_sender"`
	Months      int    `form:"months"`
	Method      string `form:"method"`
	Transaction int    `form:"transaction"`
	Account     string `form:"account"`
	Device      string `form:"device"`
}

type GetGiftPremiumLinkResponse struct {
	Ok            bool          `json:"ok"`
	Transaction   Transaction   `json:"transaction"`
	ConfirmMethod string        `json:"confirm_method"`
	ConfirmParams ConfirmParams `json:"confirm_params"`
}

// Transaction represents the top-level transaction object.
type Transaction struct {
	ValidUntil int64     `json:"validUntil"`
	From       string    `json:"from"`
	Messages   []Message `json:"messages"`
}

// Message represents a single message within the transaction.
type Message struct {
	Address string `json:"address"`
	Amount  uint64 `json:"amount"`
	Payload string `json:"payload"`
}

// ConfirmParams represents the confirmation parameters.
type ConfirmParams struct {
	ID string `json:"id"`
}

func GetGiftPremiumLink(reqId string) (result GetGiftPremiumLinkResponse, err error) {
	// 正确的 JSON 字符串 (注意引号的转义)
	accountJSON := `{"address":"0:3a7d7318b5d38d910f0cf68fc4d60b95c889b87e365ea00db05ab11fc9e5f523","chain":"-239","walletStateInit":"te6cckECFgEAAwQAAgE0ARUBFP8A9KQT9LzyyAsCAgEgAxACAUgEBwLm0AHQ0wMhcbCSXwTgItdJwSCSXwTgAtMfIYIQcGx1Z70ighBkc3RyvbCSXwXgA/pAMCD6RAHIygfL/8nQ7UTQgQFA1yH0BDBcgQEI9ApvoTGzkl8H4AXTP8glghBwbHVnupI4MOMNA4IQZHN0crqSXwbjDQUGAHgB+gD0BDD4J28iMFAKoSG+8uBQghBwbHVngx6xcIAYUATLBSbPFlj6Ahn0AMtpF8sfUmDLPyDJgED7AAYAilAEgQEI9Fkw7UTQgQFA1yDIAc8W9ADJ7VQBcrCOI4IQZHN0coMesXCAGFAFywVQA88WI/oCE8tqyx/LP8mAQPsAkl8D4gIBIAgPAgEgCQ4CAVgKCwA9sp37UTQgQFA1yH0BDACyMoHy//J0AGBAQj0Cm+hMYAIBIAwNABmtznaiaEAga5Drhf/AABmvHfaiaEAQa5DrhY/AABG4yX7UTQ1wsfgAWb0kK29qJoQICga5D6AhhHDUCAhHpJN9KZEM5pA+n/mDeBKAG3gQFImHFZ8xhAT48oMI1xgg0x/TH9MfAvgju/Jk7UTQ0x/TH9P/9ATRUUO68qFRUbryogX5AVQQZPkQ8qP4ACSkyMsfUkDLH1Iwy/9SEPQAye1U+A8B0wchwACfbFGTINdKltMH1AL7AOgw4CHAAeMAIcAC4wABwAORMOMNA6TIyx8Syx/L/xESExQAbtIH+gDU1CL5AAXIygcVy//J0Hd0gBjIywXLAiLPFlAF+gIUy2sSzMzJc/sAyEAUgQEI9FHypwIAcIEBCNcY+gDTP8hUIEeBAQj0UfKnghBub3RlcHSAGMjLBcsCUAbPFlAE+gIUy2oSyx/LP8lz+wACAGyBAQjXGPoA0z8wUiSBAQj0WfKnghBkc3RycHSAGMjLBcsCUAXPFlAD+gITy2rLHxLLP8lz+wAACvQAye1UAFEAAAAAKamjF4x3DHeYFf9NS3Wc5G0ZiWLw03yqwgISSNu5LERScRlnQA1I5nI=","publicKey":"8c770c779815ff4d4b759ce46d198962f0d37caac2021248dbb92c4452711967"}`
	deviceJSON := `{"platform":"web","appName":"tonwallet","appVersion":"1.1.49","maxProtocolVersion":2,"features":["SendTransaction",{"name":"SendTransaction","maxMessages":4}]}`
	fragmentUrl := fmt.Sprintf("https://fragment.com/api?hash=%s", global.Conf.AppConf.Hash)
	req := GetGiftPremiumLinkRequest{
		Cookie:      global.Conf.AppConf.Cookie,
		Id:          reqId,
		ShowSender:  0,
		Transaction: 1,
		Method:      "getGiftPremiumLink",
		Account:     accountJSON,
		Device:      deviceJSON,
	}
	resp, err := httpc.Do(context.Background(), http.MethodPost, fragmentUrl, req)
	if err != nil {
		return
	}
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	result = GetGiftPremiumLinkResponse{}
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return
	}
	return result, nil
}

type GetTonPaymentInfoRequest struct {
	Cookie string `header:"cookie"`
}

type GetTonPaymentInfoResponse struct {
	Version string `json:"version"`
	Body    struct {
		Type   string `json:"type"`
		Params struct {
			ValidUntil int `json:"valid_until"`
			Messages   []struct {
				Address string `json:"address"`
				Amount  uint64 `json:"amount"`
				Payload string `json:"payload"`
			} `json:"messages"`
			Source string `json:"source"`
		} `json:"params"`
		ResponseOptions struct {
			CallbackURL string `json:"callback_url"`
			Broadcast   bool   `json:"broadcast"`
		} `json:"response_options"`
		ExpiresSec int `json:"expires_sec"`
	} `json:"body"`
}

// GetTonPaymentInfo GET获取收款地址和 payload参数
func GetTonPaymentInfo(id string) (result GetTonPaymentInfoResponse, err error) {
	fragmentUrl := fmt.Sprintf("https://fragment.com/tonkeeper/rawRequest?id=%s&qr=1", id)
	req := GetTonPaymentInfoRequest{
		Cookie: global.Conf.AppConf.Cookie,
	}
	resp, err := httpc.Do(context.Background(), http.MethodGet, fragmentUrl, req)
	if err != nil {
		return
	}
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	result = GetTonPaymentInfoResponse{}
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return
	}
	return result, nil
}

type CheckOrderRequest struct {
	Cookie string `header:"cookie"`
	Id     string `form:"id"`
	Method string `form:"method"`
}

type CheckOrderResponse struct {
	Confirmed bool `json:"confirmed"`
}

// CheckOrder 检查订单是否成功
func CheckOrder(id string) (result CheckOrderResponse, err error) {
	fragmentUrl := fmt.Sprintf("https://fragment.com/api?hash=%s", global.Conf.AppConf.Hash)
	req := CheckOrderRequest{
		Cookie: global.Conf.AppConf.Cookie,
		Id:     id,
		Method: "checkReq",
	}
	resp, err := httpc.Do(context.Background(), http.MethodGet, fragmentUrl, req)
	if err != nil {
		return
	}
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	result = CheckOrderResponse{}
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return
	}
	return result, nil
}
