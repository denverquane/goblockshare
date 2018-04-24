package blockchain

import (
	"testing"
)

func TestInitialBlock(t *testing.T) {
	block := InitialBlock()

	if block.Index != 0 || block.Hash() != block.Hashed {
		t.Fail()
	}
}

func TestIsBlockSequenceValid(t *testing.T) {
	block := InitialBlock()
	failBlock := InitialBlock()

	if IsBlockSequenceValid(block, failBlock) {
		t.Fail()
	}

	//changing just the index shouldnt fix it
	failBlock.Index = block.Index + 1
	if IsBlockSequenceValid(failBlock, block) {
		t.Fail()
	}

	//changing the prev hash as well shouldnt fix it
	failBlock.PrevHash = block.Hashed
	if IsBlockSequenceValid(failBlock, block) {
		t.Fail()
	}

	//Rehashing *should* fix it!
	failBlock.Hashed = failBlock.Hash()
	if !IsBlockSequenceValid(failBlock, block) {
		t.Fail()
	}
}
