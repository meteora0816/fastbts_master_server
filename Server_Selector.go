package main

import (
	"fmt"
	"net/http"
	"strconv"
)
import "github.com/gin-gonic/gin"

var servers = []string{"101.35.92.214", "101.35.86.228", "101.35.9.228", "101.35.9.114", // Shanghai
	"49.232.210.247", "152.136.120.165", "152.136.124.173", "49.232.128.240"} // Beijing

//var servers = []string{"101.35.92.213", "101.35.86.227", "101.35.9.227", "101.35.9.113", // Shanghai
//	"49.232.210.246", "152.136.120.164", "152.136.124.172", "49.232.128.240"} // Beijing
var r *gin.Engine

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

var CISSleep = 200
var DownloadSizeSleep = 50
var TimeWindow = 2000
var TestTimeout = 8000
var MaxTrafficUse4g = 50
var MaxTrafficUse5g = 300
var MaxTrafficUseWifi = 1000
var MaxTrafficUseOthers = 1000
var KSimilar = 5
var Threshold = 0.95

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
	r.POST("/speedtest/info", func(c *gin.Context) {
		type Req struct {
			NetworkType        string   `json:"network_type"`
			ServersSortedByRTT []string `json:"servers_sorted_by_rtt"`
		}
		var req Req
		err := c.BindJSON(&req)
		fmt.Println(req.NetworkType)
		if err != nil {
			fmt.Println("err")
			return
		}
		type Res struct {
			ServerNum         int      `json:"server_num"`
			IpList            []string `json:"ip_list"`
			TestTimeout       int      `json:"test_timeout"`
			DownloadSizeSleep int      `json:"download_size_sleep"`
			CISSleep          int      `json:"cis_sleep"`
			TimeWindow        int      `json:"time_window"`
			KSimilar          int      `json:"k_similar"`
			MaxTrafficUse     int      `json:"max_traffic_use"`
			Threshold         float64  `json:"threshold"`
		}
		var res Res
		res.CISSleep = CISSleep
		res.DownloadSizeSleep = DownloadSizeSleep
		res.TimeWindow = TimeWindow
		res.TestTimeout = TestTimeout
		res.MaxTrafficUse = MaxTrafficUseOthers
		res.KSimilar = KSimilar
		res.Threshold = Threshold
		num := 8
		if req.NetworkType == "4G" || req.NetworkType == "3G" {
			num = 4
			res.MaxTrafficUse = MaxTrafficUse4g
		} else if req.NetworkType == "WiFi" {
			num = 8
			res.MaxTrafficUse = MaxTrafficUseWifi
		} else if req.NetworkType == "5G" {
			num = 8
			res.MaxTrafficUse = MaxTrafficUse5g
		} else {
			num = 8
			res.MaxTrafficUse = MaxTrafficUseOthers
		}
		res.ServerNum = min(num, len(req.ServersSortedByRTT))
		//res.IpList = servers[:res.ServerNum]
		//servers = append(servers[res.ServerNum:], servers[:res.ServerNum]...)
		//fmt.Println(servers)
		res.IpList = req.ServersSortedByRTT[:res.ServerNum]
		c.JSON(http.StatusOK, res)
	})
	r.POST("/parameter/MaxTrafficUse4g/:num", func(c *gin.Context) {
		limit := string(c.Param("num"))
		fmt.Println(limit)
		if intLimit, err := strconv.Atoi(limit); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"err":       err.Error(),
				"parameter": limit,
			})
		} else {
			MaxTrafficUse4g = intLimit
			c.JSON(http.StatusOK, gin.H{
				"parameter": MaxTrafficUse4g,
			})
		}
	})
	r.POST("/parameter/MaxTrafficUse5g/:num", func(c *gin.Context) {
		limit := string(c.Param("num"))
		fmt.Println(limit)
		if intLimit, err := strconv.Atoi(limit); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"err":       err.Error(),
				"parameter": limit,
			})
		} else {
			MaxTrafficUse5g = intLimit
			c.JSON(http.StatusOK, gin.H{
				"parameter": MaxTrafficUse5g,
			})
		}
	})
	r.POST("/parameter/MaxTrafficUseWifi/:num", func(c *gin.Context) {
		limit := string(c.Param("num"))
		fmt.Println(limit)
		if intLimit, err := strconv.Atoi(limit); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"err":       err.Error(),
				"parameter": limit,
			})
		} else {
			MaxTrafficUseWifi = intLimit
			c.JSON(http.StatusOK, gin.H{
				"parameter": MaxTrafficUseWifi,
			})
		}
	})
	r.POST("/parameter/MaxTrafficUseOthers/:num", func(c *gin.Context) {
		limit := string(c.Param("num"))
		fmt.Println(limit)
		if intLimit, err := strconv.Atoi(limit); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"err":       err.Error(),
				"parameter": limit,
			})
		} else {
			MaxTrafficUseOthers = intLimit
			c.JSON(http.StatusOK, gin.H{
				"parameter": MaxTrafficUseOthers,
			})
		}
	})
	r.POST("/parameter/TestTimeout/:num", func(c *gin.Context) {
		limit := string(c.Param("num"))
		fmt.Println(limit)
		if intLimit, err := strconv.Atoi(limit); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"err":       err.Error(),
				"parameter": limit,
			})
		} else {
			TestTimeout = intLimit
			c.JSON(http.StatusOK, gin.H{
				"parameter": TestTimeout,
			})
		}
	})
	r.POST("/parameter/KSimilar/:num", func(c *gin.Context) {
		limit := string(c.Param("num"))
		fmt.Println(limit)
		if intLimit, err := strconv.Atoi(limit); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"err":       err.Error(),
				"parameter": limit,
			})
		} else {
			KSimilar = intLimit
			c.JSON(http.StatusOK, gin.H{
				"parameter": KSimilar,
			})
		}
	})
	if err := r.Run(); err != nil {
		fmt.Println(err)
	}
}
