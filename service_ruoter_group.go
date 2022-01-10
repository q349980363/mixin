package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func init_service_bind_group() {
	r.GET("/api/search_group", func(c *gin.Context) {
		target := c.Query("target")
		if target == "" {
			c.JSON(http.StatusOK, gin.H{
				"state":   false,
				"message": "请输入关键字",
			})
			return
		}
		var list []Group
		db.Where("user_name LIKE ?", "%"+target+"%").Find(&list)
		var names []string
		for _, item := range list {
			names = append(names, item.Name)
		}
		c.JSON(http.StatusOK, gin.H{
			"state":   true,
			"list":    names,
			"message": "查询到" + strconv.Itoa(len(names)) + "条群信息.",
		})
	})
	r.POST("/api/add_group/:target", func(c *gin.Context) {
		target := c.Param("target")
		userInfo := c.Keys["userInfo"].(UserInfo)
		error := addGroup(userInfo.UserName, target)
		if error != nil {
			c.JSON(http.StatusOK, gin.H{
				"state":   false,
				"message": error.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"state":   true,
			"message": "加群请求已发送.",
		})
	})
	r.GET("/api/create_group", func(c *gin.Context) {
		target := c.Query("target")
		if target == "" {
			c.JSON(http.StatusOK, gin.H{
				"state":   false,
				"message": "请输入关键字",
			})
			return
		}
		userInfo := c.Keys["userInfo"].(UserInfo)

		if !db.First(&Group{}, &Group{
			Name: target,
		}).RecordNotFound() {
			c.JSON(http.StatusOK, gin.H{
				"state":   false,
				"message": "群已存在",
			})
			return
		}

		db.Create(&Group{
			UserName: userInfo.UserName,
			Name:     target,
		})
		db.Create(&GroupRelation{
			UserName: userInfo.UserName,
			Target:   target,
		})
		c.JSON(http.StatusOK, gin.H{
			"state":   true,
			"message": "创建群成功.",
		})
	})
	r.GET("/api/get_groupchat/:target", func(c *gin.Context) {
		target := c.Param("target")
		userInfo := c.Keys["userInfo"].(UserInfo)
		var list []GroupChat
		db.Where(&GroupChat{
			UserName: userInfo.UserName,
			Target:   target,
		}).Find(&list)
		db.Model(&GroupRelation{}).Where(&GroupRelation{
			UserName: userInfo.UserName,
			Target:   target,
		}).Update("Unread", 0)
		c.JSON(http.StatusOK, gin.H{
			"state":   true,
			"list":    list,
			"message": "",
		})
	})

	r.POST("/api/del_group/:target", func(c *gin.Context) {
		target := c.Query("target")
		userInfo := c.Keys["userInfo"].(UserInfo)
		db.Where(&GroupRelation{
			UserName: userInfo.UserName,
			Target:   target,
		}).Delete(GroupRelation{})

		c.JSON(http.StatusOK, gin.H{
			"state":   true,
			"message": ".",
		})
	})
	r.POST("/api/del_groupchat/:target", func(c *gin.Context) {
		target := c.Query("target")
		// userInfo := c.Keys["userInfo"].(UserInfo)
		db.Where(&GroupChat{
			Target: target,
		}).Delete(GroupChat{})
		c.JSON(http.StatusOK, gin.H{
			"state":   true,
			"message": ".",
		})
	})
}
