package firefly

type Data struct {
	Name      string      `json:"alertName"`
	Desc      string      `json:"alertDesc"`
	Status    string      `json:"alertStatus"`
	ApplyType string      `json:"applyType"`
	Key       string      `json:"alertKey"`
	Id        string      `json:"alertId"`
	MsgId     string      `json:"alertMsgId"`
	Level     string      `json:"alertLevel"`
	AlertFile []string    `json:"alertFile"`
	Notify    AlertNotify `json:"alertNofity"`
}

type AlertNotify struct {
	WeChats []string `json:"wechats"`
	Mobiles []string `json:"mobiles"`
	Emails  []string `json:"emails"`
}
