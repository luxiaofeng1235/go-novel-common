package models

type McOrder struct {
	Id             int64   `gorm:"column:id" json:"id"`                             // 订单id
	OrderNo        string  `gorm:"column:order_no" json:"order_no"`                 // 流水号
	TradeNo        string  `gorm:"column:trade_no" json:"trade_no"`                 // 支付平台订单号
	Subject        string  `gorm:"column:subject" json:"subject"`                   // 商品标题
	PayType        int     `gorm:"column:pay_type" json:"pay_type"`                 // 账号类型 1-支付宝 2-微信
	UserId         int64   `gorm:"column:user_id" json:"user_id"`                   // 用户ID
	CardId         int64   `gorm:"column:card_id" json:"card_id"`                   // 会员卡ID
	CardName       string  `gorm:"column:card_name" json:"card_name"`               // 会员卡名称
	CardDay        int64   `gorm:"column:card_day" json:"card_day"`                 // 会员卡天数
	TotalAmount    float64 `gorm:"column:total_amount" json:"total_amount"`         // 订单金额
	TradeAmount    float64 `gorm:"column:trade_amount" json:"trade_amount"`         // 支付成功后 平台返回订单金额
	PaySuccessTime int64   `gorm:"column:pay_success_time" json:"pay_success_time"` // 订单支付成功时间
	ClientIp       string  `gorm:"column:client_ip" json:"client_ip"`               // 客户端IP
	Status         int     `gorm:"column:status" json:"status"`                     // 订单状态:0-待付款 1-支付成功 2-订单关闭
	Addtime        int64   `gorm:"column:addtime" json:"addtime"`                   // 创建时间
	Uptime         int64   `gorm:"column:uptime" json:"uptime"`                     // 更新时间
}

func (*McOrder) TableName() string {
	return "mc_order"
}

type CreateOrderReq struct {
	Vid      int64  `form:"vid" json:"vid"`
	PayType  int    `form:"pay_type" json:"pay_type"` //支付类型 1-支付宝 2-微信
	UserId   int64  `form:"user_id" json:"user_id"`
	ClientIp string `form:"client_ip"  json:"client_ip"`
}

type UnifiedOrderReq struct {
	MchId      string `form:"mchId" json:"mchId"`            //商户号
	WayCode    int    `form:"wayCode" json:"wayCode"`        //901	通道类型，详见 通道编码
	Subject    string `form:"subject" json:"subject"`        //商品标题
	Body       string `form:"body"  json:"body"`             //商品描述
	OutTradeNo string `form:"outTradeNo"  json:"outTradeNo"` //商户生成的订单号
	Amount     int64  `form:"amount"  json:"amount"`         //支付金额 (单位: 分)，例如: 10000 即为 100.00 元
	ExtParam   string `form:"extParam"  json:"extParam"`     //扩展参数 商户扩展参数,回调时会原样返回
	ClientIp   string `form:"clientIp"  json:"clientIp"`     //客户端 IPV4 地址，尽量填写
	NotifyUrl  string `form:"notifyUrl"  json:"notifyUrl"`   //https://www.test.com/notify.htm	支付结果异步回调URL，只有传了该值才会发起回调
	ReturnUrl  string `form:"returnUrl"  json:"returnUrl"`   //https://www.test.com/return.htm	支付结果同步跳转通知URL
	ReqTime    string `form:"reqTime"  json:"reqTime"`       //请求时间 请求接口时间，13位时间戳
	Sign       string `form:"sign"  json:"sign"`             //签名值，不参与签名，详见 签名算法
}

type UnifiedOrderRes struct {
	Code    int                  `form:"code" json:"code"` //网关返回码：0=成功，其他失败
	Message string               `form:"message" json:"message"`
	Sign    string               `form:"sign" json:"sign"`
	Data    *UnifiedOrderDataRes `form:"data"  json:"data"`
}

type UnifiedOrderDataRes struct {
	MchId         string `form:"mchId" json:"mchId"`                  //商户号
	TradeNo       string `form:"tradeNo" json:"tradeNo"`              //订单号
	OutTradeNo    string `form:"outTradeNo" json:"outTradeNo"`        //流水号
	OriginTradeNo string `form:"originTradeNo"  json:"originTradeNo"` //通道订单号
	Amount        string `form:"amount"  json:"amount"`               //订单金额 (单位: 分) 例如: 10000 即为 100.00 元
	PayUrl        string `form:"payUrl"  json:"payUrl"`               //支付地址 商户需要展示该地址给付款用户。
	ExpiredTime   string `form:"expiredTime"  json:"expiredTime"`     //订单有效截止时间，13位时间戳
	SdkData       string `form:"sdkData"  json:"sdkData"`             //SDK参数内容，例如支付宝现金红包需要
}

type OrderNotifyTestReq struct {
	OrderNo     string  `form:"order_no" json:"order_no"`
	TradeNo     string  `form:"trade_no"  json:"trade_no"`
	PayType     int     `form:"pay_type" json:"pay_type"`
	Status      int     `form:"status" json:"status"`
	TradeAmount float64 `form:"trade_amount" json:"trade_amount"`
}

type OrderNotifyReq struct {
	MchId         string `form:"mchId" json:"mchId"`
	TradeNo       string `form:"trade_no"  json:"trade_no"`
	OutTradeNo    string `form:"outTradeNo" json:"outTradeNo"`       //流水号
	OriginTradeNo string `form:"originTradeNo" json:"originTradeNo"` //通道订单号
	Amount        int64  `form:"trade_amount" json:"trade_amount"`   //订单金额 (单位: 分)，例如: 10000 即为 100.00 元
	Subject       string `form:"subject" json:"subject"`
	Body          string `form:"body" json:"body"`
	ExtParam      string `form:"extParam" json:"extParam"`
	State         int    `form:"state" json:"state"`           //订单状态：0=待支付，1=支付成功，2=支付失败，3=未出码，4=异常
	NotifyTime    int64  `form:"notifyTime" json:"notifyTime"` //通知时间，13位时间戳
	Sign          string `form:"sign" json:"sign"`
}

type QueryOrderReq struct {
	OrderNo string `form:"order_no"  json:"order_no"`
}

type QueryOrderRes struct {
	Code    int                `form:"code" json:"code"` //网关返回码：0=成功，其他失败
	Message string             `form:"message" json:"message"`
	Sign    string             `form:"sign" json:"sign"`
	Data    *QueryOrderDataRes `form:"data"  json:"data"`
}

type QueryOrderDataRes struct {
	MchId         string `form:"mchId" json:"mchId"`                  //商户号
	WayCode       int    `form:"wayCode" json:"wayCode"`              //通道类型，详见 通道编码
	TradeNo       string `form:"tradeNo" json:"tradeNo"`              //订单号
	OutTradeNo    string `form:"outTradeNo"  json:"outTradeNo"`       //流水号
	OriginTradeNo string `form:"originTradeNo"  json:"originTradeNo"` //支付通道订单号
	Amount        string `form:"amount"  json:"amount"`               //订单金额 (单位: 分)，例如: 10000 即为 100.00 元
	Subject       string `form:"subject"  json:"subject"`             //商品标题测试
	Body          string `form:"body"  json:"body"`                   //商品描述
	ExtParam      string `form:"extParam"  json:"extParam"`           //商户扩展参数，回调时会原样返回
	NotifyUrl     string `form:"notifyUrl"  json:"notifyUrl"`         //支付结果异步回调URL
	PayUrl        string `form:"payUrl"  json:"payUrl"`               //支付地址，商户需要展示该地址给付款用户
	ExpiredTime   string `form:"expiredTime"  json:"expiredTime"`     //订单过期时间
	SuccessTime   string `form:"successTime"  json:"successTime"`     //支付成功时间
	CreateTime    string `form:"createTime"  json:"createTime"`       //下单时间
	State         int    `form:"state"  json:"state"`                 //订单状态：0=待支付，1=支付成功，2=支付失败，3=未出码，4=异常
	NotifyState   int    `form:"notifyState"  json:"notifyState"`     //通知状态：0=未通知，1=通知成功，2=通知失败
}

type UpdateOrderSuccessReq struct {
	OrderId        int64   `json:"id"`
	Status         int     `form:"status"     json:"status"`
	TradeNo        string  `form:"trade_no"     json:"trade_no"`
	TradeAmount    float64 `form:"trade_amount"     json:"trade_amount"`
	PaySuccessTime int64   `form:"pay_success_time"     json:"pay_success_time"`
}

type OrderListReq struct {
	OrderNo   string `form:"order_no" json:"order_no"`
	UserId    string `form:"user_id" json:"user_id"`
	PayType   string `form:"pay_type" json:"pay_type"`
	Status    string `form:"status" json:"status"`
	BeginTime string `form:"beginTime" json:"beginTime"`
	EndTime   string `form:"endTime" json:"endTime"`
	PageNum   int    `form:"pageNum" json:"pageNum"`
	PageSize  int    `form:"pageSize" json:"pageSize"`
}

type UpdateOrderReq struct {
	OrderId     int64   `form:"id"  json:"id"`
	Status      int     `form:"status" json:"status"`
	TradeAmount float64 `form:"trade_amount" json:"trade_amount"`
}
