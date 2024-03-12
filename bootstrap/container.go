package bootstrap

import (
	"context"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/viper"
)

type Container struct {
	ctx    context.Context
	infura *ethclient.Client
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
