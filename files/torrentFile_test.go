package files

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestMakeTorrentFile(t *testing.T) {
	now := time.Now()
	torr, err := MakeTorrentFileFromFile(int64(kilobyte), "../README.md", "readme.md")
	after := time.Now()
	fmt.Println("Took " + strconv.FormatFloat(after.Sub(now).Seconds(), 'f', 2, 64) + " seconds to make torrent")

	if err != nil {
		fmt.Print(err)
		return
	}

	for segmentSize := int64(2); segmentSize < 10; segmentSize++ {
		torr, _ := MakeTorrentFileFromFile(segmentSize, "../README.md", "readme.md")
		val := torr.GetDuplicatesTotal()
		fmt.Println("Size " + string(segmentSize+int64('0')) + " duplicate segments: " + strconv.FormatInt(int64(val), 10))
	}

	fmt.Println(torr.Validate())
}

//func TestAreSameTorrentBytes(t *testing.T) {
//	sampleA, _ := MakeTorrentFromBytes(2, []byte("asdfasdfhijk"), "readme.md")
//	sampleB, _ := MakeTorrentFromBytes(2, []byte("asdfasdfhijk"), "readme.md")
//
//	if !sampleA.Equals(sampleB) {
//		fmt.Println("Not equals!")
//		fmt.Println(sampleA)
//		fmt.Println(sampleB)
//		t.Fail()
//	}
//	fmt.Println(sampleA.GetDuplicatesTotal())
//	sampleC, _ := MakeTorrentFromBytes(2, []byte("asdfasf"), "readme.md")
//
//	if sampleA.Equals(sampleC) {
//		fmt.Println("Equals shorter!")
//		t.Fail()
//	}
//
//	sampleD, _ := MakeTorrentFromBytes(3, []byte("asdfasdf"), "readme.md")
//
//	if sampleA.Equals(sampleD) {
//		fmt.Println("Equals diff seg length!")
//		t.Fail()
//	}
//}
