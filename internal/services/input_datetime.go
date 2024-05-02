package services

import (
	"fmt"
	"time"

	ws "saml.dev/gome-assistant/internal/websocket"
)

/* Structs */

type InputDatetime struct {
	conn *ws.WebsocketWriter
}

/* Public API */

func (ib InputDatetime) Set(entityId string, value time.Time) {
	req := NewBaseServiceRequest(entityId)
	req.Domain = "input_datetime"
	req.Service = "set_datetime"
	req.ServiceData = map[string]any{
		"timestamp": fmt.Sprint(value.Unix()),
	}

	ib.conn.WriteMessage(req)
}

func (ib InputDatetime) Reload() {
	req := NewBaseServiceRequest("")
	req.Domain = "input_datetime"
	req.Service = "reload"
	ib.conn.WriteMessage(req)
}
