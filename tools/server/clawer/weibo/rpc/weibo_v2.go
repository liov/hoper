package rpc

type WeiboV2 struct {
	Cards []*CardGroupV2 `json:"cards"`
}

type CardGroupV2 struct {
	Mblog *MblogV2 `json:"mblog"`
}

type MblogV2 struct {
	Id  string `json:"id"`
	Mid string `json:"mid"`
}
