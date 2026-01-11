package zaplog

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"go-novel/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"net/http"
	"os"
	"time"
)

// ZincSearch 配置
var (
	ZincSearchURLFormat = "http://127.0.0.1:4080/api/%v/_doc" // 格式化字符串
	ZincSearchTimeout   = 5 * time.Second                     // HTTP 请求超时
	ZincSearchUsername  = "admin"
	ZincSearchPassword  = "SInR5cCI6IkpXV#25"
)

// SendLogToZincSearch 发送日志到 ZincSearch
func SendLogToZincSearch(fields map[string]interface{}, index string) {
	// 创建日志条目
	logEntry := map[string]interface{}{
		"@timestamp": time.Now().Format(time.RFC3339),
	}

	for k, v := range fields {
		logEntry[k] = v
	}

	data, err := json.Marshal(logEntry)
	if err != nil {
		fmt.Printf("Failed to marshal log entry: %v\n", err)
		return
	}

	// 使用局部变量构建格式化后的 URL，避免修改全局变量
	formattedURL := fmt.Sprintf(ZincSearchURLFormat, index)

	// 创建 HTTP 请求
	req, err := http.NewRequest("POST", formattedURL, bytes.NewBuffer(data))
	if err != nil {
		fmt.Printf("Failed to create HTTP request: %v\n", err)
		return
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")

	// 设置基本认证头
	credentials := fmt.Sprintf("%s:%s", ZincSearchUsername, ZincSearchPassword)
	encodedCredentials := base64.StdEncoding.EncodeToString([]byte(credentials))
	authHeader := fmt.Sprintf("Basic %s", encodedCredentials)
	req.Header.Set("Authorization", authHeader)

	// 创建 HTTP 客户端
	client := &http.Client{Timeout: ZincSearchTimeout}

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Failed to send log entry: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode >= 300 {
		fmt.Printf("Received non-success status code: %d\n", resp.StatusCode)
	}
}

// 初始化日志
func InitLogger(LogSavePath string) {
	err := utils.IsNotExistMkDir(LogSavePath)
	if err != nil {
		log.Fatalln("创建日志目录失败：", err.Error())
	}
}

func getLogFileFullPath(path string) string {
	logFileName := fmt.Sprintf("%s.%s", time.Now().Format("20060102"), "log")
	return fmt.Sprintf("%s/%s", path, logFileName)
}

// 日志配置
func LogConfig(path string) *zap.SugaredLogger {
	hook := lumberjack.Logger{
		Filename:   getLogFileFullPath(path),
		MaxSize:    64, // megabytes
		MaxBackups: 3,
		MaxAge:     5,     //days
		Compress:   false, // disabled by default
	}

	fileEncoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "name",
		CallerKey:      "line",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		FunctionKey:    "func",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}

	stdEncoderConfig := zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "name",
		CallerKey:      "line",
		StacktraceKey:  "stacktrace",
		FunctionKey:    "func",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeName:     zapcore.FullNameEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}

	fileEncoder := zapcore.NewConsoleEncoder(fileEncoderConfig)
	stdEncoder := zapcore.NewConsoleEncoder(stdEncoderConfig)

	fileWriter := zapcore.NewMultiWriteSyncer(zapcore.AddSync(&hook))
	stdWriter := zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout))

	debugLevel := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level >= zapcore.DebugLevel
	})

	infoLevel := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level >= zapcore.InfoLevel
	})

	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, fileWriter, infoLevel),
		zapcore.NewCore(stdEncoder, stdWriter, debugLevel),
	)

	zaplog := zap.New(core)
	slog := zaplog.Sugar()

	err := utils.IsNotExistMkDir(path)
	if err != nil {
		log.Println("创建日志目录失败：", err.Error())
		return slog
	}

	return slog
}
