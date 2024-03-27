package blockchain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAppendBlock(t *testing.T) {
	var (
		genesis = Block{
			Hash: "0",
		}
		bc = &Blockchain{
			GenesisBlock: genesis,
			Chain:        []Block{genesis},
			Difficulty:   1,
		}

		testCase = []BlockData{
			{
				Message: "test1",
			},
			{
				Message: "test2",
			},
			{
				Message: "test3",
			},
		}
	)

	for i, val := range testCase {

		if i == 0 {
			continue
		}

		bc.AppendBlock(val)

		assert.Equal(t, testCase[i].Message, bc.Chain[i].Data.Message)
		assert.NotEmpty(t, bc.Chain[i].Hash)
		assert.NotEmpty(t, bc.Chain[i].PreviousHash)
	}
}
