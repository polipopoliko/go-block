package websocket

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/polipopoliko/go-block/model"
	"github.com/polipopoliko/go-block/pkg/blockchain"
	ws "github.com/polipopoliko/go-block/pkg/websocket"
)

type WebsocketUsecase struct {
	bc           *blockchain.Blockchain
	hub          *ws.Hub
	publishDelay time.Duration
}

func NewWebsocketUsecase(bc *blockchain.Blockchain, hub *ws.Hub, delay time.Duration) *WebsocketUsecase {
	return &WebsocketUsecase{bc: bc, hub: hub, publishDelay: delay}
}

func (uc *WebsocketUsecase) Subscribe(ctx context.Context, cl *ws.Client) {

	defer func() {
		cl.Close()
		uc.hub.Unregister(cl)
	}()

	for {
		select {
		default:
			_, r, err := cl.Conn.ReadMessage()
			if err != nil {
				log.Println("err", err)
				if err == websocket.ErrCloseSent {
					break
				}

				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("error: %v", err)
				}
				return
			}
			msg := bytes.TrimSpace(bytes.Replace(r, []byte(string(model.Delim)), []byte(string(" ")), -1))

			var data blockchain.BlockData
			if err := json.Unmarshal(msg, &data); err != nil {
				log.Println("failed to unmarshal message ", string(msg), err)
				continue
			}

			uc.bc.AppendBlock(data)
		}
	}
}

func (uc *WebsocketUsecase) Publish(ctx context.Context, cl *ws.Client) {
	for {
		<-time.After(uc.publishDelay)

		byt, err := json.Marshal(uc.bc.Chain)
		if err != nil {
			log.Println("failed to marshal block chain", err)
			continue
		}

		cl.Conn.WriteMessage(websocket.PingMessage, []byte{})
		if err := cl.Conn.WriteMessage(websocket.TextMessage, byt); err != nil {
			if strings.Contains(err.Error(), "closed network connection") {
				log.Println("failed to init next writer", err)
				return
			}
			continue
		}
	}
}
