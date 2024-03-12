package blockchain

import "time"

type Blockchain struct {
	GenesisBlock Block
	Chain        []Block
	Difficulty   int32
}

func NewBlockchain() *Blockchain {
	genesis := Block{
		Timestamp: time.Now().Unix(),
		Hash:      "0",
		Height:    0,
	}
	genesis.calculateHash()

	return &Blockchain{
		GenesisBlock: genesis,
	}
}

func (bc *Blockchain) AppendBlock(data BlockData) {
	var (
		lb = bc.Chain[len(bc.Chain)-1]
		nb = Block{
			Data:         data,
			PreviousHash: lb.Hash,
			Timestamp:    time.Now().Unix(),
			Height:       lb.Height,
			Pow:          0,
		}
	)

	nb.mine(int(bc.Difficulty))
	bc.Chain = append(bc.Chain, nb)
}

func (bc *Blockchain) IsValid(data Block) bool {
	oldBlock := bc.Chain[data.Height-1]

	if oldBlock.Height+1 != data.Height {
		return false
	}

	if oldBlock.Hash != data.PreviousHash {
		return false
	}

	if data.calculateHash() != data.Hash {
		return false
	}

	return true
}
