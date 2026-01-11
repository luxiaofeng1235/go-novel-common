package setting_service

import (
	"bytes"
	"fmt"
	"github.com/jlaffaye/ftp"
	"log"
	"path/filepath"
	"strings"
)

// FTPClient 结构体用于封装 FTP 连接
type FTPClient struct {
	conn *ftp.ServerConn
}

// NewFTPClient 函数用于创建一个新的 FTPClient 实例
func NewFTPClient(server, user, password string) (*FTPClient, error) {
	conn, err := ftp.Dial(server)
	if err != nil {
		return nil, err
	}
	//登录
	err = conn.Login(user, password)
	if err != nil {
		return nil, err
	}
	//// 关闭FTP连接
	//if err = conn.Quit(); err != nil {
	//	log.Fatal(err)
	//}
	return &FTPClient{conn: conn}, nil
}

// UploadFtpFile 上传ftp文件
// content:文本内容 filename:上传文件名
func (c *FTPClient) UploadFtpFile(content, filename string) error {
	if filename == "" {
		return nil
	}
	//创建文远程文件夹目录
	currentDir := filepath.Dir(filename)
	currentDir = strings.ReplaceAll(currentDir, "\\", "/")
	//判断文件不存在的时候自动创建文件夹
	if currentDir != "." {
		err := c.ChangeDir(currentDir)
		if err != nil {
			//只有目录不存在的时候才进行创建
			log.Printf("当前文件的所在路径: %v", currentDir)
			err = c.CreateFolder(currentDir)
			if err != nil {
				log.Println("创建文件夹失败")
				return nil
			}
		} else {
			log.Printf("当前文件目录 path = %v 已创建，无需重复创建", currentDir)
		}
	}
	//写入文件的存储路径信息
	data := bytes.NewBufferString(content)
	err := c.conn.Stor(filename, data)
	if err != nil {
		panic(err)
		return err
	} else {
		log.Printf("文件存储的路径为：%v\n", filename)
		c.conn.Quit() //退出连接
		return nil
	}
}

// 创建ftp目录文件夹-可以创建多级目录
func (c *FTPClient) CreateFolder(dir string) error {
	err := c.conn.ChangeDir(dir)
	if err != nil {
		log.Println("远程目录", dir, "不存在或无权限访问，接下来会自动创建目录")
		// 如果服务器不支持检查目录是否存在,我们可以尝试逐级创建目录
		parts := strings.Split(dir, "/")
		currentDir := ""
		for _, part := range parts {
			currentDir = filepath.Join(currentDir, part)
			currentDir = strings.ReplaceAll(currentDir, "\\", "/")
			newDir := fmt.Sprintf("/%v", currentDir)
			if newDir != "/" { //解析当前的目录结构
				log.Printf("current dir = %v\n", newDir)
				err = c.conn.ChangeDir(newDir)
				if err != nil { //通过changeDir判断目录目录是否存在
					err = c.conn.MakeDir(newDir)
					if err != nil { //判断是否为550的返回为创建失败
						if err.Error() != "550 Create directory operation failed" {
							log.Printf("failed to create directory '%s': %w", currentDir, err)
						}
						// 如果目录已经存在,我们可以继续创建下一级目录
					}
				} else { //目录存在的情况处理
					log.Printf("this remote folder dir= %v created exists!!\n", newDir)
				}
			}
		}
		log.Printf("this remote folder dir= %v create success!\n", dir)
		return nil
	} else {
		log.Printf("this remote folder dir = %v exists，can access files\n", dir)
		return nil
	}
}

// ChangeDir 方法用于切换到指定目录
func (c *FTPClient) ChangeDir(dir string) error {
	return c.conn.ChangeDir(dir)
}

// Close 方法用于关闭 FTP 连接
func (c *FTPClient) Close() {
	c.conn.Quit()
}

// 同步上传ftp服务器 project_name 项目名称 qdh：渠道 content:设置的内容 typeA :1 隐私协议 2：用户协议
func AsyncFtpUpload(project_name, qdh, content string, typeA int64) (url string) {
	//判断是否为火龙果的域名发那会，通过项目名称来进行判断

	var serverName, userName, passWord, surl, bucket string
	serverName = "103.36.91.36:21"
	if strings.Contains(project_name, "火龙果") != false { //红龙果的ftp账号
		surl = "https://www.huolongyunwu.com" //用火龙果的域名返回
		userName = "hlgyw"
		passWord = "Shitf7YnS8aM5i8T"
		bucket = "huolongguo"
	} else { //快眼的ftp账号
		userName = "kyks"
		passWord = "5JChEJ4Ztk82i8aR"
		surl = "https://www.kuaiyankanshu.vip"
		bucket = "kuaiyan"
	}
	ftpClient, err := NewFTPClient(serverName, userName, passWord)
	//初始化实力，并进行发送请求

	if err != nil {
		log.Fatal("无法创建FTP客户端:", err)
	}
	var htmlFile string
	if typeA != 2 { //用户协议
		htmlFile = "user_agreement.html"
		content = UserAgreementTemplate(content) //编排用户协议内容
	} else { //隐私协议
		htmlFile = "privacy.html"
		content = PrivacyTemplate(content) //编排隐私协议内容
	}
	filename := fmt.Sprintf("/%v/%v/%v", bucket, qdh, htmlFile)
	err = ftpClient.UploadFtpFile(content, filename) //同步服务内容
	url = fmt.Sprintf("%v%v", surl, filename)
	log.Printf("当前访问的url=%v\n", url)
	return url
}

// 重组html页面内容
func PrivacyTemplate(html string) (content string) {
	str := `<html lang="en">
		<head>
			<title>隐私协议</title>
			<meta charset="utf-8" />
			<meta http-equiv="X-UA-Compatible" content="IE=edge" />
			<script>
				var vp = document.createElement("meta"),
					width = window.screen.width,
					design_width = 750,
					scale = width / design_width,
					content =
						"width=" +
						design_width +
						", viewport-fit=cover, user-scalable=no, initial-scale=" +
						scale +
						", maximum-scale=" +
						scale;
				vp.setAttribute("content", content), vp.setAttribute("name", "viewport");
				var head = document.querySelector("head");
				head.appendChild(vp);
			</script>
			<script>
				function getQueryString(name) {
					var reg = new RegExp("(^|&)" + name + "=([^&]*)(&|$)");
					var r = window.location.search.substr(1).match(reg);
					if (r != null) {
						return r[2];
					}
					return null;
				}
			</script>
			<script data-cfasync="false" src="/cdn-cgi/scripts/5c5dd728/cloudflare-static/email-decode.min.js"></script><script>
				var appName = decodeURI(getQueryString("appName"));
				if (appName != "" && appName != undefined && appName != "null") {
					var d = document.getElementsByClassName("appName");
					var len = d.length;
					for (let i = 0; i < len; i++) {
						d.item(i).innerText = appName.split("?")[0];
					}
				}
			</script>
			<script>
				if (/(iPhone|iPad|iPod|iOS)/i.test(navigator.userAgent)) {
				} else {
					document.getElementById("yx1").innerText = "guomushuzi@proton.me";
				}
			</script>
			<link rel="icon" href="data:;base64,=" />
			<style>
				body {
					font-size: 28px;
					box-sizing: border-box;
					padding: 20px;
					margin: 0;
				}
			</style>
		</head>
		<body>`
	str = str + html //组装内容拼接页面内容
	str = str + `</body>
		</html>`
	content = str
	return
}

// 用户协议内魔板
func UserAgreementTemplate(html string) (content string) {
	str := `<html style="font-size: 151.2px">
			  <head>
				<meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
				<meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1" />
				<meta name="renderer" content="webkit" />
				<meta
				  name="viewport"
				  content="width=device-width, user-scalable=no, initial-scale=1.0, maximum-scale=1.0"
				/>
				<meta content="yes" name="apple-mobile-web-app-capable" />
				<meta content="yes" name="apple-touch-fullscreen" />
				<meta content="telephone=no,email=no" name="format-detection" />
				<title>用户使用协议</title>
				<style>
				  html {
					color: #000;
					background: #fff;
					overflow-y: scroll;
					-webkit-text-size-adjust: 100%;
					-ms-text-size-adjust: 100%;
				  }
				  html * {
					outline: 0;
					-webkit-text-size-adjust: none;
					-webkit-tap-highlight-color: transparent;
				  }
				  * {
					margin: 0;
					padding: 0;
				  }
				  h1,h2,h3,h4,h5,h6 {
					font-size: 100%;
					font-weight: 500;
				  }
				  .content {
					width: 91.6%;
					margin: 0 auto;
					padding-top: 24px;
					margin-bottom: 24px;
					font-family: PingFangSC-Regular;
					font-size: 10px;
					color: #797d8b;
					letter-spacing: 0;
					line-height: 20px;
					text-indent: 2em;
					text-align: left;
				  }
				  p{
					margin: 0;
					font-size: 10px;
					margin-bottom: 10px;
				  }
				  .underline {
					text-decoration: underline;
				  }
				</style>
			
				<script type="text/javascript">
				  document.getElementsByTagName("html")[0].style.fontSize =
					window.innerWidth / 10 + "px";
				</script>
				<script>
				  if (/(iPhone|iPad|iPod|iOS)/i.test(navigator.userAgent)) {
				  } else {
					document.getElementById("yx1").innerText = "help_fankui#outlook.com";
				  }
				</script>
			  </head>
			  <body style="font-size: 12px">
					<!-- 页面结构 -->
					<div class="content">`
	str = str + html //拼接用户协议内容
	str = str + `    </div>
			  </body>
			</html>`
	content = str
	return
}
