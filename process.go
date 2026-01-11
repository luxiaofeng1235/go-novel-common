package main

import (
	"flag"
	"fmt"
	"go-novel/utils"
	"log"
	"os/exec"
	"strconv"
	"sync"
)

func main() {
	// 定义要执行的参数列表
	start, end := getStartEnd()
	var wg1 sync.WaitGroup
	for i := start; i <= end; i++ {
		wg1.Add(1)
		go goexec(i, &wg1)
	}
	wg1.Wait()
	log.Println("执行完了")
}

func goexec(i int, wg1 *sync.WaitGroup) {
	defer wg1.Done()
	//构建命令对象
	cmd := exec.Command("./tablepage", "-page="+strconv.Itoa(i))

	// 执行命令并获取输出
	output, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	// 输出命令执行结果
	fmt.Printf("Page %d output:\n%s\n", i, string(output))
}

func getStartEnd() (start, end int) {
	var startStr, endStr, ip, port, username, passwd string
	flag.StringVar(&startStr, "start", "1", "default :start")
	flag.StringVar(&endStr, "end", "1", "default :end")

	flag.StringVar(&ip, "ip", "", "default :ip")
	flag.StringVar(&port, "port", "", "default :port")
	flag.StringVar(&username, "username", "", "default :username")
	flag.StringVar(&passwd, "passwd", "", "default :passwd")
	flag.Parse()

	utils.S5Domain = ip
	utils.S5Port = port
	utils.S5Username = username
	utils.S5Passwd = passwd
	start, _ = strconv.Atoi(startStr)
	end, _ = strconv.Atoi(endStr)
	return
}
