package files

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestMakeTorrentFile(t *testing.T) {
	now := time.Now()
	torr, err := MakeTorrentFileFromFile(kilobyte, "../README.md")
	after := time.Now()
	fmt.Println("Took " + strconv.FormatFloat(after.Sub(now).Seconds(), 'f', 2, 64) + " seconds to make torrent")

	if err != nil {
		fmt.Print(err)
		return
	}

	for segmentSize := 2;  segmentSize < 10; segmentSize++ {
		torr, _ := MakeTorrentFileFromFile(segmentSize, "../README.md")
		val := torr.GetDuplicatesTotal() * (segmentSize * segmentSize * segmentSize)
		fmt.Println("Size " + string(segmentSize + int('0')) + " val: " + strconv.FormatInt(int64(val), 10))
	}

	fmt.Println(torr.Validate())
}

func TestAreSameTorrentBytes(t *testing.T) {
	sampleA, _ := MakeTorrentFromBytes(2, []byte("asdfasdfhijk"))
	sampleB, _ := MakeTorrentFromBytes(2, []byte("asdfasdfhijk"))

	if !sampleA.Equals(sampleB){
		fmt.Println("Not equals!")
		fmt.Println(sampleA)
		fmt.Println(sampleB)
		t.Fail()
	}
	fmt.Println(sampleA.GetDuplicatesTotal())
	sampleC, _ := MakeTorrentFromBytes(2, []byte("asdfasf"))

	if sampleA.Equals(sampleC){
		fmt.Println("Equals shorter!")
		t.Fail()
	}

	sampleD, _ := MakeTorrentFromBytes(3, []byte("asdfasdf"))

	if sampleA.Equals(sampleD){
		fmt.Println("Equals diff seg length!")
		t.Fail()
	}
}
