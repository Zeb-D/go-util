package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"
)

func main() {
	page := getUrls(1)
	fmt.Println(page)

	if err := indexArticleMap(); err != nil {
		fmt.Println("indexArticleMap fail", err)
		return
	}

	fmt.Println("len(articleMap) == total :", len(articleMap) == total)

	//排序
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for true {
			t := time.Now()
			for id, art := range articleMap {
				r := rand.Intn(10)
				if art.ViewCount > 10000 && r >= 3 {
					time.Sleep(time.Duration(rand.Intn(3)+6) * time.Second)
					continue
				} else if art.ViewCount > 8000 && r >= 6 {
					time.Sleep(time.Duration(rand.Intn(3)+3) * time.Second)
					continue
				} else if art.ViewCount > 1000 && r >= 9 {
					time.Sleep(time.Duration(rand.Intn(2)+1) * time.Second)
					continue
				}

				bd := fetch(art.Url, false)
				if time.Now().Sub(t).Seconds() <= 9 {
					time.Sleep(time.Duration(rand.Intn(5)+4) * time.Second)
				}
				if len(bd) == 0 {
					fmt.Printf("article fetch error id:%d\n", id)
				} else {
					fmt.Printf("article fetch time:%s,url:%s \n", time.Now().Sub(t).String(), art.Url)
				}
				t = time.Now()
			}
			time.Sleep(time.Duration(rand.Intn(10)+10) * time.Second)
			indexArticleMap()
		}
		wg.Done()
	}()
	wg.Wait()
}

var total = 170
var articleMap = map[int]Article{}

const pageSize = 20
const csdn_format = `https://blog.csdn.net/community/home-api/v1/get-business-list?page=%d&size=%d&businessType=blog&orderby=&noMore=false&year=&month=&username=u014229282`

var UA = []string{"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.1 (KHTML, like Gecko) Chrome/22.0.1207.1 Safari/537.1",
	"Mozilla/5.0 (X11; CrOS i686 2268.111.0) AppleWebKit/536.11 (KHTML, like Gecko) Chrome/20.0.1132.57 Safari/536.11",
	"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/536.6 (KHTML, like Gecko) Chrome/20.0.1092.0 Safari/536.6",
	"Mozilla/5.0 (Windows NT 6.2) AppleWebKit/536.6 (KHTML, like Gecko) Chrome/20.0.1090.0 Safari/536.6",
	"Mozilla/5.0 (Windows NT 6.2; WOW64) AppleWebKit/537.1 (KHTML, like Gecko) Chrome/19.77.34.5 Safari/537.1"}

func indexArticleMap() error {
	size := total/pageSize + 1
	for i := 1; i <= size; i++ {
		as := getUrls(i)
		if len(as) > 0 {
			for _, value := range as {
				if value != nil {
					articleMap[value.ArticleId] = *value
				}
			}
		}
	}
	return nil
}

//csdn爬虫
//https://blog.csdn.net/community/home-api/v1/get-business-list?page=1&size=20&businessType=blog&orderby=&noMore=false&year=&month=&username=u014229282
//一分钟内同一个IP访问算一次的计算规则
func getUrls(pageNo int) []*Article {
	url := fmt.Sprintf(csdn_format, pageNo, pageSize)
	body := fetch(url, true)
	fmt.Println(body)
	if len(body) < 1 {
		return nil
	}
	page := &CsdnPage{}
	json.Unmarshal([]byte(body), page)
	if page.Code != 200 {
		fmt.Println("page.Code", page.Code)
		return nil
	}
	total = page.Data.Total

	return page.Data.List
}

//"code":200,"message":"success","data"
type CsdnPage struct {
	Code  int    `json:"code,omitempty"`
	Value string `json:"message,omitempty"`
	Data  struct {
		Total int        `json:"total,omitempty"`
		List  []*Article `json:"list,omitempty"`
	} `json:"data,omitempty"`
}

type Article struct {
	ArticleId int    `json:"articleId,omitempty"`
	Title     string `json:"title,omitempty"`
	Url       string `json:"url,omitempty"`
	ViewCount int    `json:"viewCount,omitempty"`
}

var header = http.Header{}

func init() {
	header.Add("sec-ch-ua-mobile", "?0")
	header.Add("sec-ch-ua-platform", "macOS")
	header.Add("sec-fetch-dest", "document")
	header.Add("sec-fetch-mode", "navigate")
	header.Add("sec-fetch-site", "same-origin")
	header.Add("sec-fetch-user", "?1")
	header.Add("sec-ch-ua", `" Not A;Brand";v="99", "Chromium";v="98", "Google Chrome";v="98"`)
	articleMap[1] = Article{
		Url: "https://juejin.cn/post/7056232962852913160",
	}
	articleMap[2] = Article{
		Url: "https://juejin.cn/post/7079594713174507551",
	}
	articleMap[11] = Article{
		Url: "https://www.jianshu.com/p/c5c403706984",
	}
	articleMap[12] = Article{
		Url: "https://www.jianshu.com/p/ed47fad9c9fe",
	}
	articleMap[13] = Article{
		Url: "https://www.jianshu.com/p/6a919f9b9fab",
	}
	articleMap[14] = Article{
		Url: "https://www.jianshu.com/p/171695345bc1",
	}
	articleMap[21] = Article{
		Url: "https://www.cnblogs.com/zeb-d/articles/16028602.html",
	}
}

func fetch(url string, print bool) string {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		fmt.Println("NewRequest err", err)
		return ""
	}
	req.Header = header

	req.Header.Set("Referer", url)
	rand := rand.New(rand.NewSource(time.Now().Unix()))
	req.Header.Set("User-Agent", UA[rand.Intn(len(UA))])
	ip := genIpAddr()
	req.Header.Set("X-Forwarded-For", ip)
	req.Header.Set("X-Original-Forwarded-For", ip)
	req.Header.Set("X-Real-Ip", ip)
	req.Header.Set("Proxy-Client-IP", ip)
	req.Header.Set("WL-Proxy-Client-IP", ip)
	req.Header.Set("HTTP_CLIENT_IP", ip)
	req.Header.Set("HTTP_X_FORWARDED_FOR", ip)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Http get err:", err)
		return ""
	}
	if resp.StatusCode != 200 {
		fmt.Println("Http status code:", resp.StatusCode)
		return ""
	}
	defer resp.Body.Close()
	if !print {
		return "-1"
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Read error", err)
		return ""
	}
	return string(body)
}

func genIpAddr() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	i := r.Intn(800) + r.Intn(680)
	s := chinaIpArr[i%len(chinaIpArr)]
	ss := strings.Split(s, ".")
	var buffer bytes.Buffer
	for idx, sss := range ss {
		if sss == "0" {
			sss = fmt.Sprint(1 + r.Intn(253))
		}
		buffer.WriteString(sss)
		if idx < 3 {
			buffer.WriteString(".")
		}
	}
	return buffer.String()
}

var chinaIpArr = strings.Split(chinaIp, ",")
var chinaIp = `58.14.0.0,58.16.0.0,58.24.0.0,58.30.0.0,58.32.0.0,58.66.0.0,58.68.128.0,58.82.0.0,58.87.64.0,58.99.128.0,58.100.0.0,58.116.0.0,58.128.0.0,58.144.0.0,58.154.0.0,58.192.0.0,58.240.0.0,59.32.0.0,59.64.0.0,59.80.0.0,59.107.0.0,59.108.0.0,59.151.0.0,59.155.0.0,59.172.0.0,59.191.0.0,59.191.240.0,59.192.0.0,60.0.0.0,60.55.0.0,60.63.0.0,60.160.0.0,60.194.0.0,60.200.0.0,60.208.0.0,60.232.0.0,60.235.0.0,60.245.128.0,60.247.0.0,60.252.0.0,60.253.128.0,60.255.0.0,61.4.80.0,61.4.176.0,61.8.160.0,61.28.0.0,61.29.128.0,61.45.128.0,61.47.128.0,61.48.0.0,61.87.192.0,61.128.0.0,61.232.0.0,61.236.0.0,61.240.0.0,114.28.0.0,114.54.0.0,114.60.0.0,114.64.0.0,114.68.0.0,114.80.0.0,116.1.0.0,116.2.0.0,116.4.0.0,116.8.0.0,116.13.0.0,116.16.0.0,116.52.0.0,116.56.0.0,116.58.128.0,116.58.208.0,116.60.0.0,116.66.0.0,116.69.0.0,116.70.0.0,116.76.0.0,116.89.144.0,116.90.184.0,116.95.0.0,116.112.0.0,116.116.0.0,116.128.0.0,116.192.0.0,116.193.16.0,116.193.32.0,116.194.0.0,116.196.0.0,116.198.0.0,116.199.0.0,116.199.128.0,116.204.0.0,116.207.0.0,116.208.0.0,116.212.160.0,116.213.64.0,116.213.128.0,116.214.32.0,116.214.64.0,116.214.128.0,116.215.0.0,116.216.0.0,116.224.0.0,116.242.0.0,116.244.0.0,116.248.0.0,116.252.0.0,116.254.128.0,116.255.128.0,117.8.0.0,117.21.0.0,117.22.0.0,117.24.0.0,117.32.0.0,117.40.0.0,117.44.0.0,117.48.0.0,117.53.48.0,117.53.176.0,117.57.0.0,117.58.0.0,117.59.0.0,117.60.0.0,117.64.0.0,117.72.0.0,117.74.64.0,117.74.128.0,117.75.0.0,117.76.0.0,117.80.0.0,117.100.0.0,117.103.16.0,117.103.128.0,117.106.0.0,117.112.0.0,117.120.64.0,117.120.128.0,117.121.0.0,117.121.128.0,117.121.192.0,117.122.128.0,117.124.0.0,117.128.0.0,118.24.0.0,118.64.0.0,118.66.0.0,118.67.112.0,118.72.0.0,118.80.0.0,118.84.0.0,118.88.32.0,118.88.64.0,118.88.128.0,118.89.0.0,118.91.240.0,118.102.16.0,118.112.0.0,118.120.0.0,118.124.0.0,118.126.0.0,118.132.0.0,118.144.0.0,118.178.0.0,118.180.0.0,118.184.0.0,118.192.0.0,118.212.0.0,118.224.0.0,118.228.0.0,118.230.0.0,118.239.0.0,118.242.0.0,118.244.0.0,118.248.0.0,119.0.0.0,119.2.0.0,119.2.128.0,119.3.0.0,119.4.0.0,119.8.0.0,119.10.0.0,119.15.136.0,119.16.0.0,119.18.192.0,119.18.208.0,119.18.224.0,119.19.0.0,119.20.0.0,119.27.64.0,119.27.160.0,119.27.192.0,119.28.0.0,119.30.48.0,119.31.192.0,119.32.0.0,119.40.0.0,119.40.64.0,119.40.128.0,119.41.0.0,119.42.0.0,119.42.136.0,119.42.224.0,119.44.0.0,119.48.0.0,119.57.0.0,119.58.0.0,119.59.128.0,119.60.0.0,119.62.0.0,119.63.32.0,119.75.208.0,119.78.0.0,119.80.0.0,119.84.0.0,119.88.0.0,119.96.0.0,119.108.0.0,119.112.0.0,119.128.0.0,119.144.0.0,119.148.160.0,119.161.128.0,119.162.0.0,119.164.0.0,119.176.0.0,119.232.0.0,119.235.128.0,119.248.0.0,119.253.0.0,119.254.0.0,120.0.0.0,120.24.0.0,120.30.0.0,120.32.0.0,120.48.0.0,120.52.0.0,120.64.0.0,120.72.32.0,120.72.128.0,120.76.0.0,120.80.0.0,120.90.0.0,120.92.0.0,120.94.0.0,120.128.0.0,120.136.128.0,120.137.0.0,120.192.0.0,121.0.16.0,121.4.0.0,121.8.0.0,121.16.0.0,121.32.0.0,121.40.0.0,121.46.0.0,121.48.0.0,121.51.0.0,121.52.160.0,121.52.208.0,121.52.224.0,121.55.0.0,121.56.0.0,121.58.0.0,121.58.144.0,121.59.0.0,121.60.0.0,121.68.0.0,121.76.0.0,121.79.128.0,121.89.0.0,121.100.128.0,121.101.208.0,121.192.0.0,121.201.0.0,121.204.0.0,121.224.0.0,121.248.0.0,121.255.0.0,122.0.64.0,122.0.128.0,122.4.0.0,122.8.0.0,122.48.0.0,122.49.0.0,122.51.0.0,122.64.0.0,122.96.0.0,122.102.0.0,122.102.64.0,122.112.0.0,122.119.0.0,122.136.0.0,122.144.128.0,122.152.192.0,122.156.0.0,122.192.0.0,122.198.0.0,122.200.64.0,122.204.0.0,122.224.0.0,122.240.0.0,122.248.48.0,123.0.128.0,123.4.0.0,123.8.0.0,123.49.128.0,123.52.0.0,123.56.0.0,123.64.0.0,123.96.0.0,123.98.0.0,123.99.128.0,123.100.0.0,123.101.0.0,123.103.0.0,123.108.128.0,123.108.208.0,123.112.0.0,123.128.0.0,123.136.80.0,123.137.0.0,123.138.0.0,123.144.0.0,123.160.0.0,123.176.80.0,123.177.0.0,123.178.0.0,123.180.0.0,123.184.0.0,123.196.0.0,123.199.128.0,123.206.0.0,123.232.0.0,123.242.0.0,123.244.0.0,123.249.0.0,123.253.0.0,124.6.64.0,124.14.0.0,124.16.0.0,124.20.0.0,124.28.192.0,124.29.0.0,124.31.0.0,124.40.112.0,124.40.128.0,124.42.0.0,124.47.0.0,124.64.0.0,124.66.0.0,124.67.0.0,124.68.0.0,124.72.0.0,124.88.0.0,124.108.8.0,124.108.40.0,124.112.0.0,124.126.0.0,124.128.0.0,124.147.128.0,124.156.0.0,124.160.0.0,124.172.0.0,124.192.0.0,124.196.0.0,124.200.0.0,124.220.0.0,124.224.0.0,124.240.0.0,124.240.128.0,124.242.0.0,124.243.192.0,124.248.0.0,124.249.0.0,124.250.0.0,124.254.0.0,125.31.192.0,125.32.0.0,125.58.128.0,125.61.128.0,125.62.0.0,125.64.0.0,125.96.0.0,125.98.0.0,125.104.0.0,125.112.0.0,125.169.0.0,125.171.0.0,125.208.0.0,125.210.0.0,125.213.0.0,125.214.96.0,125.215.0.0,125.216.0.0,125.254.128.0,134.196.0.0,159.226.0.0,161.207.0.0,162.105.0.0,166.111.0.0,167.139.0.0,168.160.0.0,169.211.1.0,192.83.122.0,192.83.169.0,192.124.154.0,192.188.170.0,198.17.7.0,202.0.110.0,202.0.176.0,202.4.128.0,202.4.252.0,202.8.128.0,202.10.64.0,202.14.88.0,202.14.235.0,202.14.236.0,202.14.238.0,202.20.120.0,202.22.248.0,202.38.0.0,202.38.64.0,202.38.128.0,202.38.136.0,202.38.138.0,202.38.140.0,202.38.146.0,202.38.149.0,202.38.150.0,202.38.152.0,202.38.156.0,202.38.158.0,202.38.160.0,202.38.164.0,202.38.168.0,202.38.176.0,202.38.184.0,202.38.192.0,202.41.152.0,202.41.240.0,202.43.144.0,202.46.32.0,202.46.224.0,202.60.112.0,202.63.248.0,202.69.4.0,202.69.16.0,202.70.0.0,202.74.8.0,202.75.208.0,202.85.208.0,202.90.0.0,202.90.224.0,202.90.252.0,202.91.0.0,202.91.128.0,202.91.176.0,202.91.224.0,202.92.0.0,202.92.252.0,202.93.0.0,202.93.252.0,202.95.0.0,202.95.252.0,202.96.0.0,202.112.0.0,202.120.0.0,202.122.0.0,202.122.32.0,202.122.64.0,202.122.112.0,202.122.128.0,202.123.96.0,202.124.24.0,202.125.176.0,202.127.0.0,202.127.12.0,202.127.16.0,202.127.40.0,202.127.48.0,202.127.112.0,202.127.128.0,202.127.160.0,202.127.192.0,202.127.208.0,202.127.212.0,202.127.216.0,202.127.224.0,202.130.0.0,202.130.224.0,202.131.16.0,202.131.48.0,202.131.208.0,202.136.48.0,202.136.208.0,202.136.224.0,202.141.160.0,202.142.16.0,202.143.16.0,202.148.96.0,202.149.160.0,202.149.224.0,202.150.16.0,202.152.176.0,202.153.48.0,202.158.160.0,202.160.176.0,202.164.0.0,202.164.25.0,202.165.96.0,202.165.176.0,202.165.208.0,202.168.160.0,202.170.128.0,202.170.216.0,202.173.8.0,202.173.224.0,202.179.240.0,202.180.128.0,202.181.112.0,202.189.80.0,202.192.0.0,203.18.50.0,203.79.0.0,203.80.144.0,203.81.16.0,203.83.56.0,203.86.0.0,203.86.64.0,203.88.32.0,203.88.192.0,203.89.0.0,203.90.0.0,203.90.128.0,203.90.192.0,203.91.32.0,203.91.96.0,203.91.120.0,203.92.0.0,203.92.160.0,203.93.0.0,203.94.0.0,203.95.0.0,203.95.96.0,203.99.16.0,203.99.80.0,203.100.32.0,203.100.80.0,203.100.96.0,203.100.192.0,203.110.160.0,203.118.192.0,203.119.24.0,203.119.32.0,203.128.32.0,203.128.96.0,203.130.32.0,203.132.32.0,203.134.240.0,203.135.96.0,203.135.160.0,203.142.219.0,203.148.0.0,203.152.64.0,203.156.192.0,203.158.16.0,203.161.192.0,203.166.160.0,203.171.224.0,203.174.7.0,203.174.96.0,203.175.128.0,203.175.192.0,203.176.168.0,203.184.80.0,203.187.160.0,203.190.96.0,203.191.16.0,203.191.64.0,203.191.144.0,203.192.0.0,203.196.0.0,203.207.64.0,203.207.128.0,203.208.0.0,203.208.16.0,203.208.32.0,203.209.224.0,203.212.0.0,203.212.80.0,203.222.192.0,203.223.0.0,210.2.0.0,210.5.0.0,210.5.144.0,210.12.0.0,210.14.64.0,210.14.112.0,210.14.128.0,210.15.0.0,210.15.128.0,210.16.128.0,210.21.0.0,210.22.0.0,210.23.32.0,210.25.0.0,210.26.0.0,210.28.0.0,210.32.0.0,210.51.0.0,210.52.0.0,210.56.192.0,210.72.0.0,210.76.0.0,210.78.0.0,210.79.64.0,210.79.224.0,210.82.0.0,210.87.128.0,210.185.192.0,210.192.96.0,211.64.0.0,211.80.0.0,211.96.0.0,211.136.0.0,211.144.0.0,211.160.0.0,218.0.0.0,218.56.0.0,218.64.0.0,218.96.0.0,218.104.0.0,218.108.0.0,218.185.192.0,218.192.0.0,218.240.0.0,218.249.0.0,219.72.0.0,219.82.0.0,219.128.0.0,219.216.0.0,219.224.0.0,219.242.0.0,219.244.0.0,220.101.192.0,220.112.0.0,220.152.128.0,220.154.0.0,220.160.0.0,220.192.0.0,220.231.0.0,220.231.128.0,220.232.64.0,220.234.0.0,220.242.0.0,220.248.0.0,220.252.0.0,221.0.0.0,221.8.0.0,221.12.0.0,221.12.128.0,221.13.0.0,221.14.0.0,221.122.0.0,221.129.0.0,221.130.0.0,221.133.224.0,221.136.0.0,221.172.0.0,221.176.0.0,221.192.0.0,221.196.0.0,221.198.0.0,221.199.0.0,221.199.128.0,221.199.192.0,221.199.224.0,221.200.0.0,221.208.0.0,221.224.0.0,222.16.0.0,222.32.0.0,222.64.0.0,222.125.0.0,222.126.128.0,222.128.0.0,222.160.0.0,222.168.0.0,222.176.0.0,222.192.0.0,222.240.0.0,222.248.0.0`
