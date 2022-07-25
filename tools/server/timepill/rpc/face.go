package rpc

import (
	"github.com/actliboy/hoper/server/go/lib/utils/log"
	"github.com/actliboy/hoper/server/go/lib/utils/net/http/client"
)

type FaceRecognitionRep struct {
	Code  int  `json:"code"`
	Count int  `json:"count"`
	Fount bool `json:"fount"`
}

func FaceRecognition(url string) {
	rep := FaceRecognitionRep{}
	err := client.DoGet("http://liov.xyz:5001?url="+url, &rep)
	if err != nil {
		log.Error(err)
	}
	log.Info(rep)
	if rep.Fount {

	}
}
