package service

import (
	"fmt"
	"log"
	"time"

	"btp-saas/dao/model"
	"btp-saas/dao/query"
	"btp-saas/global"
	"btp-saas/mq"
	"btp-saas/mq/handle"

	"github.com/hibiken/asynq"
)

type CreateOrderResponse struct {
	Token        string  `json:"token"`
	ActualAmount float64 `json:"actual_amount"`
}

// CreateOrder 处理订单创建逻辑，根据支付方式调用不同的函数
func CreateOrder(order *model.Order, useEpusdt bool) (result CreateOrderResponse, err error) {
	if useEpusdt {
		result, err = CreateOrderWithEpusdt(order)
	} else {
		result, err = CreateOrderWithoutEpusdt(order)
	}
	return result, err
}

// CreateOrderWithEpusdt 创建使用 Epusdt 支付的订单
func CreateOrderWithEpusdt(order *model.Order) (result CreateOrderResponse, err error) {
	var o = query.Order
	notifyUrl := fmt.Sprintf(global.Conf.PayConf.NotifyUrl, "order")

	log.Printf("Creating order with Epusdt payment. OrderNo: %s, UsdtAmount: %f, NotifyUrl: %s", order.OrderNo, order.UsdtAmount, notifyUrl)

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

// CreateOrderWithoutEpusdt 创建不使用 Epusdt 支付的订单 (余额支付)
func CreateOrderWithoutEpusdt(order *model.Order) (result CreateOrderResponse, err error) {
	var o = query.Order

	log.Printf("Creating order without Epusdt payment (balance payment). OrderNo: %s, UsdtAmount: %f", order.OrderNo, order.UsdtAmount)

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

	// 根据你的业务逻辑，设置合适的返回值
	result = CreateOrderResponse{
		Token:        "",               // 或者设置其他合适的值
		ActualAmount: order.UsdtAmount, // 或者设置其他合适的值
	}
	log.Printf("Returning CreateOrderResponse: %+v", result)
	return result, nil // 明确返回 nil 错误
}
