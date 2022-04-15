package event

import "github.com/quanxiang-cloud/form/internal/models"

type DaprEvent struct {
	Topic           string `json:"topic"`
	Pubsubname      string `json:"pubsubname"`
	Traceid         string `json:"traceid"`
	ID              string `json:"id"`
	Datacontenttype string `json:"datacontenttype"`
	Data            Data   `json:"data"`
	Type            string `json:"type"`
	Specversion     string `json:"specversion"`
	Source          string `json:"source"`
}

type Data struct {
	*UserSpec   `json:"user,omitempty"`
	*PermitSpec `json:"permit,omitempty"`
}

type UserSpec struct {
	RoleID string `json:"roleID"`
	UserID string `json:"userID"`
	AppID  string `json:"appID"`
	Action string `json:"action"`
}

type PermitSpec struct {
	RoleID    string             `json:"roleID"`
	Path      string             `json:"path"`
	Condition models.Condition   `json:"condition"`
	Params    models.FiledPermit `json:"params"`
	Response  models.FiledPermit `json:"response"`
	Action    string             `json:"action"`
}
