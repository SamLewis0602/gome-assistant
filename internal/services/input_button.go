package services

import (
	"saml.dev/gome-assistant/internal/websocket"
)

/* Structs */

type InputButton struct {
	conn *websocket.Conn
}

func NewInputButton(conn *websocket.Conn) *InputButton {
	return &InputButton{
		conn: conn,
	}
}

/* Public API */

func (ib InputButton) Press(entityId string) {
	req := NewBaseServiceRequest(ib.conn, entityId)
	req.Domain = "input_button"
	req.Service = "press"

	ib.conn.Send(func(mw websocket.MessageWriter) error {
		req.Id = mw.NextID()
		return mw.SendMessage(req)
	})
}

func (ib InputButton) Reload() {
	req := NewBaseServiceRequest(ib.conn, "")
	req.Domain = "input_button"
	req.Service = "reload"

	ib.conn.Send(func(mw websocket.MessageWriter) error {
		req.Id = mw.NextID()
		return mw.SendMessage(req)
	})
}
