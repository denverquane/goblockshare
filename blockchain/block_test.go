package blockchain

import (
	"testing"
)

func TestInitialBlock(t *testing.T) {
	block := InitialBlock("")

	if block.Index != 0 || block.GetHash() != block.Hash {
		t.Fail()
	}
}

func TestIsBlockSequenceValid(t *testing.T) {
	block := InitialBlock("")
	failBlock := InitialBlock("")

	if IsBlockSequenceValid(block, failBlock) {
		t.Fail()
	}

	//changing just the index shouldnt fix it
	failBlock.Index = block.Index + 1
	if IsBlockSequenceValid(failBlock, block) {
		t.Fail()
	}

	//changing the prev hash as well shouldnt fix it
	failBlock.PrevHash = block.Hash
	if IsBlockSequenceValid(failBlock, block) {
		t.Fail()
	}

	//Rehashing *should* fix it!
	failBlock.Hash = failBlock.GetHash()
	if !IsBlockSequenceValid(failBlock, block) {
		t.Fail()
	}
}
