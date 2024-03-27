package server

import (
	"log"
	"net/http"

	"github.com/polipopoliko/go-block/bootstrap"
	"github.com/polipopoliko/go-block/internal/handler"
	"github.com/polipopoliko/go-block/internal/usecase/websocket"
	ws "github.com/polipopoliko/go-block/pkg/websocket"
	"github.com/spf13/viper"
)

func Serve(c *bootstrap.Container) {

	var (
		readBuffer  = viper.GetInt("server.websocket.read_buffer")
		writeBuffer = viper.GetInt("server.websocket.write_buffer")
		usecase     = websocket.NewWebsocketUsecase(c.GetBlockChain(), c.GetHub(), viper.GetDuration("server.websocket.write_delay"))
		handler     = handler.NewWebSocketHandler(c.GetHub(), usecase, readBuffer, writeBuffer)
	)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		var (
			wsCl = ws.NewWebsocketClient(readBuffer, writeBuffer, c.GetHub(), w, r)
		)
		handler.ConnectWS(r.Context(), wsCl)
	})

	go c.GetHub().Run()
	server := &http.Server{
		Addr:              viper.GetString("server.websocket.port"),
		ReadHeaderTimeout: viper.GetDuration("server.websocket.read_header_timeout"),
	}

	if err := server.ListenAndServe(); err != nil {
		log.Println("failed to listen and serve ws", err)
	}
}
