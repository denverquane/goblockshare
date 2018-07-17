package transaction

import (
	"crypto/sha256"
	"math/big"
	"github.com/denverquane/GoBlockShare/files"
)

type TorrentTransaction struct {
	Origin   	OriginInfo // needed to say who I am (WITHIN the transaction)
	Transaction	Transaction
	R, S     *big.Int // signature of the transaction, should be separate from the actual "message" components
}

func (st TorrentTransaction) SetRS(r *big.Int, s *big.Int) SignableTransaction {
	st.R = r
	st.S = s
	return st
}

func (st TorrentTransaction) GetRS() (*big.Int, *big.Int) {
	return st.R, st.S
}

func (st TorrentTransaction) GetHash(haveRSbeenSet bool) []byte {
	h := sha256.New()
	h.Write(st.Origin.GetRawBytes())
	h.Write(st.Transaction.GetRawBytes())

	//Filters the cases where we just want the hash for non-signing purposes
	//(if the transaction hasn't been signed, we shouldn't hash R and S as they don't matter)
	if haveRSbeenSet {
		h.Write(st.R.Bytes())
		h.Write(st.S.Bytes())
	}
	return h.Sum(nil)
}

func (st TorrentTransaction) GetOrigin() OriginInfo {
	return st.Origin
}

func (st TorrentTransaction) GetType() string {
	return st.Transaction.GetType()
}

func (st TorrentTransaction) ToString() string {
	return st.Origin.ToString() + "\"txref\":[],\n" +
		string(st.Transaction.GetRawBytes()) + "\",\n\"r\":" + st.R.String() + ",\n\"s\":" +
		st.S.String() + "\n}\n"
}



type PublishTorrentTrans struct {
	Torrent files.TorrentFile
}
func (tt PublishTorrentTrans) GetType() string {
	return "PUBLISH_TORRENT"
}
func (tt PublishTorrentTrans) GetRawBytes() []byte {
	return tt.Torrent.GetRawBytes()
}
func (tt PublishTorrentTrans) ToString() string {
	return tt.Torrent.ToString()
}


type SharedLayerTrans struct {
	SharedLayerHash []byte
	Recipient		Base64Address
}
func (lt SharedLayerTrans) GetType() string {
	return "SHARED_LAYER"
}
func (lt SharedLayerTrans) GetRawBytes() []byte {
	return []byte(string(lt.SharedLayerHash) + string(lt.Recipient))
}
func (lt SharedLayerTrans) ToString() string {
	return "Shared " + string(lt.SharedLayerHash) + " with " + string(lt.Recipient)
}


type LayerRepTrans struct {
	TxID		string //the original transaction when the layer was shared with "me"
	WasLayerValid	bool //TODO should probably be a more complex message later
}
func (rt LayerRepTrans) GetType() string {
	return "LAYER_REP"
}
func (rt LayerRepTrans) GetRawBytes() []byte {
	return []byte(rt.TxID + string(boolToByte(rt.WasLayerValid)))
}
func (rt LayerRepTrans) ToString() string {
	return "Gave " + string(boolToByte(rt.WasLayerValid)) + " rep for layer shared in TX: " + rt.TxID
}

func boolToByte(v bool) byte {
	if v {
		return 't'
	} else {
		return 'f'
	}
}

type RepMessage struct {
	WasValid	 bool
	HighQuality  bool
	AccurateName bool
}

func (rm RepMessage) toBytes() []byte {
	b := make([]byte, 3)
	b[0] = boolToByte(rm.WasValid)
	b[1] = boolToByte(rm.HighQuality)
	b[2] = boolToByte(rm.AccurateName)
	return b
}

type TorrentRepTrans struct {
	TxID		string //the original transaction when the layer was shared with "me"
	RepMessage	RepMessage //TODO should probably be a more complex message later (esp for torrents
}
func (rt TorrentRepTrans) GetType() string {
	return "TORRENT_REP"
}


func (rt TorrentRepTrans) GetRawBytes() []byte {
	return []byte(rt.TxID + string(rt.RepMessage.toBytes()))
}
func (rt TorrentRepTrans) ToString() string {
	return "Gave " + string(rt.RepMessage.toBytes()) + " rep for torrent shared in TX: " + rt.TxID
}
