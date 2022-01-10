package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"gopkg.in/olahol/melody.v1"
)

//用户组 将ws会话信息与用户信息匹配
var session__SessionToUserInfo = make(map[*melody.Session]UserInfo)
var session__SessionToGroup = make(map[*melody.Session][]string)
var session__GroupToSession = make(map[string][]*melody.Session)

// gin类 http 服务器类
var r *gin.Engine

// ws 服务器的
var m *melody.Melody

func init_service() {
	r = gin.Default()
	r.Use(Cors())
	m = melody.New()
	m.Upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	init_service_bind()
	init_service_bind1()
}

var authPath = []string{
	"/api/",
	"/admin/",
}

func init_service_bind() {

	r.POST("/empty", http_empty)
	r.POST("/login", http_login)
	r.POST("/register", http_register)
	r.GET("/ws", func(c *gin.Context) {

		m.HandleRequest(c.Writer, c.Request)
		c.ClientIP()
	})
	r.Use(func(c *gin.Context) {

		for _, url := range authPath {
			if strings.Contains(c.Request.URL.Path, url) {
				var userInfo UserInfo
				token := c.GetHeader("token")
				if token == "" {
					c.JSON(http.StatusInternalServerError, gin.H{
						"state":   false,
						"message": "token 不存在",
					})
					c.Abort()
					return
				}
				if db.First(&userInfo, &UserInfo{ToKen: token}).RecordNotFound() {
					c.JSON(http.StatusInternalServerError, gin.H{
						"state":      false,
						"tokenError": true,
						"message":    "token 不正确",
					})
					c.Abort()
					return
				}
				//TODO token 有效 将用户信息写入会话
				c.Set("userInfo", userInfo)
			}
		}

		c.Next()
	})
	r.Use(func(c *gin.Context) {
		if strings.Contains(c.Request.URL.Path, "/admin/") {
			userInfo := c.Keys["userInfo"].(UserInfo)
			if userInfo.Tags != "admin" {
				c.String(http.StatusNotFound, "")
				c.Abort()
				return
			}
		}
		c.Next()
	})
	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		path = strings.ReplaceAll(path, "/", "\\")
		_log(c.Request.URL.Path)
		c.File(".\\webroot\\" + path)
	})
	m.HandleConnect(func(s *melody.Session) {
		token := s.Request.URL.Query()["toKen"][0]
		var userInfo UserInfo
		if db.First(&userInfo, &UserInfo{ToKen: token}).RecordNotFound() {
			json, _ := json.Marshal(Chat{Data: "Token错误请退出登录.", Model: gorm.Model{CreatedAt: time.Now()}, UserName: "系统消息"})
			s.Write(json)
			s.Write([]byte("close"))
			s.Close()
			return
		}
		lock.Lock()
		session__SessionToUserInfo[s] = userInfo
		lock.Unlock()
		HandleConnect(s, userInfo)
	})
	m.HandleDisconnect(func(s *melody.Session) {
		lock.Lock()
		delete(session__SessionToUserInfo, s)
		list := session__SessionToGroup[s]
		list_ := make([]string, len(list))
		copy(list_, list)
		for _, groupName := range list_ {
			sessionDelGroup(s, groupName, true)
		}
		lock.Unlock()
	})
	m.HandleMessage(func(s *melody.Session, b []byte) {
		userInfo := sessionGetUserInfo(s)
		HandleMessage(s, b, string(b), userInfo)
	})

}
func sessionGetS(name string) []*melody.Session {
	var old []*melody.Session
	lock.Lock()
	old = session__GroupToSession[name]
	lock.Unlock()
	return old
}

func sessionGetUserInfo(s *melody.Session) UserInfo {
	var userInfo UserInfo
	lock.Lock()
	userInfo = session__SessionToUserInfo[s]
	lock.Unlock()
	return userInfo
}

func sessionAddGroup(s *melody.Session, name string) {
	lock.Lock()
	{
		if len(session__GroupToSession[name]) == 0 {
			session__GroupToSession[name] = make([]*melody.Session, 0)
		}
		list := session__GroupToSession[name]
		session__GroupToSession[name] = append(list, s)
	}
	{
		if len(session__SessionToGroup[s]) == 0 {
			session__SessionToGroup[s] = make([]string, 0)
		}
		list := session__SessionToGroup[s]
		session__SessionToGroup[s] = append(list, name)
	}
	lock.Unlock()
}

func broadcastGroup(name string, msg []byte) {
	var list []*melody.Session
	lock.Lock()
	list = session__GroupToSession[name]
	lock.Unlock()
	go func() {
		for _, v := range list {
			v.Write(msg)
		}
	}()
}
func broadcastGroupJson(name string, v interface{}) {
	json, _ := json.Marshal(v)
	broadcastGroup(name, json)
}

func sessionDelGroup(s *melody.Session, name string, unLock bool) {
	if !unLock {
		lock.Lock()
	}
	{
		old := session__GroupToSession[name]
		newData := old
		for k, v := range session__GroupToSession[name] {
			if v == s {
				newData = append(old[:k], old[k+1:]...)
				break
			}
		}
		if len(newData) > 0 {
			session__GroupToSession[name] = newData
		} else {
			delete(session__GroupToSession, name)
		}
	}
	{
		old := session__SessionToGroup[s]
		newData := old
		for k, v := range session__SessionToGroup[s] {
			if v == name {
				newData = append(old[:k], old[k+1:]...)
				break
			}
		}
		if len(newData) > 0 {
			session__SessionToGroup[s] = newData
		} else {
			delete(session__SessionToGroup, s)
		}
	}
	if !unLock {
		lock.Unlock()
	}
}

type SessionUser struct {
	UserInfo UserInfo
	Groups   []string
}
