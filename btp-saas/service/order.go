package service

import (
	"fmt"
	"log"
	"time"

	"github.com/buyaobilian1/autobuytgvip-saas/btp-saas/dao/model"
	"github.com/buyaobilian1/autobuytgvip-saas/btp-saas/dao/query"
	"github.com/buyaobilian1/autobuytgvip-saas/btp-saas/global"
	"github.com/buyaobilian1/autobuytgvip-saas/btp-saas/mq"
	"github.com/buyaobilian1/autobuytgvip-saas/btp-saas/mq/handle"
	"github.com/hibiken/asynq"
)

type CreateOrderResponse struct {
	Token        string  `json:"token"`
	ActualAmount float64 `json:"actual_amount"`
}

func CreateOrder(order *model.Order) (result CreateOrderResponse, err error) {
	var o = query.Order
	notifyUrl := fmt.Sprintf(global.Conf.PayConf.NotifyUrl, "order")

	log.Printf("Creating order with OrderNo: %s, UsdtAmount: %f, NotifyUrl: %s", order.OrderNo, order.UsdtAmount, notifyUrl)

	payment, err := CreateEpusdtPayment(order.OrderNo, order.UsdtAmount, notifyUrl)
	if err != nil {
		log.Printf("Failed to create Epusdt payment: %v", err)
		return // 确保在有错误时返回
	}

	log.Printf("Epusdt payment response: %+v", payment) // 打印支付响应

	order.UsdtAmount = payment.Data.ActualAmount
	log.Printf("Updating order UsdtAmount to: %f", order.UsdtAmount)

	err = o.Create(order)
	if err != nil {
		log.Printf("Failed to create order in database: %v", err)
		return // 确保在有错误时返回
	}

	log.Printf("Order created successfully in database with ID: %d", order.ID) //假设model.Order有ID字段

	task, _ := handle.NewOrderExpirationTask(order.OrderNo)
	_, err = mq.QueueClient.Enqueue(task, asynq.ProcessIn(time.Minute*time.Duration(global.Conf.AppConf.OrderExpireMinute)))

	if err != nil {
		log.Printf("Failed to enqueue order expiration task: %v", err)
		//这里可以选择是否返回错误，根据业务需求决定
	} else {
		log.Printf("Order expiration task enqueued successfully for OrderNo: %s", order.OrderNo)
	}

	result = CreateOrderResponse{
		Token:        payment.Data.Token,
		ActualAmount: payment.Data.ActualAmount,
	}
	log.Printf("Returning CreateOrderResponse: %+v", result)
	return result, nil // 明确返回 nil 错误
}
