package wallet

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"github.com/denverquane/GoBlockShare/blockchain/transaction"
	//"crypto/sha256"
	//"crypto/x509"
	//"encoding/hex"
)

//TODO when making a channel, we need to be aware of which users actually have the token that we sent them, and only
// bother parsing/processing those senders' messages.

//TODO any issuing of new tokens needs to be a validated transaction first, to make sure people don't make numerous
//tokens all under the same name

type ChannelHandshakeStage int

const (
	Uninitialized ChannelHandshakeStage = iota
	ReceivedTokenAndChannelPub
	SentMyEncryptedPub
	ReceivedEncryptedPrivAndAddr
)

// This struct is for storing the information required to interact with a channel
type ChannelRecord struct {
	status            ChannelHandshakeStage

	channelCreatorAddress transaction.Base64Address //where we get the channel info from

	channelPublic  rsa.PublicKey //how to encrypt messages to post onto the channel
	channelPrivate rsa.PrivateKey //how we decrypt channel messages TODO Encrypt and save locally
	channelAddress transaction.Base64Address //where we post channel messages

	myPublic  rsa.PublicKey //how the channel creator sends the info to me
	myPrivate rsa.PrivateKey //how I decrypt info the channel creator sends me TODO Encrypt and save locally
	myAddress transaction.PersonalAddress //where I expect to receive channel info
}

func GenerateNewChannelRecord(tokenName string, theiraddress transaction.Base64Address,
	myaddress transaction.PersonalAddress) ChannelRecord {

	record := ChannelRecord{status:Uninitialized, channelCreatorAddress:theiraddress, myAddress: myaddress}

	private, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		fmt.Println("Problem generating RSA key!")
		return record
	}

	record.myPrivate = *private
	record.myPublic = record.myPrivate.PublicKey
	record.myAddress = myaddress

	fmt.Println("Generated a new channel record for: " + tokenName)

	return record
}

//func (cr ChannelRecord) MakeSendMyPublicTransaction() transaction.SignedTransaction {
//	if cr.status != ReceivedTokenAndChannelPub {
//		fmt.Println("Shouldn't send my key until we have a token and the channel's public key...")
//	}
//	pubKeyBytes, _ := x509.MarshalPKIXPublicKey(cr.myPublic)
//	//format my public key
//	encrypted, _ := rsa.EncryptOAEP(sha256.New(), rand.Reader, &cr.channelPublic, pubKeyBytes, []byte("key"))
//	//encrypt my public key using the channel's public key
//	origin := transaction.OriginInfo{*cr.myAddress.PublicKey.X, *cr.myAddress.PublicKey.Y, cr.myAddress.Address}
//	signed := transaction.SignedTransaction{origin, cr.channelCreatorAddress, 0.0, "", hex.EncodeToString(encrypted), nil, nil}
//	signed = signed.SignMessage(&cr.myAddress.PrivateKey)
//	return signed
//}
