package popin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"popindata/function"
	"strconv"
	"strings"
)

//popinJSONObj 爬取的json数据结构
type popinJSONObj struct {
	TotalMoney string `json:"total_charge"`
}

//POP POPin的相关信息
type POP struct {
	Account      string   `json:"popin_account"`
	Password     string   `json:"popin_password"`
	CampaignList []string `json:"popin_CampaignList"`
	Cookie       string   `json:"popin_cookie"`
}

//Popin 实例化一个POP对象
var Popin POP

func init() {
	var file *os.File
	var err error
	var configByte []byte
	if file, err = os.Open("config.json"); err != nil {
		log.Fatal(err)
	}
	if configByte, err = ioutil.ReadAll(file); err != nil {
		log.Fatal(err)
	}
	if err = json.Unmarshal(configByte, &Popin); err != nil {
		log.Fatal(err)
	}
}

//GetAllMoney 获取popin指定campaign的总消耗,最常用的
func (pop *POP) GetAllMoney() (money int) {
	for _, v := range pop.CampaignList {
		j := getPopinMoney(pop.Cookie, pop.Account, v)
		money += j
	}
	return
}

//GetPopinMoney 爬取所需要的数据（消耗额）,核心函数,不要手贱的去"优化"
func getPopinMoney(popcookie, account, campaign string) int {
	campaignurl := "https://dashboard.popin.cc/discovery/accounts-tw/index.php/" + account + "/c/" + campaign + "/getDiscoveryReports"
	referer := "https://dashboard.popin.cc/discovery/accounts-tw/index.php/" + account + "/manageAgency?subPage=/" + account + "/campaigns/listCampaign"

	client := &http.Client{}
	var resp *http.Response
	req, err := http.NewRequest("GET", campaignurl, nil)
	if err != nil {
		fmt.Println("http.NewRequest Failed:", err)
		return 0
	}

	req.Header.Add("Host", "dashboard.popin.cc")
	req.Header.Add("User-Agent", function.GetRandomUserAgent())
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.8,zh-TW;q=0.7,zh-HK;q=0.5,en-US;q=0.3,en;q=0.2")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Referer", referer)
	req.Header.Add("Cookie", popcookie)
	req.Header.Add("Pragma", "no-cache")
	req.Header.Add("Cache-Control", "no-cache")

	//发送请求数据之服务器，并获取响应数据
	resp, err = client.Do(req)
	if err != nil {
		fmt.Println("Client Do Failed:", err)
		return 0
	}
	if resp.StatusCode != 200 {
		fmt.Println("status code error: ", resp.StatusCode, resp.Status)
		return 0
	}
	defer resp.Body.Close()

	//var body string
	//switch resp.Header.Get("Content-Encoding") {
	//case "gzip":
	//	reader, _ := gzip.NewReader(resp.Body)
	//	for {
	//		buf := make([]byte, 1024)
	//		n, err := reader.Read(buf)
	//
	//		if err != nil && err != io.EOF {
	//			panic(err)
	//		}
	//
	//		if n == 0 {
	//			break
	//		}
	//		body += string(buf)
	//	}
	//default:
	//	bodyByte, _ := ioutil.ReadAll(resp.Body)
	//	body = string(bodyByte)
	//}

	bodyByte, _ := ioutil.ReadAll(resp.Body)
	//开始处理信息
	var popJSON popinJSONObj
	//解码json数据，将字节切片映射到指定结构上
	e := json.Unmarshal(bodyByte, &popJSON)
	if e != nil {
		fmt.Println("json.Unmarshal failed")
		return 0
	}

	//fmt.Println(popJSON.TotalMoney)

	moneystring := strings.Replace(popJSON.TotalMoney, ",", "", -1)
	moneystring = moneystring[:len(moneystring)-3]

	moneyint, err := strconv.Atoi(moneystring)
	if err != nil {
		fmt.Println("strconv.Atoi failed:", err)
		return 0
	}

	return moneyint
}

//GetPopinCookie 获取Cookie
//func (pop *POP) getPopinCookie()(popcookie string) {
//	loginurl := "https://dashboard.popin.cc/discovery/accounts-tw/index.php"
//	// 填充表单，类似于net/url
//	args := &fasthttp.Args{}
//	args.Add("_token",pop.LoginToken)
//	args.Add("userid", pop.Account)
//	args.Add("password", pop.Password)
//	args.Add("autoLoginEnable", "on")
//
//	req := fasthttp.AcquireRequest()
//	resp := fasthttp.AcquireResponse()
//	defer func(){
//		// 用完需要释放资源
//		fasthttp.ReleaseResponse(resp)
//		fasthttp.ReleaseRequest(req)
//	}()
//	// 默认是application/x-www-form-urlencoded
//	req.Header.SetContentType("application/json")
//	req.Header.SetMethod("POST")
//	req.Header.Set("Host", "dashboard.popin.cc")
//	req.Header.Set("User-Agent", function.GetRandomUserAgent())
//	req.Header.Set("Origin", "https://dashboard.popin.cc")
//	req.Header.Set("Referer", "https://dashboard.popin.cc/discovery/accounts-tw/index.php")
//	req.Header.Set("Pragma", "no-cache")
//	req.Header.Set("Accept", "*/*")
//
//	req.SetRequestURI(loginurl)
//
//	fasthttp.Do(req,resp)
//
//	status, Resp, err := fasthttp.Post(nil, loginurl, args)
//	if err != nil||status != fasthttp.StatusOK{
//		fmt.Println("请求失败:", err,",状态值：",status)
//		return
//	}
//	fmt.Println(string(Resp))
//
//	//data := make(url.Values)
//	//data["userid"] = []string{pop.Account}
//	//data["password"] = []string{pop.Password}
//	//data["autoLoginEnable"] = []string{"on"}
//	//res, err := http.PostForm(loginurl, data)
//	////设置http中header参数，可以再此添加cookie等值
//	//res.Header.Add("Host", "dashboard.popin.cc")
//	//res.Header.Add("User-Agent", function.GetRandomUserAgent())
//	//res.Header.Add("Origin", "https://dashboard.popin.cc")
//	//res.Header.Add("Referer", "https://dashboard.popin.cc/discovery/accounts-tw/index.php")
//	//res.Header.Add("Pragma", "no-cache")
//	//res.Header.Add("Accept", "*/*")
//	//res.Header.Add("User-Agent", "***")
//	//res.Header.Add("http.socket.timeou", 5000)
//	//if err != nil {
//	//	fmt.Println(err.Error())
//	//	return
//	//}
//	//defer res.Body.Close()
//	//cookieSlices := res.Header["Set-Cookie"]
//	//PHPSESSID := cookieSlices[0]
//	//laravel_session := cookieSlices[1]
//	//reg1 := regexp.MustCompile("PHPSESSID=\\S*")
//	//reg2 := regexp.MustCompile("laravel_session=\\S*")
//	//cookie1 := reg1.FindAllString(PHPSESSID, -1)
//	//cookie2 := reg2.FindAllString(laravel_session, -1)
//	//var buffer bytes.Buffer
//	//buffer.WriteString(cookie1[0])
//	//buffer.WriteString(" ")
//	//buffer.WriteString(cookie2[0])
//	//buffer.WriteString(" __stripe_mid=abbaa6ef-6d4b-44f7-82b7-6f3d7da3ca68; autoLoginEnable=true; ticket=94c85e7803; __zlcmid=xdj2wiKdOzUtvl")
//	//popcookie = buffer.String()
//	//if strings.TrimSpace(cookie);cookie==""{
//	//	return errors.New("cookie not found")
//	//}
//	return
//}
