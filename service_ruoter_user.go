package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func init_service_bind_user() {
	r.GET("/api/search_user", func(c *gin.Context) {
		target := c.Query("target")
		if target == "" {
			c.JSON(http.StatusOK, gin.H{
				"state":   false,
				"message": "请输入关键字",
			})
			return
		}
		var users []UserInfo
		db.Where("user_name LIKE ?", "%"+target+"%").Find(&users)
		var usernames []string
		for _, user := range users {
			usernames = append(usernames, user.UserName)
		}
		c.JSON(http.StatusOK, gin.H{
			"state":   true,
			"list":    usernames,
			"message": "查询到" + strconv.Itoa(len(usernames)) + "条好友信息.",
		})
	})

	r.GET("/api/get_session", func(c *gin.Context) {
		userInfo := c.Keys["userInfo"].(UserInfo)
		var friends []Friends
		db.Order("last_chat_at desc").Find(&friends, &Friends{
			UserName: userInfo.UserName,
		})
		var group []GroupRelation
		db.Order("last_chat_at desc").Find(&group, &GroupRelation{
			UserName: userInfo.UserName,
		})
		c.JSON(http.StatusOK, gin.H{
			"state": true,
			"data": gin.H{
				"friends": friends,
				"group":   group,
			},
			"message": "",
		})
	})

}
