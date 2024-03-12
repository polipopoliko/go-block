package usecase

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/polipopoliko/go-block/pkg/blockchain"
)

type Reader interface {
	Reader(ctx context.Context, data string) error
}

type Writer interface {
	Writer(ctx context.Context) ([]byte, error)
}

type ReaderWriter interface {
	Reader
	Writer
}

type StreamReaderWriter struct {
	Mutex      *sync.Mutex
	Blockchain *blockchain.Blockchain
}

func NewStreamUsecase(mtx *sync.Mutex, bc *blockchain.Blockchain) *StreamReaderWriter {
	return &StreamReaderWriter{
		Mutex:      mtx,
		Blockchain: bc,
	}
}

func (r *StreamReaderWriter) Reader(ctx context.Context, data string) error {

	var block *blockchain.Blockchain
	if err := json.Unmarshal([]byte(data), &block); err != nil {
		log.Println("invalid unmarshal input", err, data)
		return err
	}

	r.Mutex.Lock()
	if len(block.Chain) > len(r.Blockchain.Chain) {
		r.Blockchain = block

		log.Println("synchronizing locak blockchain with incoming chain ", r.Blockchain.Chain, block)
	}
	r.Mutex.Unlock()

	return nil
}

func (r *StreamReaderWriter) Writer(ctx context.Context) ([]byte, error) {
	bytes, err := json.Marshal(r.Blockchain)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}
