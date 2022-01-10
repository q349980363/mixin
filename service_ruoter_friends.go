package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func init_service_bind_friends() {

	r.GET("/api/get_friendschat/:target", func(c *gin.Context) {
		target := c.Param("target")
		userInfo := c.Keys["userInfo"].(UserInfo)
		var list []FriendsChat
		db.Where(&FriendsChat{
			UserName: userInfo.UserName,
			Target:   target,
		}).Or(&FriendsChat{
			UserName: target,
			Target:   userInfo.UserName,
		}).Find(&list)
		db.Model(&Friends{}).Where(&Friends{
			UserName: userInfo.UserName,
			Target:   target,
		}).Update("Unread", 0)
		c.JSON(http.StatusOK, gin.H{
			"state":   true,
			"list":    list,
			"message": "",
		})
	})

	r.POST("/api/add_friends/:target", func(c *gin.Context) {
		target := c.Param("target")
		userInfo := c.Keys["userInfo"].(UserInfo)
		error := addFriends(userInfo.UserName, target)
		if error != nil {
			c.JSON(http.StatusOK, gin.H{
				"state":   false,
				"message": error.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"state":   true,
			"message": "添加请求已发送.",
		})
	})
	r.POST("/api/del_friends/:target", func(c *gin.Context) {
		target := c.Query("target")
		userInfo := c.Keys["userInfo"].(UserInfo)
		db.Where(&Friends{
			UserName: userInfo.UserName,
			Target:   target,
		}).Or(&Friends{
			UserName: target,
			Target:   userInfo.UserName,
		}).Delete(Friends{})

		c.JSON(http.StatusOK, gin.H{
			"state":   true,
			"message": ".",
		})
	})
	r.POST("/api/del_friendschat/:target", func(c *gin.Context) {
		target := c.Query("target")
		userInfo := c.Keys["userInfo"].(UserInfo)
		db.Where(&FriendsChat{
			UserName: userInfo.UserName,
			Target:   target,
		}).Or(&FriendsChat{
			UserName: target,
			Target:   userInfo.UserName,
		}).Delete(FriendsChat{})
		c.JSON(http.StatusOK, gin.H{
			"state":   true,
			"message": ".",
		})
	})
}
