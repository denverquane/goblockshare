package transaction

import (
	"crypto/sha256"
	"math/big"
	"strconv"
)

type SignedTransaction struct {
	Origin   OriginInfo
	DestAddr Base64Address
	Quantity float64
	Payload  string
	R, S     *big.Int
}

func (st SignedTransaction) SetRS(r *big.Int, s *big.Int) SignableTransaction {
	st.R = r
	st.S = s
	return st
}

func (st SignedTransaction) GetRS() (*big.Int, *big.Int) {
	return st.R, st.S
}

func (st SignedTransaction) GetHash(haveRSbeenSet bool) []byte {
	h := sha256.New()
	h.Write(st.Origin.GetHash())
	h.Write([]byte(st.DestAddr))
	// -1 as the precision arg gets the # to 64bit precision intuitively
	h.Write([]byte(strconv.FormatFloat(st.Quantity, 'f', -1, 64)))

	//Filters the cases where we just want the hash for non-signing purposes
	//(if the transaction hasn't been signed, we shouldn't hash R and S as they don't matter)
	if haveRSbeenSet {
		h.Write(st.R.Bytes())
		h.Write(st.S.Bytes())
	}
	h.Write([]byte(st.Payload))
	return h.Sum(nil)
}

func (st SignedTransaction) GetOrigin() OriginInfo {
	return st.Origin
}

func (st SignedTransaction) ToString() string {
	return st.Origin.ToString() + "\"txref\":[],\n\"quantity\":" +
		strconv.FormatFloat(st.Quantity, 'f', -1, 64) + ",\n\"currency\":\"" +
		"\",\n\"payload\":\"" + st.Payload + "\",\n\"r\":" + st.R.String() + ",\n\"s\":" +
		st.S.String() + ",\n\"destAddr\":\"" +
		string(st.DestAddr) + "\"\n}\n"
}
