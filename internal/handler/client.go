package handler

import (
	"context"

	ws "github.com/polipopoliko/go-block/pkg/websocket"
)

type WebSocketProcessor interface {
	Subscribe(ctx context.Context, cl *ws.Client)
	Publish(ctx context.Context, cl *ws.Client)
}

type WebSocketHandler struct {
	hub                     *ws.Hub
	readBuffer, writeBuffer int
	uc                      WebSocketProcessor
}

func NewWebSocketHandler(hub *ws.Hub, uc WebSocketProcessor, readBuffer, writeBuffer int) *WebSocketHandler {
	return &WebSocketHandler{
		hub: hub,
		uc:  uc,
	}
}

func (h *WebSocketHandler) ConnectWS(ctx context.Context, cl *ws.Client) {
	h.hub.Register(cl)

	go h.uc.Publish(ctx, cl)
	go h.uc.Subscribe(ctx, cl)
}
