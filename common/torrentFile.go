package common

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strconv"
)

type LayerFileMetadata struct {
	fileUrl string //used specifically when receiving layers from other sources (without having an entire file yet)
	Begin   int64
	Size    int64
}

func (lm LayerFileMetadata) GetUrl() string {
	return lm.fileUrl
}

type TorrentFile struct {
	Name          string
	LayerByteSize int64
	TotalByteSize int64
	TotalHash     string
	//Total hash is dependent on the name, the segment size, the raw data, and nothing else
	//-the other aspects of a torrent, such as the hashed duplicates and ordering of segments should be LOCAL
	//(we already hash the raw data in order)

	LayerHashKeys []string //fine to expose publicly to say what layers the torrent has

	layerHashMaps map[string]LayerFileMetadata //don't expose; reveals entire torrent structure
	url           string
}

var kilobyte = 1000
var megabyte = 1000000

//GetLayerHashMap exposes the hashes and file metadata associated with a torrent file's layers, but without
//exposing the fields to manipulation, or exposure by mux when listing the TorrentFile as a web response
func (torr TorrentFile) GetLayerHashMap() map[string]LayerFileMetadata {
	return torr.layerHashMaps
}

func (torr TorrentFile) GetUrl() string { //just don't want mux to automatically disclose the url (like in mux)
	return torr.url
}

func (torr TorrentFile) ToString() string {
	a := "Torrent segment size: " + strconv.FormatInt(int64(torr.LayerByteSize), 10) + "\n"
	for _, v := range torr.LayerHashKeys {
		a += v + "\n"
	}
	return a
}

func Int64ToByteArr(i int64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(i))
	return b
}

func (torr TorrentFile) GetRawBytes() []byte {
	ret := make([]byte, 0)
	ret = append(ret, []byte(torr.Name)...)
	ret = append(ret, Int64ToByteArr(torr.LayerByteSize)...)

	for _, v := range torr.LayerHashKeys {
		ret = append(ret, []byte(v)...) //write the hashed key

		//DON'T write the fileurl

		ret = append(ret, Int64ToByteArr(torr.layerHashMaps[v].Begin)...)

		ret = append(ret, Int64ToByteArr(torr.layerHashMaps[v].Size)...)
	}
	return ret
}

func MakeTorrentFileFromFile(layerByteSize int64, url string, name string) (TorrentFile, error) {
	file, err := os.Open(url)
	defer file.Close()
	if err != nil {
		return TorrentFile{}, err
	}

	torr := TorrentFile{name, layerByteSize, 0, "", make([]string, 0),
		make(map[string]LayerFileMetadata, 0), url}
	readbytes := layerByteSize
	total := sha256.New()
	total.Write([]byte(name))
	total.Write(Int64ToByteArr(layerByteSize))
	for offset := int64(0); readbytes == layerByteSize; {
		buffer := make([]byte, layerByteSize)
		temp, _ := file.ReadAt(buffer, offset)
		readbytes = int64(temp)
		if err != nil {
			return TorrentFile{}, err
		}
		torr.appendNewSegment(buffer[0:readbytes], offset, readbytes)
		total.Write(buffer[0:readbytes])
		offset += readbytes
		torr.TotalByteSize += readbytes
		//fmt.Println("Read " + strconv.FormatInt(off / 1000, 10) + " kilobytes so far")
	}
	torr.TotalHash = hex.EncodeToString(total.Sum(nil))
	return torr, nil
}

func MakeTorrentFromBytes(layerByteSize int64, data []byte, name string) (TorrentFile, error) {

	torr := TorrentFile{name, layerByteSize, 0, "", make([]string, 0),
		make(map[string]LayerFileMetadata, 0), ""}

	var offset int64
	total := sha256.New()
	total.Write([]byte(name))
	total.Write(Int64ToByteArr(layerByteSize))
	for offset = 0; offset+layerByteSize < int64(len(data)); {
		segment := data[offset : offset+layerByteSize]
		torr.appendNewSegment(segment, offset, layerByteSize)
		total.Write(segment)
		offset += layerByteSize
	}
	torr.appendNewSegment(data[offset:], offset, int64(len(data))-offset)
	total.Write(data[offset:])
	torr.TotalHash = hex.EncodeToString(total.Sum(nil))
	return torr, nil
}

func (torr1 TorrentFile) Equals(torr2 TorrentFile) bool {
	h1 := torr1.GetRawBytes()
	h2 := torr2.GetRawBytes()

	for i, v := range h1 {
		if v != h2[i] {
			return false
		}
	}
	return true
}

func (torr TorrentFile) Validate() (bool, error) {
	file, err := os.Open(torr.url)
	defer file.Close()
	if err != nil {
		return false, err
	}

	for hash, raw := range torr.layerHashMaps {
		bytes := make([]byte, raw.Size)
		read, _ := file.ReadAt(bytes, raw.Begin)
		if hex.EncodeToString(hashSegment(bytes)) != hash || int64(read) != raw.Size {
			return false, nil
		}
	}
	return true, nil
}

//appendNewSegment adds new raw data to the torrentfile by hashing it and storing it in a map.
func (file *TorrentFile) appendNewSegment(segData []byte, min int64, max int64) {
	hexHashed := hex.EncodeToString(hashSegment(segData))      //hash the data
	file.LayerHashKeys = append(file.LayerHashKeys, hexHashed) //record the hash in order

	if _, ok := file.layerHashMaps[hexHashed]; ok { //we've generated this hash before
		fmt.Println("Duplicate hash")
	} else {
		file.layerHashMaps[hexHashed] = LayerFileMetadata{file.url, min, max}
	}
}

func hashSegment(seg []byte) []byte {
	h := sha256.New()
	h.Write(seg)
	return h.Sum(nil)
}

func AppendLayerDataToFile(layerId string, data []byte) LayerFileMetadata {
	h := sha256.New()
	h.Write(data)
	if hex.EncodeToString(h.Sum(nil)) != layerId {
		fmt.Println("Id " + layerId + " hash doesn't match data hash!")
		return LayerFileMetadata{}
	}

	f, err := os.OpenFile("layers.data", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Println("Couldn't open layers file for writing new layer")
		return LayerFileMetadata{}
	}
	defer f.Close()

	info, _ := f.Stat()
	offset := info.Size()

	var keyData = []byte(layerId + ":")
	offset += int64(len(keyData)) //how much we are actually offset before we write raw data
	writeSize := len(data)

	if _, err = f.WriteString(string(keyData) + string(data) + "\n"); err != nil {
		panic(err)
	}
	return LayerFileMetadata{"layers.data", offset, int64(writeSize)}
}
