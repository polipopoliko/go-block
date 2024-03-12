package nodes

import (
	"crypto/rand"
	"io"
	"log"
	mrand "math/rand"
	"strconv"
	"strings"

	libp2p "github.com/libp2p/go-libp2p"
	crypto "github.com/libp2p/go-libp2p/core/crypto"
	host "github.com/libp2p/go-libp2p/core/host"
	ma "github.com/multiformats/go-multiaddr"
)

type Nodes struct {
	ConnectedPorts []string
	private        crypto.PrivKey
	host           host.Host
	addr           ma.Multiaddr
	port           int
	seed           int
}

func NewNodes(port int, seed int, connectedPorts []string) *Nodes {

	if connectedPorts == nil {
		connectedPorts = []string{}
	}

	return &Nodes{
		ConnectedPorts: connectedPorts,
		port:           port,
		seed:           seed,
	}
}

func (n Nodes) GetConnectedPorts() []string {

	if n.ConnectedPorts == nil {
		return []string{}
	}

	res := []string{}
	for _, val := range n.ConnectedPorts {
		res = append(res, val)
	}
	return res
}

func (n Nodes) GetStringPort() string {
	return defaultAddr + strconv.Itoa(n.port)
}

func (n *Nodes) GetPrivateKey() (crypto.PrivKey, error) {

	if n.private == nil {
		var reader io.Reader
		if n.seed <= 0 {
			reader = rand.Reader
		}
		reader = mrand.New(mrand.NewSource(int64(n.seed)))

		pri, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, reader)
		if err != nil {
			return nil, err
		}

		n.private = pri
	}

	return n.private, nil
}

func (n *Nodes) GetHost() (host.Host, error) {

	if n.host != nil {
		return n.host, nil
	}

	pri, err := n.GetPrivateKey()
	if err != nil {
		log.Println("get private key errr", err)
		return nil, err
	}

	host, err := libp2p.New(libp2p.ListenAddrStrings(n.GetStringPort()), libp2p.Identity(pri))
	if err != nil {
		log.Println("libp2p", err)
		return nil, err
	}
	n.host = host

	log.Println("Node ID", n.host.ID().ShortString())

	return n.host, nil
}

func (n *Nodes) GetAddress() (ma.Multiaddr, error) {

	host, err := n.GetHost()
	if err != nil {
		return nil, err
	}

	addr, err := ma.NewMultiaddr("/ipfs/" + host.ID().String())
	if err != nil {
		return nil, err
	}

	// select the address starting with "ip4"
	for _, i := range host.Addrs() {
		if strings.HasPrefix(i.String(), "/ip4") {
			addr = i.Encapsulate(addr)
			break
		}
	}
	n.addr = addr

	log.Println("Node Address ", n.addr.String())
	return n.addr, err
}

func (n *Nodes) GetDialAddress() ([]ma.Multiaddr, error) {

	res := make([]ma.Multiaddr, len(n.GetConnectedPorts()))
	for i, val := range n.GetConnectedPorts() {
		ma, err := ma.NewMultiaddr(val)
		if err != nil {
			log.Println("dialing address failed", val, err)
			continue
		}
		res[i] = ma
	}
	return res, nil
}
