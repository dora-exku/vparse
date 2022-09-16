package vparse

type CallFunc func(...any) (string, error)

func New(v string) Parse {
	switch v {
	case "iqiyi":
		return &IqiyiParse{}
	case "tencent":
		return &TencentParse{}
	default:
		return nil
	}
}
