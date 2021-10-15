package main

import (
	"fmt"
	"net/http"
)
import "github.com/gin-gonic/gin"

var servers = []string{"101.35.92.214", "101.35.86.228", "101.35.9.228", "101.35.9.114", // Shanghai
	"49.232.210.247", "152.136.120.165", "152.136.124.173", "49.232.128.240"} // Beijing

var r *gin.Engine

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	fmt.Println("start")
	r = gin.Default()
	r.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, `hello, FastBTS!`)
	})
	r.GET("/speedtest/iplist/available", func(c *gin.Context) {
		type Res struct {
			ServerNum int      `json:"server_num"`
			IpList    []string `json:"ip_list"`
			ClientIP  string   `json:"client_ip"`
		}
		var res Res
		res.ServerNum = len(servers)
		res.IpList = servers
		res.ClientIP = c.ClientIP()
		c.JSON(http.StatusOK, res)
	})
	r.GET("/speedtest/info", func(c *gin.Context) {
		type Req struct {
			NetworkType        string   `json:"network_type"`
			ServersSortedByRTT []string `json:"servers_sorted_by_rtt"`
		}
		var req Req
		err := c.BindJSON(&req)
		if err != nil {
			fmt.Println("err")
			return
		}
		type Res struct {
			ServerNum int      `json:"server_num"`
			IpList    []string `json:"ip_list"`
		}
		var res Res
		num := 8
		if req.NetworkType == "4G" {
			num = 4
		} else if req.NetworkType == "WiFi" {
			num = 8
		} else if req.NetworkType == "5G" {
			num = 8
		}
		res.ServerNum = min(num, len(req.ServersSortedByRTT))
		//res.IpList = servers[:res.ServerNum]
		//servers = append(servers[res.ServerNum:], servers[:res.ServerNum]...)
		//fmt.Println(servers)
		res.IpList = req.ServersSortedByRTT[:res.ServerNum]
		c.JSON(http.StatusOK, res)
	})
	if err := r.Run(); err != nil {
		fmt.Println(err)
	}
}
