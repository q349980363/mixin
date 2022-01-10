package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"gopkg.in/olahol/melody.v1"
)

func init_service_bind1() {
	//留给客户端测试是否拥有权限

	r.GET("/api/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"state":   true,
			"message": "233",
		})
	})
	r.GET("/admin/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"state":   true,
			"message": "233",
		})
	})
	init_service_bind_group()
	init_service_bind_user()
	init_service_bind_friends()

	r.GET("/api/get_systemchat", func(c *gin.Context) {
		userInfo := c.Keys["userInfo"].(UserInfo)
		var systemChat []SystemChat
		db.Find(&systemChat, &SystemChat{
			UserName: userInfo.UserName,
		})
		db.Model(&Friends{}).Where(&Friends{
			UserName: userInfo.UserName,
			Target:   "系统消息",
		}).Update("Unread", 0)
		c.JSON(http.StatusOK, gin.H{
			"state":   true,
			"list":    systemChat,
			"message": "",
		})
	})

	r.POST("/api/operation_systemchat/:id", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		json := make(map[string]string)
		c.ShouldBind(&json)
		userInfo := c.Keys["userInfo"].(UserInfo)

		error := operationSystemChat(userInfo, id, json)
		if error != nil {
			c.JSON(http.StatusOK, gin.H{
				"state":   false,
				"message": error.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"state":   true,
			"message": "操作成功.",
		})
	})

}

func http_empty(c *gin.Context) {
	db.Where("true").Delete(Chat{})
	m.Broadcast([]byte("empty"))
}

func http_login(c *gin.Context) {
	json := make(map[string]string)
	c.ShouldBind(&json)

	_log("用户登录")
	_log("userName:" + json["userName"])
	_log("passWord:" + json["passWord"])

	if json["userName"] == "" || json["passWord"] == "" {
		_log("用户名或密码不能为空")
		c.JSON(http.StatusOK, gin.H{
			"state":   false,
			"message": "用户名或密码不能为空",
		})
		return
	}

	var userInfo UserInfo
	if db.First(&userInfo, &UserInfo{UserName: json["userName"]}).RecordNotFound() {
		_log("用户不存在")
		c.JSON(http.StatusOK, gin.H{
			"state":   false,
			"message": "用户不存在",
		})
		return
	}
	if json["passWord"] != userInfo.PassWord {
		_log("密码不正确")
		c.JSON(http.StatusOK, gin.H{
			"state":   false,
			"message": "密码不正确",
		})
		return
	}
	_log("登录成功")
	c.JSON(http.StatusOK, gin.H{
		"state":   true,
		"message": "登录成功",
		"data":    userInfo,
	})
}

func http_register(c *gin.Context) {
	json := make(map[string]string)
	c.ShouldBind(&json)

	_log("注册用户")
	_log("userName:" + json["userName"])
	_log("passWord:" + json["passWord"])
	if json["userName"] == "" || json["passWord"] == "" {
		_log("用户名或密码不能为空")
		c.JSON(http.StatusOK, gin.H{
			"state":   false,
			"message": "用户名或密码不能为空",
		})
		return
	}
	if !db.First(&UserInfo{}, &UserInfo{UserName: json["userName"]}).RecordNotFound() {
		_log("用户已存在")
		c.JSON(http.StatusOK, gin.H{
			"state":   false,
			"message": "用户已存在",
		})
		return
	}

	userInfo := UserInfo{
		UserName: json["userName"],
		PassWord: json["passWord"],
		ToKen:    randomString(32),
	}
	db.Create(&userInfo)
	c.JSON(http.StatusOK, gin.H{
		"state":   true,
		"message": "注册成功",
		"data":    userInfo,
	})
}

func HandleConnect(s *melody.Session, userInfo UserInfo) {
	_log("用户连接 " + userInfo.UserName)
	{
		//改为世界频道广播
		// globals
		ip, ipAddr := getIpAddr(s)
		json, _ := json.Marshal(Chat{Data: "用户 [" + userInfo.UserName + "] 进入", Model: gorm.Model{CreatedAt: time.Now()}, UserName: "系统消息", Ip: ip, IpAddr: ipAddr})
		m.BroadcastOthers(json, s)
	}
	var chats []Chat
	db.Limit("-100").Find(&chats)
	json, _ := json.Marshal(gin.H{
		"type": "globals",
		"data": chats,
	})
	s.Write(json)
	sessionAddGroup(s, "user_"+userInfo.UserName)

	var groups []GroupRelation
	db.Find(&groups, &GroupRelation{
		UserName: userInfo.UserName,
	})
	for _, v := range groups {
		sessionAddGroup(s, "group_"+v.Target)
	}
}

func HandleMessage(s *melody.Session, bytes []byte, msg string, userInfo UserInfo) {
	_log("收到消息[" + msg + "]")
	if msg == "" {
		return
	}
	requestData := make(map[string]string)
	_ = json.Unmarshal(bytes, &requestData)

	target := requestData["target"]
	ip, ipAddr := getIpAddr(s)
	fmt.Println(requestData)
	if requestData["txt"] == "" {
		_log("消息内容为空")
		return
	}

	switch requestData["type"] {
	case "friends":
		{
			sendFriendsChat(&FriendsChat{
				Type:     "txt",
				UserName: userInfo.UserName,
				Target:   target,
				Data:     requestData["txt"],
				Ip:       ip,
				IpAddr:   ipAddr,
			})
			return
		}
	case "group":
		{
			sendGroupChat(&GroupChat{
				Type:     "txt",
				UserName: userInfo.UserName,
				Target:   target,
				Data:     requestData["txt"],
				Ip:       ip,
				IpAddr:   ipAddr,
			})
			return
		}
	case "global":
		{
			sendGlobalChat(&GroupChat{
				Type:     "txt",
				UserName: userInfo.UserName,
				Target:   "global",
				Data:     requestData["txt"],
				Ip:       ip,
				IpAddr:   ipAddr,
			})
			return
		}
	case "":
	}

	chat := &Chat{Type: "txt", UserName: userInfo.UserName, Data: msg, Ip: ip, IpAddr: ipAddr}
	db.Create(chat)
	json, _ := json.Marshal(gin.H{
		"type": "global",
		"data": chat,
	})
	m.Broadcast(json)
}
