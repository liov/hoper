package pick

import (
	"encoding/json"
	"log"
)

type PMHeader struct {
	Key     string `json:"key"`
	Value   string `json:"value"`
	Type    string `json:"type"`
	Disable bool   `json:"disable"`
}

type PMBody struct {
	Mode string `json:"mode"`
	Raw  string `json:"raw"`
}

type PMRequest struct {
	Method      string     `json:"method"`
	Header      []PMHeader `json:"header"`
	Body        PMBody     `json:"body"`
	Description string     `json:"description"`
	URL         PMUrl      `json:"url"`
}

type PMUrl struct {
	Raw      string    `json:"raw"`
	Protocol string    `json:"protocol"`
	Host     []string  `json:"host"`
	Path     []string  `json:"path"`
	Query    []PMQuery `json:"query"`
}

type PMProtocolProfileBehavior struct {
	DisableBodyPruning bool `json:"disableBodyPruning"`
}

type PMApi struct {
	Name                    string                    `json:"name"`
	ProtocolProfileBehavior PMProtocolProfileBehavior `json:"protocolProfileBehavior"`
	Request                 PMRequest                 `json:"request"`
	Response                []interface{}             `json:"response"`
}

type PMCategory struct {
	Name                    string                    `json:"name"`
	ProtocolProfileBehavior PMProtocolProfileBehavior `json:"protocolProfileBehavior"`
	Items                   []PMApi                   `json:"item"`
}

type PMInfo struct {
	PostManId string `json:"_postman_id"`
	Name      string `json:"name"`
	Schema    string `json:"schema"`
}

type PMFile struct {
	Info  PMInfo       `json:"info"`
	Items []PMCategory `json:"item"`
}

type PMQuery struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func genPostman() string {
	pObj := PMFile{}
	pObj.Info.Schema = "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"

	ret, err := json.MarshalIndent(pObj, "", "    ")
	if err != nil {
		log.Println(err.Error())
	}
	return string(ret)
}
