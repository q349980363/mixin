package main

import (
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func _log(txt ...interface{}) {
	currentTime := time.Now()
	for _, value := range txt {
		fmt.Println("[MiXin] "+currentTime.Format("[2006/1/02 15:04]    "), value)
	}
}

// Returns an int >= min, < max
func randomInt(min, max int) int {
	return min + rand.Intn(max-min)
}

// Generate a random string of A-Z chars with len = l
func randomString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		bytes[i] = byte(randomInt(65, 90))
	}
	return string(bytes)
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", "*") // 可将将 * 替换为指定的域名
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Header("Access-Control-Allow-Headers", "token,Origin, X-Requested-With, Content-Type, Accept, Authorization")
		c.Header("Access-Control-Expose-Headers", "token,Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()

	}
}

func Ips() (map[string]string, error) {
	ips := make(map[string]string)
	//返回 interface 结构体对象的列表，包含了全部网卡信息
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	//遍历全部网卡
	for _, i := range interfaces {
		if i.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if i.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		// Addrs() 方法返回一个网卡上全部的IP列表
		address, err := i.Addrs()
		if err != nil {
			return nil, err
		}

		//遍历一个网卡上全部的IP列表，组合为一个字符串，放入对应网卡名称的map中
		for _, v := range address {
			if ipnet, ok := v.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					ips[i.Name] += ipnet.IP.String()
				}
			}
		}
	}
	return ips, nil
}

func init_ip() {
	_log("http://127.0.0.1:8080/")
	var ips = make(map[string]string)
	ips, _ = Ips()

	//key是网卡名称，value是网卡IP
	for k, v := range ips {
		_log("http://" + v + ":8080/    " + k)
	}
}
