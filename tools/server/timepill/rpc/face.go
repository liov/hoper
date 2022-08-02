package rpc

import (
	"github.com/actliboy/hoper/server/go/lib/utils/log"
	"github.com/actliboy/hoper/server/go/lib/utils/net/http/client"
)

type FaceDetectionRep struct {
	Code  int  `json:"code"`
	Count int  `json:"count"`
	Found bool `json:"found"`
}

func FaceDetection(url string) *FaceDetectionRep {
	rep := FaceDetectionRep{}
	err := client.Get("http://liov.xyz:5001?url="+url, &rep)
	if err != nil {
		log.Error(err)
	}
	log.Info(rep)
	return &rep
}
