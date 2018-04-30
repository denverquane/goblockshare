package transaction

import (
	"crypto/sha256"
	"math/big"
	"strconv"
)

type ChannelCreationTransaction struct {
	Origin                OriginInfo
	ChannelName           string
	TokenName             string
	RecipientsQuantityMap map[Base64Address]float64
	ExchangeRate          float64
	R, S                  *big.Int
}

func (cct ChannelCreationTransaction) GetHash(haveRSbeenSet bool) []byte {
	h := sha256.New()
	h.Write(cct.Origin.GetHash())
	h.Write([]byte(cct.ChannelName))
	h.Write([]byte(cct.TokenName))
	for addr, amt := range cct.RecipientsQuantityMap {
		h.Write([]byte(addr))
		h.Write([]byte(strconv.FormatFloat(amt, 'f', -1, 64)))
	}

	// -1 as the precision arg gets the # to 64bit precision intuitively
	h.Write([]byte(strconv.FormatFloat(cct.ExchangeRate, 'f', -1, 64)))

	//Filters the cases where we just want the hash for non-signing purposes
	//(if the transaction hasn't been signed, we shouldn't hash R and S as they don't matter)
	if haveRSbeenSet {
		h.Write(cct.R.Bytes())
		h.Write(cct.S.Bytes())
	}
	return h.Sum(nil)
}

func (cct ChannelCreationTransaction) SetRS(r *big.Int, s *big.Int) SignableTransaction {
	cct.R = r
	cct.S = s
	return cct
}

func (cct ChannelCreationTransaction) GetRS() (*big.Int, *big.Int) {
	return cct.R, cct.S
}

func (cct ChannelCreationTransaction) GetOrigin() OriginInfo {
	return cct.Origin
}

func (cct ChannelCreationTransaction) ToString() string {
	return cct.ChannelName
}
