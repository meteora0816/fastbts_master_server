package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"sync"
)
import "github.com/gin-gonic/gin"

var servers = []string{"101.35.92.214", "101.35.86.228", "101.35.9.228", "101.35.9.114", "110.42.169.86", "121.5.26.137", "1.116.117.183", "1.15.30.244", // Shanghai
					"49.232.210.247", "152.136.120.165", "152.136.124.173", "49.232.128.240", "81.70.193.140", "49.232.129.114", "62.234.117.45", "81.70.55.189"} // Beijing
const maxBandwidth float64 = 200 // bandwidth limit for each server

type SafeMap struct {
	bandwidthUsed map[string]float64
	mux           sync.RWMutex
}

var r *gin.Engine

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type Config struct {
	CISSleep            int     `yaml:"cis_sleep"`
	DownloadSizeSleep   int     `yaml:"download_size_sleep"`
	TimeWindow          int     `yaml:"time_window"`
	TestTimeout         int     `yaml:"test_timeout"`
	GetInfoInterval     int     `yaml:"get_info_interval"`
	MaxTrafficUse4g     int     `yaml:"max_traffic_use_4_g"`
	MaxTrafficUse5g     int     `yaml:"max_traffic_use_5_g"`
	MaxTrafficUseWifi   int     `yaml:"max_traffic_use_wifi"`
	MaxTrafficUseOthers int     `yaml:"max_traffic_use_others"`
	KSimilar            int     `yaml:"k_similar"`
	Threshold           float64 `yaml:"threshold"`
}

var GlobalConfig Config

func init() {
	config, err := ioutil.ReadFile("./config.yaml")
	if err != nil {
		fmt.Print(err)
	}
	err = yaml.Unmarshal(config, &GlobalConfig)
	if err != nil {
		fmt.Print(err)
	} else {
		fmt.Println(GlobalConfig)
	}
}

func getBandwidthUsed(ip string) float64 {
	resp, err := http.Get("http://" + ip + ":8000/bandwidth")
	if err != nil {
		fmt.Println(err)
		return 10000
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(body))
	bandwidthUsed, _ := strconv.ParseFloat(string(body), 64)
	//fmt.Println(bandwidthUsed)
	return bandwidthUsed * 8 / 1024 / 1024
}

type bdu struct {
	ip string
	bd float64 // bandwidthUsed
}

func SS(eBandwidth float64, ServersSortedByRTT []string) (int, []string) {
	var bandwidthUsed []bdu
	wg := sync.WaitGroup{}
	wg.Add(len(servers))
	buCh := make(chan bdu, len(servers))
	for _, ip := range servers {
		ip := ip
		go func() {
			bandwidthUsed := bdu{ip: ip, bd: getBandwidthUsed(ip)}
			buCh <- bandwidthUsed
			wg.Done()
		}()
	}
	wg.Wait()
	close(buCh)
	for bu := range buCh {
		bandwidthUsed = append(bandwidthUsed, bu)
		//fmt.Println(bu.ip, bu.bd)
	}
	//fmt.Println(bandwidthUsed)
	sort.Slice(bandwidthUsed, func(i, j int) bool {
		return bandwidthUsed[i].bd < bandwidthUsed[j].bd
	})
	fmt.Println(bandwidthUsed)
	num := 0
	var ipList []string
	for _, bu := range bandwidthUsed {
		rest := maxBandwidth - bu.bd
		if rest <= 0 {
			continue
		} else {
			num++
			ipList = append(ipList, bu.ip)
			eBandwidth -= rest
			if eBandwidth <= 0 {
				break
			}
		}
	}
	if eBandwidth > 0 {
		return -1, nil
	} else {
		return num, ipList
	}
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
	r.POST("/speedtest/info", func(c *gin.Context) {
		type Req struct {
			NetworkType        string   `json:"network_type"`
			ServersSortedByRTT []string `json:"servers_sorted_by_rtt"`
		}
		var req Req
		err := c.BindJSON(&req)
		//fmt.Println(req.NetworkType)
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
			GetInfoInterval   int      `json:"get_info_interval"`
		}
		var res Res
		res.CISSleep = GlobalConfig.CISSleep
		res.DownloadSizeSleep = GlobalConfig.DownloadSizeSleep
		res.TimeWindow = GlobalConfig.TimeWindow
		res.TestTimeout = GlobalConfig.TestTimeout
		res.MaxTrafficUse = GlobalConfig.MaxTrafficUseOthers
		res.KSimilar = GlobalConfig.KSimilar
		res.Threshold = GlobalConfig.Threshold
		res.GetInfoInterval = GlobalConfig.GetInfoInterval
		var eBandwidth float64
		if req.NetworkType == "LTE" || req.NetworkType == "3G" || req.NetworkType == "2G" {
			eBandwidth = 400
			res.MaxTrafficUse = GlobalConfig.MaxTrafficUse4g
		} else if req.NetworkType == "WIFI" {
			eBandwidth = 1500
			res.MaxTrafficUse = GlobalConfig.MaxTrafficUseWifi
		} else if req.NetworkType == "5G" {
			eBandwidth = 1000
			res.MaxTrafficUse = GlobalConfig.MaxTrafficUse5g
		} else {
			eBandwidth = 500
			res.MaxTrafficUse = GlobalConfig.MaxTrafficUseOthers
		}
		res.ServerNum, res.IpList = SS(eBandwidth, req.ServersSortedByRTT)
		//num := 4
		//if req.NetworkType == "LTE" || req.NetworkType == "3G" || req.NetworkType == "2G" {
		//	num = 2
		//	res.MaxTrafficUse = GlobalConfig.MaxTrafficUse4g
		//} else if req.NetworkType == "WIFI" {
		//	num = 6
		//	res.MaxTrafficUse = GlobalConfig.MaxTrafficUseWifi
		//} else if req.NetworkType == "5G" {
		//	num = 6
		//	res.MaxTrafficUse = GlobalConfig.MaxTrafficUse5g
		//} else {
		//	num = 6
		//	res.MaxTrafficUse = GlobalConfig.MaxTrafficUseOthers
		//}
		//res.ServerNum = min(num, len(req.ServersSortedByRTT))
		////res.IpList = servers[:res.ServerNum]
		////servers = append(servers[res.ServerNum:], servers[:res.ServerNum]...)
		////fmt.Println(req.ServersSortedByRTT)
		////fmt.Println(lastUse)
		//lastUse.mux.Lock()
		//defer lastUse.mux.Unlock()
		//var doNotUse []string
		//tp := 0
		//for _, ip := range req.ServersSortedByRTT {
		//	//fmt.Println(ip, time.Since(lastUse.lastUseTime[ip]).Milliseconds(), int64(GlobalConfig.TestTimeout))
		//	if tp < res.ServerNum && time.Since(lastUse.lastUseTime[ip]).Milliseconds() >= int64(GlobalConfig.TestTimeout) {
		//		tp++
		//		lastUse.lastUseTime[ip] = time.Now()
		//		res.IpList = append(res.IpList, ip)
		//	} else {
		//		doNotUse = append(doNotUse, ip)
		//	}
		//}
		//for _, ip := range doNotUse {
		//	if tp < res.ServerNum {
		//		tp++
		//		lastUse.lastUseTime[ip] = time.Now()
		//		res.IpList = append(res.IpList, ip)
		//	}
		//}
		//fmt.Println(res.IpList)
		c.JSON(http.StatusOK, res)
	})
	if err := r.Run(); err != nil {
		fmt.Println(err)
	}
}
