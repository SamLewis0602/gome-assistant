package services

import (
	"saml.dev/gome-assistant/internal/websocket"
)

/* Structs */

type InputBoolean struct {
	conn *websocket.Conn
}

func NewInputBoolean(conn *websocket.Conn) *InputBoolean {
	return &InputBoolean{
		conn: conn,
	}
}

/* Public API */

func (ib InputBoolean) TurnOn(entityId string) {
	req := NewBaseServiceRequest(ib.conn, entityId)
	req.Domain = "input_boolean"
	req.Service = "turn_on"

	ib.conn.Send(func(mw websocket.MessageWriter) error {
		req.Id = mw.NextID()
		return mw.SendMessage(req)
	})
}

func (ib InputBoolean) Toggle(entityId string) {
	req := NewBaseServiceRequest(ib.conn, entityId)
	req.Domain = "input_boolean"
	req.Service = "toggle"

	ib.conn.Send(func(mw websocket.MessageWriter) error {
		req.Id = mw.NextID()
		return mw.SendMessage(req)
	})
}

func (ib InputBoolean) TurnOff(entityId string) {
	req := NewBaseServiceRequest(ib.conn, entityId)
	req.Domain = "input_boolean"
	req.Service = "turn_off"

	ib.conn.Send(func(mw websocket.MessageWriter) error {
		req.Id = mw.NextID()
		return mw.SendMessage(req)
	})
}

func (ib InputBoolean) Reload() {
	req := NewBaseServiceRequest(ib.conn, "")
	req.Domain = "input_boolean"
	req.Service = "reload"

	ib.conn.Send(func(mw websocket.MessageWriter) error {
		req.Id = mw.NextID()
		return mw.SendMessage(req)
	})
}
