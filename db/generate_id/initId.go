package generate_id

var ThirdBianOrder *Worker
var ThirHBOrder *Worker
var PlatformOrder *Worker

const (
	WORKER_ID           = 1 //机器中心
	PLATFORM_ORDER_ID   = 1 //平台订单ID
	THIRD_BIAN_ORDER_ID = 2 //币安客户端订单ID
	THIRD_HB_ORDER_ID   = 3 //火币客户端订单ID
)

func InitId() {
	ThirdBianOrder = NewWorker(WORKER_ID, THIRD_BIAN_ORDER_ID)
	ThirHBOrder = NewWorker(WORKER_ID, THIRD_HB_ORDER_ID)
	PlatformOrder = NewWorker(WORKER_ID, PLATFORM_ORDER_ID)
}
