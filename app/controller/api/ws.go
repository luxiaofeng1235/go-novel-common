package api

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
	"go-novel/app/models"
	"go-novel/app/service/api/book_service"
	"go-novel/app/service/api/notice_service"
	"go-novel/global"
	"go-novel/utils"
)

type Ws struct{}

func (ws *Ws) HandleRequest(c *gin.Context) {
	err := global.Ws.HandleRequest(c.Writer, c.Request)
	if err != nil {
		global.Errlog.Error(err.Error())
	}
}

func WsHandleMsg() {
	global.Ws.HandleConnect(func(s *melody.Session) {
		if s.Request.URL.Path == utils.ApiWs {
			userId := s.Request.URL.Query().Get("user_id")
			utils.UserList[utils.WsRemoteIp(s)] = utils.FormatInt64(userId)
			utils.NodeList[utils.WsRemoteIp(s)] = s
			global.Wslog.Infoln(s.Request.URL.Path, utils.WsRemoteIp(s), userId, "在线")
		}
	})

	global.Ws.HandleDisconnect(func(s *melody.Session) {
		ip := utils.WsRemoteIp(s)
		if s.Request.URL.Path == utils.ApiWs {
			delete(utils.UserList, ip)
			delete(utils.NodeList, ip)
		}
		global.Wslog.Infoln(s.Request.URL.Path, ip, "离线")
	})

	global.Ws.HandlePong(func(s *melody.Session) {
		//fmt.Println("pong", s.Request.UserAgent())
		//s.Write([]byte(fmt.Sprint(websocket.PingMessage)))
	})

	global.Ws.HandleMessage(func(s *melody.Session, message []byte) {
		global.Wslog.Infoln("收到消息", utils.WsRemoteIp(s), s.Request.URL.Path, string(message))
		if len(message) < 4 {
			utils.WsFail(s, nil, "", "", "未知的请求方法")
			return
		}
		wsApi := new(Ws)
		method := string(message[:4])
		msg := message[4:]
		//log.Println(msg)
		//0x02 //init初始化返回错误信息
		//0x03 //获取未读数命令
		//0x04 //发送未读数
		//0x05 //获取小说章节
		//0x06 //发送小说章节
		if s.Request.URL.Path == utils.ApiWs {
			switch method {
			case utils.Zone:
				//wsApi.AddNode(s, msg)
			case utils.Zthree:
				wsApi.SendNoReadCount(s, msg)
			case utils.Zfive:
				wsApi.SendChapters(s, msg)
			default:
				utils.WsFail(s, nil, method, "", "未知的请求协议")
			}
		}
	})
}

func (ws *Ws) SendNoReadCount(s *melody.Session, msg []byte) {
	var req models.SendNoReadCount
	// 参数绑定
	var err error
	if len(msg) > 0 && json.Valid([]byte(msg)) {
		err = json.Unmarshal(msg, &req)
	}
	userId := utils.UserList[utils.WsRemoteIp(s)]
	if userId <= 0 {
		err = fmt.Errorf("%v", "获取user_id错误")
		utils.WsFail(s, err, utils.Zfour, req.Key, "获取用户信息失败")
		return
	}

	totalCount := notice_service.GetNoReadNotifyCount(userId, "") //关注消息不算在内
	noticeCount := notice_service.GetNoReadNotifyCount(userId, utils.Notice)
	praiseCount := notice_service.GetNoReadNotifyCount(userId, utils.Praise)
	commentCount := notice_service.GetNoReadNotifyCount(userId, utils.Comment)
	global.Requestlog.Info("user_id:", userId, " totalCount:", totalCount, "noticeCount", ":", noticeCount, "praiseCount:", praiseCount, "commentCount:", commentCount)
	res := gin.H{
		"totalCount":   totalCount,
		"noticeCount":  noticeCount,
		"praiseCount":  praiseCount,
		"commentCount": commentCount,
	}
	utils.WsSuccess(s, utils.Zfour, req.Key, res, "ok")
	return
}

func (ws *Ws) SendChapters(s *melody.Session, msg []byte) {
	var req models.SendChapters
	// 参数绑定
	var err error
	if len(msg) > 0 && json.Valid([]byte(msg)) {
		err = json.Unmarshal(msg, &req)
	}
	userId := utils.UserList[utils.WsRemoteIp(s)]
	if userId <= 0 {
		err = fmt.Errorf("%v", "获取user_id错误")
		utils.WsFail(s, err, utils.Zsix, req.Key, "获取用户信息失败")
		return
	}

	chapterRes, bookUrl, chapterTextReg, err := book_service.Chapter(req.BookId, req.SourceId, req.Sort)
	if err != nil {
		utils.WsFail(s, err, utils.Zsix, req.Key, err.Error())
		return
	}

	dataCount := len(chapterRes)
	batchSize := 2000
	// 计算批次数量
	numBatches := dataCount / batchSize
	if dataCount%batchSize != 0 {
		numBatches++
	}
	// 按批次发送数据
	for i := 0; i < numBatches; i++ {
		// 计算当前批次的起始索引和结束索引
		startIndex := i * batchSize
		endIndex := (i + 1) * batchSize
		if endIndex > dataCount {
			endIndex = dataCount
		}

		// 获取当前批次的数据切片
		batchData := chapterRes[startIndex:endIndex]
		res := gin.H{
			"chapters":       batchData,
			"bookUrl":        bookUrl,
			"chapterTextReg": chapterTextReg,
		}
		utils.WsSuccess(s, utils.Zsix, req.Key, res, "ok")

	}
	return
}
