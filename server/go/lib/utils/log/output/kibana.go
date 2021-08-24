package output

import (
	"time"

	"go.uber.org/zap/zapcore"
)

type Kibana struct {
	Interface   string        `json:"interface,omitempty"`
	ProcessTime time.Duration `json:"processTime,omitempty"`
	Param       string        `json:"param,omitempty"`
	Result      string        `json:"result,omitempty"`
	Other       string        `json:"other,omitempty"`
}

func (k *Kibana) MarshalLogObject(e zapcore.ObjectEncoder) error {
	if k.Interface != "" {
		e.AddString("interface", k.Interface)
	}
	if k.ProcessTime != 0 {
		e.AddDuration("processTime", k.ProcessTime)
	}
	if k.Param != "" {
		e.AddString("body", k.Param)
	}
	if k.Result != "" {
		e.AddString("result", k.Result)
	}
	if k.Other != "" {
		e.AddString("other", k.Other)
	}
	return nil
}
