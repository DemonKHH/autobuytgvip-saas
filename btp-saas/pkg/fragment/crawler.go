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
	Amount  string `json:"amount"`
	Payload string `json:"payload"`
}

// ConfirmParams represents the confirmation parameters.
type ConfirmParams struct {
	ID string `json:"id"`
}

func GetGiftPremiumLink(reqId string) (result GetGiftPremiumLinkResponse, err error) {
	// 正确的 JSON 字符串 (注意引号的转义)
	accountJSON := `{"address":"0:2dd4d2055fa91e2624d6bc93a160c182dd475aaa150fa631b4a62407f5283c45","chain":"-239","walletStateInit":"te6cckECFgEAAwQAAgE0AgEAUQAAAAApqaMXt6WWfqVvrCnbM5weIa4rzy0JyQvpSdhS04sAEPVotThAART/APSkE/S88sgLAwIBIAkEBPjygwjXGCDTH9Mf0x8C+CO78mTtRNDTH9Mf0//0BNFRQ7ryoVFRuvKiBfkBVBBk+RDyo/gAJKTIyx9SQMsfUjDL/1IQ9ADJ7VT4DwHTByHAAJ9sUZMg10qW0wfUAvsA6DDgIcAB4wAhwALjAAHAA5Ew4w0DpMjLHxLLH8v/CAcGBQAK9ADJ7VQAbIEBCNcY+gDTPzBSJIEBCPRZ8qeCEGRzdHJwdIAYyMsFywJQBc8WUAP6AhPLassfEss/yXP7AABwgQEI1xj6ANM/yFQgR4EBCPRR8qeCEG5vdGVwdIAYyMsFywJQBs8WUAT6AhTLahLLH8s/yXP7AAIAbtIH+gDU1CL5AAXIygcVy//J0Hd0gBjIywXLAiLPFlAF+gIUy2sSzMzJc/sAyEAUgQEI9FHypwICAUgTCgIBIAwLAFm9JCtvaiaECAoGuQ+gIYRw1AgIR6STfSmRDOaQPp/5g3gSgBt4EBSJhxWfMYQCASAODQARuMl+1E0NcLH4AgFYEg8CASAREAAZrx32omhAEGuQ64WPwAAZrc52omhAIGuQ64X/wAA9sp37UTQgQFA1yH0BDACyMoHy//J0AGBAQj0Cm+hMYALm0AHQ0wMhcbCSXwTgItdJwSCSXwTgAtMfIYIQcGx1Z70ighBkc3RyvbCSXwXgA/pAMCD6RAHIygfL/8nQ7UTQgQFA1yH0BDBcgQEI9ApvoTGzkl8H4AXTP8glghBwbHVnupI4MOMNA4IQZHN0crqSXwbjDRUUAIpQBIEBCPRZMO1E0IEBQNcgyAHPFvQAye1UAXKwjiOCEGRzdHKDHrFwgBhQBcsFUAPPFiP6AhPLassfyz/JgED7AJJfA+IAeAH6APQEMPgnbyIwUAqhIb7y4FCCEHBsdWeDHrFwgBhQBMsFJs8WWPoCGfQAy2kXyx9SYMs/IMmAQPsABgwneNg=","publicKey":"b7a5967ea56fac29db339c1e21ae2bcf2d09c90be949d852d38b0010f568b538"}`
	deviceJSON := `{"platform":"windows","appName":"tonkeeper","appVersion":"4.2.2","maxProtocolVersion":2,"features":["SendTransaction",{"name":"SendTransaction","maxMessages":4,"extraCurrencySupported":true},{"name":"SignData","types":["text","binary","cell"]}]}`
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
				Amount  string `json:"amount"`
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
