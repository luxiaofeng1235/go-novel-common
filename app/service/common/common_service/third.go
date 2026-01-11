package common_service

//这里面是写关于第三方的调用的一些共用函数和业务类
import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-novel/app/models"
	"go-novel/app/service/api/user_service"
	"go-novel/app/service/api/version_service"
	"go-novel/global"
	"go-novel/pkg/config"
	"go-novel/utils"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

/*
* @note 上报小米的事件推送
* @param id int64 包ID
* @param c object gin对象信息
* @param device_type string 客户端类型
* @param package_name string 包名
* @param mark string  渠道号
* @param is_guest int 是否为游客身份 1：是 0：否 只有非游客身份才进行调用自定义激活
* @return unknow
 */
func AsyncXiaomiReportEvent(c *gin.Context, device_type, package_name, mark string, is_guest int) {
	if device_type == "" || package_name == "" {
		return
	}
	prefix := "xiaomi"
	//判断是否为小米的渠道
	if strings.Contains(mark, prefix) != false {
		log.Printf("检测当前mark= %v是小米的渠道，进行上报", mark)
		uuid := utils.GetRequestHeaderByName(c, "Uuid")     //获取客户端中的UUID
		userinfo, _ := user_service.GetUserByDeviceid(uuid) //根据UUID获取用户信息
		log.Printf("用户的上报状态信息：%v\n", userinfo.ReportStatus)
		if userinfo.Id != 0 {
			//获取对应的渠道信息
			versionInfo, _ := version_service.GetVersionByQdh(device_type, package_name, mark)
			log.Println("小米渠道：【自定义新增激活】上报流程 start.....")
			//上报小米的的新增数据
			utils.AsyncXiaomoReport(c, versionInfo, package_name, mark, "APP_ACTIVE_NEW")
			//只有是非游客状态才进行调用
			if is_guest == 0 {
				log.Printf("每次冷启动来自动上报小米推送事件：【自定义激活】事件,下面开始流程.............")
				//每次启动来进行上报
				utils.AsyncXiaomoReport(c, versionInfo, package_name, mark, "APP_ACTIVE")
			}
		} else {
			log.Printf("当前用户设备ID= %v 尚未存储在库里面 ---记录不存在\n", uuid)
		}
	} else {
		log.Printf("当前渠道= %v 不是小米的，不需要上报\n", mark)
	}
}

// 百度点击统计记录
func SyncBaiduClickData(c *gin.Context, callbackUrl string) {
	params := c.Request.URL.Query()
	imei := params.Get("imei_md5")    //imei信息
	oaid := params.Get("oaid")        //唯一的标识ID
	accountId := params.Get("userid") //账户id
	aid := params.Get("aid")          //创意ID
	ip := params.Get("ip")            //所属IP
	userAgent := params.Get("ua")     //用户客户端标识
	//转换aid参数
	aidInt, err := strconv.Atoi(aid)
	if err != nil {
		aidInt = 0
	}
	global.BaiduClicklog.Info("oaid = ", oaid, " imei = ", imei, " aid = ", aid, " ip = ", ip, " accountId = ", accountId)
	global.BaiduClicklog.Info("ua = ", userAgent)
	channel := config.BAIDU_CHANNEL //获取百度的渠道统计
	callback, _ := GetShenmaClickInfo(oaid, imei, channel)
	if callback.Id != 0 {
		global.BaiduClicklog.Info("【baidu】此设备oaid = ", oaid, "已经记录点击记录了，不需要重复保存哈")
	} else {
		accountId64, err := strconv.Atoi(accountId)
		if err != nil {
			accountId64 = 0
		}
		//待添加的数据
		sadd := models.ClickCallback{
			Oaid:            oaid,            //oaid编号
			Imei:            imei,            //ime编码
			Ua:              userAgent,       //UA客户端标识
			Callback:        callbackUrl,     //获取回调地址
			Ip:              ip,              //所属IP
			PlaceType:       0,               //投放类型默认为0
			Channel:         channel,         //当前渠道
			AdvertisementId: accountId64,     //广告ID对应账户ID
			CreativeId:      aidInt,          //广告创意ID
			CreatedAt:       utils.GetUnix(), //创建事件
		}
		global.SmClicklog.Info("回调callbak地址", sadd.Callback)
		var dataRes *models.ClickCallback
		dataRes = &sadd
		//添加神马平台的关联数据
		err = dataRes.AddCallbackShenma()
		global.SmClicklog.Info("记录百度的回调数据入库啦，去查看吧")
		if err != nil {
			global.SmClicklog.Error("shenma error insert data:", err)
		}
	}
}

/*
* @note 上报神马的点击数据同步写在这里面判断
* @param c gin.Contentxt 请求数据
* @param callbackUrl 回调函数
* @return object ,total , err
 */
func SyncShenmaClickData(c *gin.Context, callbackUrl string) {
	params := c.Request.URL.Query()
	oaid := params.Get("oaid") //oaid参数
	imei := params.Get("imei") //imei参数
	ip := params.Get("ip")
	clickTime := params.Get("time") //处理获取的点击时间
	acid := params.Get("acid")      //广告主账号ID对应我们库里的advertiser_id
	cid := params.Get("cid")        //广告创意ID creative_id
	gid := params.Get("gid")        //广告组ID group_id
	aid := params.Get("aid")        //广告ID，对应我们库里的advertisement_id
	clickTimeNew, err := strconv.ParseInt(clickTime, 10, 64)
	if err != nil {
		clickTimeNew = 0
	}
	cidInt, err := strconv.Atoi(cid)
	if err != nil {
		cidInt = 0
	}
	gidInt, err := strconv.Atoi(gid)
	if err != nil {
		gidInt = 0
	}
	aidInt, err := strconv.Atoi(aid)
	if err != nil {
		aidInt = 0
	}
	userAgent := params.Get("ua") //获取用户的UA信息
	global.SmClicklog.Info("oaid = ", oaid, " imei = ", imei, " ip = ", ip, " Clicktime = ", clickTime, " acid = ", acid, " cid =", cid, " aid =", aid, " gid = ", gid)
	global.SmClicklog.Info("ua = ", userAgent)
	channel := config.SM_CHALLEN
	//获取激活的记录状态信息
	callback, _ := GetShenmaClickInfo(oaid, imei, channel)
	if callback.Id != 0 {
		global.SmClicklog.Info("此设备oaid=", oaid, "已经记录点击记录了，不需要重复保存哈")
	} else {
		global.SmClicklog.Info("同步神马点击数据 oaid=", oaid, "流程开始")
		//这里进行首次激活上报
		sadd := models.ClickCallback{
			Oaid:            oaid,            //oaid编号
			Imei:            imei,            //ime编码
			AdvertiserId:    acid,            //广告主ID，对应对方的广告账户ID
			AdvertisementId: aidInt,          //广告ID
			CreativeId:      cidInt,          //广告创意ID
			GroupId:         gidInt,          //群组ID
			Ua:              userAgent,       //UA客户端标识
			ClickTime:       clickTimeNew,    //点击时间
			Callback:        callbackUrl,     //获取回调地址
			Ip:              ip,              //所属IP
			PlaceType:       0,               //投放类型默认为0
			Channel:         channel,         //当前渠道
			CreatedAt:       utils.GetUnix(), //创建事件
		}
		global.SmClicklog.Info("IP地址", sadd.Ip)
		global.SmClicklog.Info("回调callbak地址", sadd.Callback)
		var dataRes *models.ClickCallback
		dataRes = &sadd
		//添加神马平台的关联数据
		err = dataRes.AddCallbackShenma()
		global.SmClicklog.Info("记录神马平台的回调数据入库啦，去查看吧")
		if err != nil {
			global.SmClicklog.Error("shenma error insert data:", err)
		}
	}
}

// 上报神马平台
func SmReportEvent(c *gin.Context, package_name, mark, oaid, imei, channel, cvType string) {
	global.SmClicklog.Info("package_name:", package_name, "  mark:", mark, "  oaid:", oaid, "  imei:", imei, "channel:", cvType)
	prefix := "sm"
	if strings.Contains(mark, prefix) != false {
		userinfo, err := user_service.GetUserByOaid(oaid)
		if err != nil {
			return
		}
		if userinfo.Id > 0 && userinfo.ReportStatus == 0 {
			//只有用户首次激活才进行上报
			callback := models.ClickCallback{}
			if oaid != "" {
				callback, err = callback.GetCountByImeiAndOaidType(oaid, channel)
				if err != nil {
					global.SmClicklog.Error("GetCountByImeiAndOaid error:", err)
					return
				}
			} else {
				if imei == "" {
					global.SmClicklog.Error("imei为空，无法查询sm回调记录")
					return
				}
				callback, err = callback.GetCountByImeiAndOaidType(imei, cvType)
				if err != nil {
					global.SmClicklog.Error("【sm】 GetCountByImeiAndOaid 2 error:", err)
					return
				}
			}
			if callback.Id != 0 {
				callBackUrl := callback.Callback //获取神马平台的回调地址
				if callBackUrl != "" {
					//激活判断
					if cvType == "active" {
						//激活上报
						global.SmClicklog.Info("首次激活上报，开始报送激活数据，开始下面的流程...........")
						global.SmClicklog.Info("******************神马平台激活回调地址**************", callBackUrl)
						//上报激活地址
						activateUrl := utils.GetReplaceChaojihuiCallbak(callBackUrl, 1, "")
						global.SmClicklog.Info("【shenmapingtai】转换后的激活上报地址 : ", activateUrl)
						//请求激活的回调
						utils.GetBaiduResponse(activateUrl)
						//激活完成后更新用户状态
						userId := userinfo.Id //用户id
						mapData := make(map[string]interface{})
						mapData["report_status"] = 1 //已上报标记
						mapData["uptime"] = utils.GetUnix()
						_ = user_service.UpdateUserByUserId(userId, mapData) //根据用户ID进行更新
					}
				} else {
					global.SmClicklog.Error("sm 激活地址为空，暂时不需要上报")
				}
			} else {
				global.SmClicklog.Error("sm 回调记录不存在:", callback)
			}
		} else {
			global.SmClicklog.Error("channel = ", channel, " 已经上报过激活回传了，不需要再次上报 user_id: ", userinfo.Id, " status: ", userinfo.ReportStatus)
		}
		//留存上报
		if cvType == "retention" {
			callback, _ := GetShenmaClickInfo(oaid, imei, channel)
			if callback.Id != 0 {
				callBackUrl := callback.Callback //获取他的拼装地址
				//获取留存地址上报
				retainUrl := utils.GetReplaceChaojihuiCallbak(callBackUrl, 1001, "")
				global.SmClicklog.Info("【shenmapingtai】转换后的留存次日上报地址 : ", retainUrl)
				utils.GetBaiduResponse(retainUrl) //请求留存
			} else {
				global.SmClicklog.Info("【神马平台监测的记录信息为空哦 ")
			}
		}
	}
}

func VivoReportedEvent(c *gin.Context, package_name, mark, oaid, imei, cvType string) {
	global.VivoClicklog.Info("package_name:", package_name, "  mark:", mark, "  oaid:", oaid, "  imei:", imei)
	prefix := "vivo"
	if strings.Contains(mark, prefix) != false {
		userinfo, err := user_service.GetUserByOaid(oaid)
		if err != nil {
			return
		}
		if (userinfo.Id > 0 && userinfo.ReportStatus == 0) || cvType == config.VIVO_RETENTION_1 {
			//根据oaid and imei查询vivo回调记录
			callback := models.ClickCallback{}
			if oaid != "" {
				callback, err = callback.GetCountByImeiAndOaidType(oaid, "vivo")
				if err != nil {
					global.VivoClicklog.Error("GetCountByImeiAndOaid error:", err)
					return
				}
			} else {
				if imei == "" {
					global.VivoClicklog.Error("imei为空，无法查询vivo回调记录")
					return
				}
				callback, err = callback.GetCountByImeiAndOaidType(imei, "vivo")
				if err != nil {
					global.VivoClicklog.Error("GetCountByImeiAndOaid 2 error:", err)
					return
				}
			}

			if callback.Id != 0 {

				var userid string
				var UserIdType string
				if oaid == "" {
					userid = imei
					if len(userid) >= 17 {
						UserIdType = "IMEI_MD5"
					} else {
						UserIdType = "IMEI"
					}
				} else {
					userid = oaid
					//长度大于10，则认为是MD5
					if len(userid) >= 32 {
						UserIdType = "OAID"
					} else {
						UserIdType = "OAID_MD5"
					}
				}
				global.VivoClicklog.Info("UserIdType->", UserIdType)
				global.VivoClicklog.Info("userid->", userid)
				dataList := models.DataList{
					UserIdType: UserIdType,
					UserId:     userid,
					CvType:     cvType,
					CvTime:     time.Now().UnixMilli(),
					CreativeId: strconv.Itoa(callback.CreativeId),
					RequestId:  callback.RequestId,
					ExtParam: models.ExtParam{
						PayAmount: "0",
					},
				}

				returnData := models.RequestData{
					DataList: []models.DataList{dataList},
				}
				returnData.PkgName = "com.ttyhsgapp.kxs"
				returnData.SrcId = "ds-202409046187"
				returnData.SrcType = "APP"
				returnData.DataFrom = "1"
				timestamp := strconv.Itoa(int(time.Now().UnixMilli()))
				nonce := utils.GetRandomString(10)
				advertiserId := "fd18da3fe8f79d93fc6e"
				accessToken, _ := GetVivoAccessToken()
				url := "https://marketing-api.vivo.com.cn/openapi/v1/advertiser/behavior/upload"

				re, err := VivoReport(returnData, url, accessToken, timestamp, nonce, advertiserId)
				if err != nil {
					global.VivoClicklog.Error("VivoReport error :", err)
					return
				}
				if re.Code != 0 {
					global.VivoClicklog.Error("VivoReport error code:", re.Code, "message:", re.Message)
					return
				}
				global.VivoClicklog.Info("VivoReport success:", re)
				if cvType == config.VIVO_ACTIVATION { //留存无需更改
					userId := userinfo.Id //用户id
					mapData := make(map[string]interface{})
					mapData["report_status"] = 1 //已上报标记
					mapData["uptime"] = utils.GetUnix()
					_ = user_service.UpdateUserByUserId(userId, mapData) //根据用户ID进行更新
				}
			} else {
				global.VivoClicklog.Error("vivo回调记录不存在:", callback)
			}
		} else {
			global.VivoClicklog.Error("已经上报过了，不需要再次上报 user_id:", userinfo.Id, "   status:", userinfo.ReportStatus)
		}
	}
}
func VivoReport(returnData models.RequestData, url string, accessToken string, timestamp string, nonce string, advertiserId string) (models.Response, error) {
	// 序列化为 JSON
	returnJSON, err := json.Marshal(returnData)
	if err != nil {
		global.VivoClicklog.Error("JSON marshal error:", err)
		return models.Response{}, err
	}
	global.VivoClicklog.Info("body", string(returnJSON))
	// 创建请求体
	body := bytes.NewBuffer(returnJSON)

	// 构建请求 URL
	fullURL := fmt.Sprintf("%s?access_token=%s&timestamp=%s&nonce=%s&advertiser_id=%s",
		url, accessToken, timestamp, nonce, advertiserId)
	global.VivoClicklog.Info("fullUrl", fullURL)
	// 创建 HTTP 请求
	re, err := http.NewRequest("POST", fullURL, body)
	if err != nil {
		global.VivoClicklog.Error("Request error:", err)
		return models.Response{}, err
	}

	// 设置请求头
	re.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(re)
	if err != nil {
		fmt.Println("Request execution error:", err)
		return models.Response{}, err
	}
	defer resp.Body.Close()

	// 读取响应体
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		global.VivoClicklog.Error("Read response error:", err)
		return models.Response{}, err
	}
	//{"code":70100,"message":"请求参数错误，详情列表不能为空!"}
	response := models.Response{}
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		global.VivoClicklog.Error("JSON unmarshal error:", err)
	}

	return response, nil
}

// 获取神马平台的统计信息
func GetShenmaClickInfo(oaid, imei, channel string) (info models.ClickCallback, err error) {
	callback := models.ClickCallback{}
	if oaid != "" { //如果oaid不为空用这个判断
		fmt.Println(oaid, channel)
		callback, err = callback.GetCountByImeiAndOaidType(oaid, channel)
		if err != nil {
			return
		}
		info = callback
	} else if imei != "" { //如果imei不为空用这个判断
		callback, err = callback.GetCountByImeiAndOaidType(imei, channel)
		if err != nil {
			return
		}
		info = callback
	}
	return info, nil

}

// 保存vivo检测数据
func SaveVivoClick(req []models.ClickCallback) {
	for _, v := range req {
		var c *models.ClickCallback
		c = &v
		c.Channel = "vivo"
		//查询请求id是否存在
		is_requestId, _ := c.GetCountByRequestId(c.RequestId)
		if is_requestId != 0 {
			global.VivoClicklog.Error("重复请求:", c.RequestId)
			continue
		}
		err := c.AddCallback()
		if err != nil {
			global.VivoClicklog.Error("error insert data:", err)
			return
		}
	}
}
