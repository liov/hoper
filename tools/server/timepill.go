package main

import (
	"github.com/actliboy/hoper/server/go/lib/utils/log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func getRequest(id int) (*http.Request, error) {
	req, err := http.NewRequest("GET", "https://timepill.net/diary/"+strconv.Itoa(id), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("authority", "timepill.net")
	req.Header.Add("sec-ch-ua", "\"Chromium\";v=\"21\", \" Not;A Brand\";v=\"99\"")
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.69 Safari/537.36")
	req.Header.Add("sec-ch-ua-platform", "\"Windows\"")
	req.Header.Add("accept", "*/*")
	req.Header.Add("sec-fetch-site", "same-origin")
	req.Header.Add("sec-fetch-mode", "cors")
	req.Header.Add("sec-fetch-dest", "empty")
	req.Header.Add("referer", "https://timepill.net/")
	req.Header.Add("accept-language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Add("cookie", "BAIDU_SSP_lcr=https://www.baidu.com/link?url=Uua85jNMPG50cziVC70uyMuyazxJU0g_vBFJdpqLvuHix5VtsEiUC5FiymVbYevG&wd=&eqid=d0e82a380001efb70000000661ea68df; Hm_lvt_524f61d4da0fc26075b878c3d69484db=1642752226; XSRF-TOKEN=wU38YxWWrFcbJ3NbF9CyvOdqzHzLWv3yCL5fprAW; laravel_session=g2nkB6YyDAFESnrxW7HJRyfllOdxZqteL54F1kEv; Hm_lpvt_524f61d4da0fc26075b878c3d69484db=1642752240")
	http.DefaultClient.Timeout = 300 * time.Second
	return req, nil
}

func commentRequest(id int) (*http.Request, error) {
	req, err := http.NewRequest("POST", "https://timepill.net/comment/add/"+strconv.Itoa(id), strings.NewReader("content=%E5%A5%BD%E4%B9%85%E4%B8%8D%E8%A7%81&recipient_id=0&has_cmt="))
	if err != nil {
		return nil, err
	}
	req.Header.Add("authority", "timepill.net")
	req.Header.Add("sec-ch-ua", "\"Chromium\";v=\"21\", \" Not;A Brand\";v=\"99\"")
	req.Header.Add("accept", "*/*")
	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	req.Header.Add("x-requested-with", "XMLHttpRequest")
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.69 Safari/537.36")
	req.Header.Add("sec-ch-ua-platform", "\"Windows\"")
	req.Header.Add("accept", "*/*")
	req.Header.Add("sec-fetch-site", "same-origin")
	req.Header.Add("sec-fetch-mode", "cors")
	req.Header.Add("sec-fetch-dest", "empty")
	req.Header.Add("referer", "https://timepill.net/")
	req.Header.Add("accept-language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Add("cookie", "BAIDU_SSP_lcr=https://www.baidu.com/link?url=Uua85jNMPG50cziVC70uyMuyazxJU0g_vBFJdpqLvuHix5VtsEiUC5FiymVbYevG&wd=&eqid=d0e82a380001efb70000000661ea68df; Hm_lvt_524f61d4da0fc26075b878c3d69484db=1642752226; XSRF-TOKEN=wU38YxWWrFcbJ3NbF9CyvOdqzHzLWv3yCL5fprAW; laravel_session=g2nkB6YyDAFESnrxW7HJRyfllOdxZqteL54F1kEv; Hm_lpvt_524f61d4da0fc26075b878c3d69484db=1642752240")
	http.DefaultClient.Timeout = 300 * time.Second
	return req, nil
}

var id = 19607110

func main() {
	client := http.DefaultClient

	for {
		getReq, err := getRequest(id)
		if err != nil {
			log.Error(err)
			continue
		}
		resp, err := client.Do(getReq)
		if err != nil {
			log.Error(err, "id:", id)
			continue
		}
		if resp.StatusCode != 200 {
			resp.Body.Close()
			log.Errorf("返回错误，状态码：%d,id:%d", resp.StatusCode, id)
			id++
			continue
		}
		if resp.StatusCode == 200 {
			commentReq, err := commentRequest(id)
			if err != nil {
				log.Error(err)
				continue
			}
			resp, err = client.Do(commentReq)
			if err != nil {
				log.Error(err, "id:", id)
				continue
			}
			if resp.StatusCode != 200 {
				log.Errorf("返回错误，状态码：%d,id:%d", resp.StatusCode, id)
			}
			resp.Body.Close()
			id++
		}
	}
}
