package rpc

import (
	"github.com/liov/hoper/server/go/lib/utils/log"
	"github.com/liov/hoper/server/go/lib/utils/net/http/client"
)

type FaceDetectionRep struct {
	Code  int  `json:"code"`
	Count int  `json:"count"`
	Found bool `json:"found"`
}

func FaceDetection(url string) *FaceDetectionRep {
	rep := FaceDetectionRep{}
	err := client.SimpleGet("http://liov.xyz:5001?url="+url, &rep)
	if err != nil {
		log.Error(err)
	}
	log.Info(rep)
	return &rep
}
