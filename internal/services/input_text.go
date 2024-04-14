package services

import (
	"saml.dev/gome-assistant/internal/websocket"
)

/* Structs */

type InputText struct {
	conn *websocket.Conn
}

func NewInputText(conn *websocket.Conn) *InputText {
	return &InputText{
		conn: conn,
	}
}

/* Public API */

func (ib InputText) Set(entityID string, value string) {
	req := CallServiceRequest{}
	req.Domain = "input_text"
	req.Service = "set_value"
	req.Target.EntityID = entityID
	req.ServiceData = map[string]any{
		"value": value,
	}

	ib.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}

func (ib InputText) Reload() {
	req := CallServiceRequest{}
	req.Domain = "input_text"
	req.Service = "reload"

	ib.conn.Send(func(mw websocket.MessageWriter) error {
		req.ID = mw.NextID()
		return mw.SendMessage(req)
	})
}
