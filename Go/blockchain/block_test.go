package blockchain

//import (
//	"testing"
//)

//func TestInitialBlock(t *testing.T) {
//	block := InitialBlock(UserPassPair{"user"})
//
//	if block.Index != 0 || block.PrevHash != "" || calcHash(block) != block.Hash{
//		t.Fail()
//	}
//}
//
//func TestIsBlockSequenceValid(t *testing.T) {
//	block := InitialBlock()
//	failBlock := InitialBlock()
//
//	if IsBlockSequenceValid(block, failBlock) {
//		t.Fail()
//	}
//
//	//changing just the index shouldnt fix it
//	failBlock.Index = block.Index + 1
//	if IsBlockSequenceValid(failBlock, block) {
//		t.Fail()
//	}
//
//	//changing the prev hash as well shouldnt fix it
//	failBlock.PrevHash = block.Hash
//	if IsBlockSequenceValid(failBlock, block) {
//		t.Fail()
//	}
//
//	//Rehashing *should* fix it!
//	failBlock.Hash = calcHash(failBlock)
//	if !IsBlockSequenceValid(failBlock, block) {
//		t.Fail()
//	}
//
//}
//
//func TestGenerateBlock(t *testing.T) {
//	block := InitialBlock()
//	newBlock := GenerateBlock(block, SampleTransaction())
//
//	if !IsBlockSequenceValid(newBlock, block) {
//		t.Fail()
//	}
//}
