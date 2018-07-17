package files

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"strconv"
	"encoding/binary"
)

type LayerFileMetadata struct {
	Begin int64
	Offset int64
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

	layerHashMaps  map[string]LayerFileMetadata //don't expose; reveals entire torrent structure
	duplicatesMaps map[string]int    //don't reveal; should only be used internally/locally
	url	string
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
	ret = append(ret, Int64ToByteArr(torr.LayerByteSize)...)

	for _, v := range torr.LayerHashKeys {
		ret = append(ret, []byte(v)...)             //write the hashed key

		ret = append(ret, Int64ToByteArr(torr.layerHashMaps[v].Begin)...)

		ret = append(ret, Int64ToByteArr(torr.layerHashMaps[v].Offset)...)
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
		make(map[string]LayerFileMetadata, 0), make(map[string]int), url}
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

//func MakeTorrentFromBytes(segByteSize int64, data []byte, name string) (TorrentFile, error) {
//	if segByteSize > int64(len(data)) {
//		return TorrentFile{}, errors.New("Segment too long")
//	}
//
//	torr := TorrentFile{name, segByteSize, 0,"", make([]string, 0), make(map[string][]byte, 0), make(map[string]int),}
//
//	var offset int64
//	total := sha256.New()
//	total.Write([]byte(name))
//	//TODO write segbytesize (dont allow torrents to change segmentation size)
//	for offset = 0; offset+segByteSize < int64(len(data)); {
//		segment := data[offset : offset+segByteSize]
//
//		offset += segByteSize
//		torr.appendNewSegment(segment)
//		total.Write(segment)
//	}
//	torr.appendNewSegment(data[offset:])
//	total.Write(data[offset:])
//	torr.TotalHash = hex.EncodeToString(total.Sum(nil))
//	return torr, nil
//}

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
		bytes := make([]byte, raw.Offset)
		read, _ := file.ReadAt(bytes, raw.Begin)
		if hex.EncodeToString(hashSegment(bytes)) != hash || int64(read) != raw.Offset{
			return false, nil
		}
	}
	return true, nil
}

func (file TorrentFile) GetDuplicatesTotal() int {
	total := 0
	for _, v := range file.duplicatesMaps {
		total += v
	}
	return total
}

//appendNewSegment adds new raw data to the torrentfile by hashing it and storing it in a map. If the same hash has
//been computed previously, then the counter is incremented for that particular entry (to allow for analyzing if there
//are common layers of file, and thus save on storage/bandwidth)
func (file *TorrentFile) appendNewSegment(segData []byte, min int64, max int64) {
	hexHashed := hex.EncodeToString(hashSegment(segData))      //hash the data
	file.LayerHashKeys = append(file.LayerHashKeys, hexHashed) //record the hash in order

	if _, ok := file.layerHashMaps[hexHashed]; ok { //we've generated this hash before
		if _, okk := file.duplicatesMaps[hexHashed]; okk { //we've made the entry before
			file.duplicatesMaps[hexHashed]++
		} else {
			file.duplicatesMaps[hexHashed] = 1 //this is the 2nd occurrence (counter starts at 1 for "1st duplicate")
		}
	} else {
		file.layerHashMaps[hexHashed] = LayerFileMetadata{min, max}
	}
}

func hashSegment(seg []byte) []byte {
	h := sha256.New()
	h.Write(seg)
	return h.Sum(nil)
}
