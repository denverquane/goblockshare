package blockchain

import (
	"testing"
)

func TestInitialBlock(t *testing.T) {
	block := InitialBlock(nil)

	if block.Index != 0 || calcHash(block) != block.Hash {
		t.Fail()
	}
}

func TestIsBlockSequenceValid(t *testing.T) {
	block := InitialBlock(nil)
	failBlock := InitialBlock(nil)

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
	failBlock.Hash = calcHash(failBlock)
	if !IsBlockSequenceValid(failBlock, block) {
		t.Fail()
	}
}

func TestGenerateBlock(t *testing.T) {
	array := make([]UserPassPair, 1)
	array[0] = UserPassPair{"user", "pass"}
	block := InitialBlock(array)
	newBlock, _ := GenerateBlock(block, []AuthTransaction{SampleAuthTransaction("user", "pass")})

	if !IsBlockSequenceValid(newBlock, block) {
		t.Fail()
	}
}
