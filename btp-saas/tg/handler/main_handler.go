package handler

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"btp-saas/dao/model"
	"btp-saas/dao/query"
	"btp-saas/mq"
	"btp-saas/mq/handle"
	"btp-saas/pkg/fragment"
	"btp-saas/pkg/id"
	"btp-saas/pkg/image"
	"btp-saas/service"

	tele "gopkg.in/telebot.v3"
)

func StartHandler(ctx tele.Context) error {
	_, err := service.FindOrCreateUserByTgCtx(ctx)
	if err != nil {
		return err
	}
	dbUser, err := service.FindUserByBotId(ctx.Bot().Me.ID)
	if err != nil {
		return err
	}
	var startFormatText = `
	❤️本机器人向您提供Telegram Premium会员自动开通服务！
	
	当前价格：
	*  3个月 / %s U*
	*  6个月 / %s U*
	*12个月 / %s U*🔥
	
	请选择下方菜单：`
	price3 := EscapeText(tele.ModeMarkdownV2, Float64Format(*dbUser.ThreeMonthPrice))
	price6 := EscapeText(tele.ModeMarkdownV2, Float64Format(*dbUser.SixMonthPrice))
	price12 := EscapeText(tele.ModeMarkdownV2, Float64Format(*dbUser.TwelveMonthPrice))
	startText := fmt.Sprintf(startFormatText, price3, price6, price12)
	keyboards := &tele.ReplyMarkup{ResizeKeyboard: true}
	keyboards.Reply(
		// keyboards.Row(RechargeKeyboard, MineKeyboard),
		keyboards.Row(AgentKeyboard, CooperationKeyboard),
	)
	_ = ctx.Send(startText, keyboards)

	startText2 := "尊贵的Telegram用户您好！\n\n请选择为谁开通/续费TG会员："
	replyMarkup := &tele.ReplyMarkup{}
	btnBuyMyself := replyMarkup.Data("✈️此账号开通", ByMyselfBtnId, "@"+ctx.Sender().Username)
	btnGift := replyMarkup.Data("🎁赠送给他人", GiftOtherBtnId)
	replyMarkup.Inline(
		replyMarkup.Row(btnBuyMyself, btnGift),
		replyMarkup.Row(CloseBtn, SupportBtn),
	)
	return ctx.Send(EscapeText(tele.ModeMarkdownV2, startText2), replyMarkup)
}

func BuyMyselfHandler(ctx tele.Context) error {
	arr := strings.Split(ctx.Data(), "|")
	return ShowTgUserInfo(ctx, arr[1])
}

func GiftOtherHandler(ctx tele.Context) error {
	var giftFormatText = `
请直接发送你需要开通会员的Telegram用户名：

*提示：用户名以@开头，如 %s*
`
	giftText := fmt.Sprintf(giftFormatText, "@"+ctx.Sender().Username)
	replyMarkup := &tele.ReplyMarkup{
		ForceReply:     true,
		Placeholder:    "请输入Tg用户名",
		ResizeKeyboard: true,
	}

	return ctx.Send(giftText, replyMarkup)
}

func BuyThreeMonthHandler(ctx tele.Context) error {
	return CreateTelegramPremiumOrder(ctx)
}

func BuySixMonthHandler(ctx tele.Context) error {
	return CreateTelegramPremiumOrder(ctx)
}

func BuyTwelveMonthHandler(ctx tele.Context) error {
	return CreateTelegramPremiumOrder(ctx)
}

func CreateTelegramPremiumOrder(ctx tele.Context) error {
	var u, p = query.User, query.Param
	dbAgentUser, err := u.Where(u.BotID.Eq(ctx.Bot().Me.ID)).First()
	if err != nil {
		log.Printf("[db] 查询数据失败, %v\n", err)
		return err
	}
	dbUser, err := u.Where(u.TgID.Eq(ctx.Sender().ID)).First()
	if err != nil {
		log.Printf("[db] 查询数据失败, %v\n", err)
		return err
	}
	basePriceObj, err := p.Where(p.K.Eq("base_price")).First()
	if err != nil {
		log.Printf("[db] 查询数据失败. %v\n", err)
		return err
	}

	params := strings.Split(ctx.Data(), "|")
	distUsername := params[1]
	vipMonth, _ := strconv.Atoi(params[3])
	var usdtAmount, baseAmount float64
	if vipMonth == 3 {
		usdtAmount = *dbAgentUser.ThreeMonthPrice
		baseAmount, err = strconv.ParseFloat(*basePriceObj.V1, 64)
		if err != nil {
			log.Printf("base price3 设置失败. %v\n", err)
			return err
		}
	} else if vipMonth == 6 {
		usdtAmount = *dbAgentUser.SixMonthPrice
		baseAmount, err = strconv.ParseFloat(*basePriceObj.V2, 64)
		if err != nil {
			log.Printf("base price3 设置失败. %v\n", err)
			return err
		}
	} else if vipMonth == 12 {
		usdtAmount = *dbAgentUser.TwelveMonthPrice
		baseAmount, err = strconv.ParseFloat(*basePriceObj.V3, 64)
		if err != nil {
			log.Printf("base price3 设置失败. %v\n", err)
			return err
		}
	} else {
		return errors.New("套餐错误")
	}

	order := &model.Order{
		OrderNo:           id.GenerateId(1),
		UserID:            dbUser.ID,
		AgentUserID:       dbAgentUser.ID,
		BotID:             ctx.Bot().Me.ID,
		ReceiveTgUsername: distUsername,
		VipMonth:          int32(vipMonth),
		UsdtAmount:        usdtAmount,
		BaseAmount:        baseAmount,
		Status:            1,
		CreatedAt:         time.Now(),
		ExpiredAt:         time.Now().Add(10 * time.Minute),
		TgChatID:          ctx.Chat().ID,
		TgMsgID:           0, //先置零，登消息发出去后得到消息id后再更新
	}
	log.Printf("order: %+v\n", order)

	err = HandleBalancePay(ctx, order)
	if err != nil {
		log.Printf("余额支付失败，尝试USDT支付: %v\n", err)
		// 只有当余额支付失败时才进行USDT支付
		err = HandleUsdtPay(ctx, order)
		if err != nil {
			return err // 返回USDT支付的错误，如果也失败了
		}
	}

	return nil
}

// 使用USDT支付
func HandleUsdtPay(ctx tele.Context, order *model.Order) error {
	orderFormatText := `❗️❗️❗️请注意：网络必须是TRC\-20，否则无法到账
❗️❗️❗️请注意，金额必须与下面的一致（一位都不能少）
👇*请向以下地址转账 %s USDT*

%s

👆点击复制上面地址进行支付，或者扫描上面二维码支付。
`
	var o = query.Order
	res, err := service.CreateOrder(order, true)
	if err != nil {
		log.Printf("创建订单失败: %v\n", err)
		return ctx.Respond(&tele.CallbackResponse{
			Text:      "系统繁忙，订单创建失败，请重试",
			ShowAlert: true,
		})
	}

	replyMarkup := &tele.ReplyMarkup{}
	replyMarkup.Inline(
		replyMarkup.Row(SupportBtn),
	)

	amountStr := Float64Format(res.ActualAmount)
	replyText := fmt.Sprintf(orderFormatText, EscapeText(tele.ModeMarkdownV2, amountStr), res.Token)
	context := &tele.Photo{
		File:    tele.FromReader(image.GenQrcode(res.Token)),
		Caption: replyText,
	}
	msg, _ := ctx.Bot().Send(ctx.Recipient(), context, replyMarkup)
	_, err = o.Where(o.OrderNo.Eq(order.OrderNo)).Update(o.TgMsgID, msg.ID)
	return err
}

// 使用余额支付
func HandleBalancePay(ctx tele.Context, order *model.Order) error {
	var o = query.Order
	var u = query.User
	res, err := service.CreateOrder(order, false)
	if err != nil {
		log.Printf("创建订单失败: %v\n", err)
		return ctx.Respond(&tele.CallbackResponse{
			Text:      "系统繁忙，订单创建失败，请重试",
			ShowAlert: true,
		})
	}
	dbUser, err := service.FindOrCreateUserByTgCtx(ctx)
	if err != nil {
		log.Printf("[db] 查询失败. : %v, dbuser: %v", err, dbUser)
		return err // 返回 FindOrCreateUserByTgCtx 的错误
	}
	log.Printf("[order] 进行余额支付，当前余额：%+v\n", dbUser.Balance)
	if dbUser.Balance < res.ActualAmount {
		log.Printf("[order] 余额不足，切换为USDT支付\n")
		ctx.Bot().Send(ctx.Recipient(), EscapeText(tele.ModeMarkdownV2, "余额不足，切换为USDT支付"))
		return errors.New("余额不足") // 返回一个错误，触发USDT支付
	}
	_, err = u.Where(u.ID.Eq(dbUser.ID)).Update(u.Balance, u.Balance.Sub(res.ActualAmount))
	if err != nil {
		log.Printf("[db] 更新余额失败. %v\n", err)
		return err // 返回数据库错误
	}
	_, err = o.Where(o.OrderNo.Eq(order.OrderNo), o.Status.Eq(1)).Update(o.Status, 2)
	if err != nil {
		log.Printf("[db] 更新订单失败. %v\n", err)
		return err // 返回数据库错误
	}

	msg, _ := ctx.Bot().Send(ctx.Recipient(), EscapeText(tele.ModeMarkdownV2, "🎉🎉🎉支付成功，正在为您开通会员..."))
	_, err = o.Where(o.OrderNo.Eq(order.OrderNo)).Update(o.TgMsgID, msg.ID)
	task, _ := handle.NewGiftTelegramPremiumTask(order.OrderNo)
	_, _ = mq.QueueClient.Enqueue(task)
	return err
}

func ShowTgUserInfo(ctx tele.Context, username string) error {
	var u = query.User
	dbUser, err := u.Where(u.BotID.Eq(ctx.Bot().Me.ID)).First()
	if err != nil {
		log.Printf("query db fail: %v\n", err)
		return err
	}
	var replyText string
	var replyFormatText = `
开通用户：%s
用户昵称：%s

确定为此用户 开通/续费 Telegram Premium会员吗？
`
	userInfo, err := fragment.SearchPremiumGiftRecipient(username, 3)
	if err != nil {
		log.Printf("fail to get premium gift recipient, %v", err)
		return nil
	}
	if userInfo.Error == "No Telegram users found." {
		return ctx.Send(EscapeText(tele.ModeMarkdownV2, "用户名不存在."))
	}
	if userInfo.Error == "This account is already subscribed to Telegram Premium." {
		return ctx.Send(EscapeText(tele.ModeMarkdownV2, "此账号已经订阅会员."))
	}

	replyText = EscapeText(tele.ModeMarkdownV2, fmt.Sprintf(replyFormatText, username, userInfo.Found.Name))

	replyMarkup := &tele.ReplyMarkup{}
	btnBuy3Month := replyMarkup.Data(fmt.Sprintf("3个月 / %s U", Float64Format(*dbUser.ThreeMonthPrice)), BuyThreeMonthBtnId, username, fmt.Sprintf("%d", ctx.Sender().ID), "3")
	btnBuy6Month := replyMarkup.Data(fmt.Sprintf("6个月 / %s U", Float64Format(*dbUser.SixMonthPrice)), BuySixMonthBtnId, username, fmt.Sprintf("%d", ctx.Sender().ID), "6")
	btnBuy12Month := replyMarkup.Data(fmt.Sprintf("12个月 / %s U🔥", Float64Format(*dbUser.TwelveMonthPrice)), BuyTwelveMonthBtnId, username, fmt.Sprintf("%d", ctx.Sender().ID), "12")

	replyMarkup.Inline(
		replyMarkup.Row(btnBuy3Month, btnBuy6Month),
		replyMarkup.Row(btnBuy12Month),
		replyMarkup.Row(CloseBtn, SupportBtn),
	)

	return ctx.Send(replyText, replyMarkup)
}
