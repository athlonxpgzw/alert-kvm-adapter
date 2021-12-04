package kvm

type Data struct {
	Topic     string     `json:"topic"`
	Params    Parameters `json:"params"`
	Timestamp int64      `json:"timestamp"`
}

type Parameters struct {
	Zone            string `json:"zone"`
	Level           string `json:"level"`
	Mesg            string `json:"message"`
	NodeID          string `json:"node_id"`
	NodeIP          string `json:"node_ip"`
	Duration        string `json:"duration"`
	DurationRaw     int64  `json:"duration_raw"`
	Synopsis        string `json:"synopsis"`
	NodeType        int64  `json:"node_type"`
	AlertType       string `json:"alert_type"`
	Suggestion      string `json:"suggestion"`
	EventCount      int64  `json:"event_count"`
	AggsEventID     int64  `json:"aggs_event_id"`
	SendDuration    string `json:"send_duration"`
	SendDurationRaw int64  `json:"send_duration_raw"`
	I18n            I18n   `json:"synopsis_i18n"`
	FirstNotify     string `json:"first_notify_time"`
	LastestNotify   string `json:"latest_notify_time"`
	StatusTime      string `json:"status_time"`
}

type I18n struct {
	EN   string `json:"en"`
	ZHCN string `json:"zh-cn"`
}
