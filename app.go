package main

import (
	"math/rand"
	"sync"
	"time"
)

//全局锁
var lock = new(sync.Mutex)

func main() {

	//初始化随机数发生器
	rand.Seed(time.Now().UnixNano())

	init_qqwry()

	init_db()
	defer db.Close()
	init_service()

	_log("服务器启动,端口8080")
	//获取本机所有IP并输出提示信息
	init_ip()
	r.Run(":8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

	for i := 0; i < 10; i++ {
		_log("##################服务器关闭##################")
	}
}
