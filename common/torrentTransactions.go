package common

import (

)

type TorrentTransaction interface {
	GetRawBytes() []byte
	ToString() string
}

//RepMessage represents a message that conveys the quality of a Torrent, in terms of the accuracy of the name, the
//validity of the file itself (does it form a cohesive file that runs/opens as expected?), and if the content itself is
//high quality
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

//PublishTorrentTrans represents a transaction that shows a node/user posting a new torrent onto the blockchain
//This allows other users and nodes to discover what layers comprise what type of file or content, so that the individual
//layers can be directly/"privately" requested from nodes (users don't request layers on the blockchain)
type PublishTorrentTrans struct {
	Torrent TorrentFile
}
func (tt PublishTorrentTrans) GetRawBytes() []byte {
	return tt.Torrent.GetRawBytes()
}
func (tt PublishTorrentTrans) ToString() string {
	return tt.Torrent.ToString()
}
//super cool compile-time interface check!
var _ TorrentTransaction = PublishTorrentTrans{}

//SharedLayerTrans represents a transaction that shows a node sharing a torrent layer with a particular address
//This can be used by other nodes to see if a node has indicated sharing with an address, and whether or not that address
//left feedback on the share, if that address actually received the layer, etc.
type SharedLayerTrans struct {
	SharedLayerHash []byte
	Recipient		string
}
func (lt SharedLayerTrans) GetRawBytes() []byte {
	return []byte(string(lt.SharedLayerHash) + string(lt.Recipient))
}
func (lt SharedLayerTrans) ToString() string {
	return "Shared " + string(lt.SharedLayerHash) + " with " + string(lt.Recipient)
}

var _ TorrentTransaction = SharedLayerTrans{}

//LayerRepTrans represents an address leaving feedback on a layer that was shared with it. This can be used to determine
//if a particular node/address is reputable when sharing layers that are valid and hash correctly,
// or even if they shared the layer at all
type LayerRepTrans struct {
	TxID		string //the original transaction when the layer was shared with "me"
	WasLayerReceived bool
	WasLayerValid	bool
}
func (rt LayerRepTrans) GetRawBytes() []byte {
	return []byte(rt.TxID + string(boolToByte(rt.WasLayerReceived)) + string(boolToByte(rt.WasLayerValid)))
}
func (rt LayerRepTrans) ToString() string {
	return "Received?: " + string(boolToByte(rt.WasLayerReceived)) + " and valid?: " +
		string(boolToByte(rt.WasLayerValid)) + " for layer shared in TX: " + rt.TxID
}

var _ TorrentTransaction = LayerRepTrans{}

//TorrentRepTrans represents a transaction for feedback on a Torrent. See RepMessage for more details
type TorrentRepTrans struct {
	TxID		string //the original transaction when the layer was shared with "me"
	RepMessage	RepMessage
}
func (rt TorrentRepTrans) GetRawBytes() []byte {
	return []byte(rt.TxID + string(rt.RepMessage.toBytes()))
}
func (rt TorrentRepTrans) ToString() string {
	return "Gave " + string(rt.RepMessage.toBytes()) + " rep for torrent shared in TX: " + rt.TxID
}

var _ TorrentTransaction = TorrentRepTrans{}