package bootstrap

import (
	"context"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/polipopoliko/go-block/pkg/blockchain"
	ws "github.com/polipopoliko/go-block/pkg/websocket"
	"github.com/spf13/viper"
)

type Container struct {
	ctx    context.Context
	infura *ethclient.Client
	hub    *ws.Hub
	bc     *blockchain.Blockchain
	ws     *ws.Client
}

func Init() *Container {
	initConfig()

	return &Container{
		ctx: context.Background(),
	}
}

func initConfig() {
	// TODO set up option for remote config
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	viper.ReadInConfig()
}

func (c *Container) Context() context.Context {
	return c.ctx
}

func (c *Container) Close() {
	if c.infura != nil {
		c.infura.Close()
	}
}

func (c *Container) GetHub() *ws.Hub {
	if c.hub != nil {
		return c.hub
	}

	c.hub = ws.NewHub()
	return c.hub
}

func (c *Container) GetBlockChain() *blockchain.Blockchain {
	if c.bc != nil {
		return c.bc
	}

	c.bc = blockchain.NewBlockchain(viper.GetInt32("blockchain.difficulty"))

	return c.bc
}
