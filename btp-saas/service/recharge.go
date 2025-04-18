package service

import (
	"fmt"
	"time"

	"btp-saas/dao/model"
	"btp-saas/dao/query"
	"btp-saas/global"
	"btp-saas/mq"
	"btp-saas/mq/handle"

	"github.com/hibiken/asynq"
)

type CreateRechargeOrderResponse struct {
	Token        string  `json:"token"`
	ActualAmount float64 `json:"actual_amount"`
}

func CreateRechargeOrder(v *model.Recharge) (result CreateOrderResponse, err error) {
	var r = query.Recharge
	notifyUrl := fmt.Sprintf(global.Conf.PayConf.NotifyUrl, "recharge")
	payment, err := CreateEpusdtPayment(v.OrderNo, v.Amount, notifyUrl)
	if err != nil {
		return
	}

	v.ActualAmount = payment.Data.ActualAmount
	err = r.Create(v)
	if err != nil {
		return
	}

	result.Token = payment.Data.Token
	result.ActualAmount = payment.Data.ActualAmount

	// 加入延时队列
	task, _ := handle.NewRechargeExpirationTask(v.OrderNo)
	_, _ = mq.QueueClient.Enqueue(task, asynq.ProcessIn(time.Minute*time.Duration(global.Conf.AppConf.RechargeExpireMinute)))

	return
}
