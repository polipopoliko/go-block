package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/multiformats/go-multiaddr"
	"github.com/polipopoliko/go-block/bootstrap"
	"github.com/polipopoliko/go-block/internal/nodes"
	"github.com/polipopoliko/go-block/internal/stream"
	usecase "github.com/polipopoliko/go-block/internal/usecase/stream"
	"github.com/polipopoliko/go-block/pkg/blockchain"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	gologging "github.com/whyrusleeping/go-logging"
)

var (
	seed    int
	port    int
	dial    []string
	nodeCmd = func(container *bootstrap.Container) *cobra.Command {
		return &cobra.Command{
			PreRun: func(cmd *cobra.Command, args []string) {
				log.Println("running node command")
			},
			Short: "running nodes on p2p network",
			Use:   "p2p",
			PreRunE: func(cmd *cobra.Command, args []string) error {
				gologging.SetLevel(gologging.DEBUG, "p2p")
				return nil
			},
			Run: func(cmd *cobra.Command, args []string) {

				var (
					node        = nodes.NewNodes(port, seed, dial)
					ctx, cancel = context.WithCancel(context.Background())
					mtx         = &sync.Mutex{}
					bc          = blockchain.NewBlockchain()
				)
				defer cancel()

				host, err := node.GetHost()
				if err != nil {
					log.Fatal("failed to get host", err)
				}

				if _, err = node.GetAddress(); err != nil {
					log.Fatal("failed to get address", err)
				}

				var (
					uc   = usecase.NewStreamUsecase(mtx, bc)
					strm = stream.NewStream(ctx, uc, mtx, viper.GetDuration("configuration.sync_delay"))
				)

				host.SetStreamHandler("/p2p/1.0.0", strm.StreamHandler())

				dialAddr, err := node.GetDialAddress()
				if err != nil {
					log.Println("failed to get dial address", err)
				}
				if len(dialAddr) > 0 {
					pid, err := dialAddr[0].ValueForProtocol(multiaddr.P_IPFS)
					if err != nil {
						log.Fatalln(err)
					}

					peerID, err := peer.Decode(pid)
					if err != nil {
						log.Fatalln(err)
					}

					targetPeerAddr, _ := multiaddr.NewMultiaddr(
						fmt.Sprintf("/ipfs/%s", peerID.String()))
					host.Peerstore().AddAddr(peerID, dialAddr[0].Decapsulate(targetPeerAddr), peerstore.PermanentAddrTTL)

					s, err := host.NewStream(context.Background(), peerID, "/p2p/1.0.0")
					if err != nil {
						log.Fatalln(err)
					}

					pmtx := &sync.Mutex{}
					pstrm := stream.NewStream(ctx, uc, pmtx, viper.GetDuration("configuration.sync_delay"))
					pstrm.StreamHandler()(s)
				}

				sig := make(chan os.Signal, 1)
				signal.Notify(sig, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)

				<-sig
			},
		}
	}
)
