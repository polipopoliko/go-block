package bootstrap

import (
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/viper"
)

func (c *Container) Infura() *ethclient.Client {

	if c.infura == nil {
		ecl, err := ethclient.DialContext(c.Context(), viper.GetString("infura.url"))
		if err != nil {
			log.Fatal("failed to connect to infura")
		}

		nb, err := ecl.BlockNumber(c.ctx)
		if err != nil {
			log.Fatal("failed to ping infura")
		}
		log.Println("number of blocks", nb)

		c.infura = ecl
	}

	return c.infura
}
