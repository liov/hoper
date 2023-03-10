package client

import "testing"

func TestFetch(t *testing.T) {
	_, err := GetFile("https://upos-sz-mirrorcos.bilivideo.com/upgcxcode/43/48/173034843/173034843-1-30016.m4s?e=ig8euxZM2rNcNbdlhoNvNC8BqJIzNbfqXBvEqxTEto8BTrNvN0GvT90W5JZMkX_YN0MvXg8gNEV4NC8xNEV4N03eN0B5tZlqNxTEto8BTrNvNeZVuJ10Kj_g2UB02J0mN0B5tZlqNCNEto8BTrNvNC7MTX502C8f2jmMQJ6mqF2fka1mqx6gqj0eN0B599M=&uipk=5&nbs=1&deadline=1678438731&gen=playurlv2&os=cosbv&oi=2032348538&trid=68d2d7928903414dab9b545c7326cb52u&mid=0&platform=pc&upsig=474cdcddd021a71d48408637d0e88cf8&uparams=e,uipk,nbs,deadline,gen,os,oi,trid,mid,platform&bvc=vod&nettype=0&orderid=0,3&buvid=&build=0&agrr=1&bw=49052&logo=80000000")
	if err != nil {
		t.Log(err)
	}
}
