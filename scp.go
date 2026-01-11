package main

import (
	"fmt"
	"go-novel/db"
	"go-novel/utils"
	"log"
	"os/exec"
	"time"
)

func main() {
	addr, passwd, defaultdb := db.GetRedis()
	addr = "206.119.66.243:6379"
	passwd = "o9kHvO95bP"
	defaultdb = 0
	db.InitRedis(addr, passwd, defaultdb)
	db.InitZapLog()
	log.Println("开始执行")
	for {
		utils.RefreshS5Proxys()
		source := "E:\\go\\src\\go-novel\\S5Proxys.json"
		destination := "root@103.36.91.35:/www/wwwroot/novel"
		//remoteCommand := "cd /www/wwwroot/novel && sh starttemp0.sh && sh starttemp1.sh && sh starttemp2.sh" +
		//	"&& sh starttemp3.sh && sh starttemp4.sh && sh starttemp5.sh && sh starttemp6.sh && sh starttemp7.sh " +
		//	"&& sh starttemp8.sh && sh starttemp9.sh && sh starttemp10.sh && sh starttemp3.sh && sh starttemp11.sh " +
		//	"&& sh starttemp12.sh && sh starttemp13.sh && sh starttemp14.sh && sh starttemp15.sh && sh starttemp16.sh " +
		//	"&& sh starttemp17.sh && sh starttemp18.sh" +
		//	"&& sh starttable1.sh && sh starttable2.sh && sh starttable3.sh && sh starttable4.sh && sh starttable5.sh " +
		//	"&& sh starttable6.sh && sh starttable7.sh && sh starttable8.sh && sh starttable9.sh && sh starttable10.sh " +
		//	"&& sh starttable11.sh && sh starttable12.sh && sh starttable13.sh && sh starttable14.sh && sh starttable15.sh " +
		//	"&& sh starttable16.sh && sh starttable17.sh && sh starttable18.sh && sh starttable19.sh" +
		//	"&& sh starttable21.sh && sh starttable22.sh && sh starttable23.sh && sh starttable24.sh && sh starttable25.sh " +
		//	"&& sh starttable26.sh && sh starttable27.sh && sh starttable28.sh && sh starttable29.sh && sh starttable30.sh " +
		//	"&& sh starttable31.sh && sh starttable32.sh && sh starttable33.sh && sh starttable34.sh && sh starttable35.sh " +
		//	"&& sh starttable36.sh && sh starttable37.sh && sh starttable38.sh"

		// 构建 SCP 命令
		cmd := exec.Command("scp", source, destination)

		// 运行 SCP 命令并等待完成
		err := cmd.Run()
		if err != nil {
			log.Println(err)
			continue
		}

		fmt.Println("File copied successfully to remote server.")

		// 构建 SSH 命令
		//sshCmd := exec.Command("ssh", "root@103.36.91.35", remoteCommand)
		//// 运行 SSH 命令并等待完成
		//err = sshCmd.Run()
		//if err != nil {
		//	log.Println(err)
		//	continue
		//}
		//fmt.Println("Remote process restarted.")
		time.Sleep(time.Minute * 15)
	}
}
