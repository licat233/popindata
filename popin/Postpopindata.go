package popin

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"popindata/function"
	"regexp"
	"strconv"
	"strings"
)

//GetAllMoneyPost 获取popin指定campaign的总消耗,最常用的
func (pop *POP) GetAllMoneyPost(d string) (money int) {
	for _, v := range pop.CampaignList {
		j := Postpopindata(pop.Cookie, pop.Account, v, d)
		money += j
	}
	return
}

//Postpopindata 爬取所需要的数据（消耗额）,核心函数,不要手贱的去"优化"
func Postpopindata(popcookie, account, campaign, thedate string) int {
	campaignurl := "https://dashboard.popin.cc/discovery/accounts-tw/index.php/" + account + "/c/" + campaign + "/getCampaignDetails"
	referer := "https://dashboard.popin.cc/discovery/accounts-tw/index.php/" + account + "/manageAgency?subPage=/" + account + "/campaigns/listCampaign"
	client := &http.Client{}

	postdata := strings.NewReader("start=" + thedate + "&stop=" + thedate)
	req, err := http.NewRequest("POST", campaignurl, postdata)
	if err != nil {
		fmt.Println("http.NewRequest Failed:", err)
		return 0
	}

	req.Header.Set("Host", "dashboard.popin.cc")
	req.Header.Set("User-Agent", function.GetRandomUserAgent())
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Referer", referer)
	req.Header.Set("Cookie", popcookie)
	req.Header.Set("Origin", "https://dashboard.popin.cc")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("X-CSRF-TOKEN", "IYmSCESEIwZImiLPW53RTzEeoJmNnaN0EXUiAAYk")
	//发送请求数据之服务器，并获取响应数据
	var resp *http.Response
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

	bodyByte, _ := ioutil.ReadAll(resp.Body)
	body := string(bodyByte)
	r := regexp.MustCompile(`"charge":{(.*?)}`)
	d := r.FindString(body)
	r1 := regexp.MustCompile(`"mobile_click":(.*?),`)
	d1 := r1.FindString(d)
	r2 := regexp.MustCompile(`[0-9]+`)
	d2 := r2.FindString(d1)
	if len(d2) == 0 {
		return 0
	}
	d3, ok := strconv.Atoi(d2)
	if ok != nil {
		return 0
	}
	return d3
}
