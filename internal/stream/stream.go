package stream

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p/core/network"
	usecase "github.com/polipopoliko/go-block/internal/usecase/stream"
	"github.com/polipopoliko/go-block/model"
	"github.com/polipopoliko/go-block/pkg/blockchain"
)

type Stream struct {
	BC    []blockchain.Block
	Mutex *sync.Mutex
	*bufio.ReadWriter
	stream    network.Stream
	Rw        usecase.ReaderWriter
	SyncDelay time.Duration
	Ctx       context.Context
}

func NewStream(ctx context.Context, rw usecase.ReaderWriter, mutex *sync.Mutex, dur time.Duration) *Stream {
	stream := &Stream{
		Ctx:       ctx,
		Mutex:     mutex,
		BC:        []blockchain.Block{},
		Rw:        rw,
		SyncDelay: dur,
	}

	return stream
}

func (s *Stream) StreamHandler() func(network.Stream) {
	return func(sh network.Stream) {
		s.stream = sh
		s.ReadWriter = bufio.NewReadWriter(bufio.NewReader(s.stream), bufio.NewWriter(s.stream))

		go s.Read(s.Ctx)
		go s.Write(s.Ctx)
	}
}

func (s *Stream) Read(ctx context.Context) {
	defer s.close()
	for {
		select {
		case <-ctx.Done():
			log.Println("context done, closing reader")
			return
		default:
			data, err := s.ReadWriter.ReadString(model.Delim)
			if err != nil || data == "" || data == string(model.Delim) {
				if strings.Contains(err.Error(), "stream reset") {
					time.Sleep(1 * time.Second)
				}
				log.Println("read string err or data is empty", err, data)
				continue
			}

			log.Println("reader err ", s.Rw.Reader(ctx, data))
		}
	}
}

func (s *Stream) Write(ctx context.Context) {
	defer s.close()
	for {
		select {
		case <-ctx.Done():
			log.Println("context done, closing writer")
			return
		case <-time.After(s.SyncDelay):
			data, err := s.Rw.Writer(ctx)
			if err != nil {
				log.Println("writer error", err, data)
				continue
			}

			cnt, err := s.ReadWriter.WriteString(fmt.Sprintf("%s\n", string(data)))
			log.Printf("data written %v, err %v\n", cnt, err)
			if err != nil {
				continue
			}

			log.Println("writer flush err", s.ReadWriter.Flush())
		}
	}
}

func (s *Stream) close() {
	if s.stream == nil {
		return
	}
	s.close()
}
