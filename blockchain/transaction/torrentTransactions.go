package transaction

import (
	"github.com/denverquane/GoBlockShare/files"
)

type TorrentTransaction interface {
	GetType() string
	GetRawBytes() []byte
	ToString() string
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

func boolToByte(v bool) byte {
	if v {
		return 't'
	} else {
		return 'f'
	}
}

type PublishTorrentTrans struct {
	Torrent files.TorrentFile
	Type 	string
}
func (tt PublishTorrentTrans) GetType() string {
	return tt.Type
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
	Type 			string
}
func (lt SharedLayerTrans) GetType() string {
	return lt.Type
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
	Type 		string
}
func (rt LayerRepTrans) GetType() string {
	return rt.Type
}
func (rt LayerRepTrans) GetRawBytes() []byte {
	return []byte(rt.TxID + string(boolToByte(rt.WasLayerValid)))
}
func (rt LayerRepTrans) ToString() string {
	return "Gave " + string(boolToByte(rt.WasLayerValid)) + " rep for layer shared in TX: " + rt.TxID
}

type TorrentRepTrans struct {
	TxID		string //the original transaction when the layer was shared with "me"
	RepMessage	RepMessage //TODO should probably be a more complex message later (esp for torrents
	Type 		string
}
func (rt TorrentRepTrans) GetType() string {
	return rt.Type
}
func (rt TorrentRepTrans) GetRawBytes() []byte {
	return []byte(rt.TxID + string(rt.RepMessage.toBytes()))
}
func (rt TorrentRepTrans) ToString() string {
	return "Gave " + string(rt.RepMessage.toBytes()) + " rep for torrent shared in TX: " + rt.TxID
}
