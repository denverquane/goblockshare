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
	h.Write(st.Origin.GetHash())
	h.Write(st.Transaction.GetHash())

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
		st.Transaction.ToString() + "\",\n\"r\":" + st.R.String() + ",\n\"s\":" +
		st.S.String() + "\n}\n"
}



type PublishTorrentTrans struct {
	Torrent files.TorrentFile
}
func (tt PublishTorrentTrans) GetType() string {
	return "PUBLISH_TORRENT"
}
func (tt PublishTorrentTrans) GetHash() []byte {
	return tt.Torrent.GetHash()
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
func (lt SharedLayerTrans) GetHash() []byte {
	h := sha256.New()
	h.Write(lt.SharedLayerHash)
	h.Write([]byte(lt.Recipient))
	return h.Sum(nil)
}
func (lt SharedLayerTrans) ToString() string {
	return "Shared " + string(lt.SharedLayerHash) + " with " + string(lt.Recipient)
}


type LayerRepTrans struct {
	TxID		string //the original transaction when the layer was shared with "me"
	RepMessage	string //TODO should probably be a more complex message later
}
func (rt LayerRepTrans) GetType() string {
	return "LAYER_REP"
}
func (rt LayerRepTrans) GetHash() []byte {
	h := sha256.New()
	h.Write([]byte(rt.TxID))
	h.Write([]byte(rt.RepMessage))
	return h.Sum(nil)
}
func (rt LayerRepTrans) ToString() string {
	return "Gave " + rt.RepMessage + " rep for layer shared in TX: " + rt.TxID
}

type TorrentRepTrans struct {
	TxID		string //the original transaction when the layer was shared with "me"
	RepMessage	string //TODO should probably be a more complex message later (esp for torrents
}
func (rt TorrentRepTrans) GetType() string {
	return "TORRENT_REP"
}
func (rt TorrentRepTrans) GetHash() []byte {
	h := sha256.New()
	h.Write([]byte(rt.TxID))
	h.Write([]byte(rt.RepMessage))
	return h.Sum(nil)
}
func (rt TorrentRepTrans) ToString() string {
	return "Gave " + rt.RepMessage + " rep for torrent shared in TX: " + rt.TxID
}
