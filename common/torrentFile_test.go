package common

import (
	"fmt"
	"os"
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

	fmt.Println(torr.Validate())
}

func TestFileVsBytesTorrentFile(t *testing.T) {
	torr, err := MakeTorrentFileFromFile(int64(kilobyte), "../README.md", "readme.md")
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	file, err := os.Open("../README.md")
	info, err := file.Stat()
	bytes := make([]byte, info.Size())
	file.Read(bytes)

	torr2, err := MakeTorrentFromBytes(int64(kilobyte), bytes, "readme.md")

	if !torr.Equals(torr2) {
		t.Fail()
	}
	fmt.Println(torr.layerHashMaps)
	fmt.Println(torr2.layerHashMaps)
}

func TestTorrentFile_Equals(t *testing.T) {
	sampleA, _ := MakeTorrentFromBytes(2, []byte("asdfasdfhijk"), "readme.md")
	sampleB, _ := MakeTorrentFromBytes(2, []byte("asdfasdfhijk"), "readme.md")

	if !sampleA.Equals(sampleB) {
		t.Fail()
	}

	sampleA.Name = "readme"

	if sampleA.Equals(sampleB) {
		t.Fail()
	}
}

func TestAreSameTorrentBytes(t *testing.T) {
	sampleA, _ := MakeTorrentFromBytes(2, []byte("asdfasdfhijk"), "readme.md")
	sampleB, _ := MakeTorrentFromBytes(2, []byte("asdfasdfhijk"), "readme.md")

	if !sampleA.Equals(sampleB) {
		fmt.Println("Not equals!")
		fmt.Println(sampleA)
		fmt.Println(sampleB)
		t.Fail()
	}

	sampleC, _ := MakeTorrentFromBytes(2, []byte("asdfasf"), "readme.md")

	if sampleA.Equals(sampleC) {
		fmt.Println("Equals shorter!")
		t.Fail()
	}

	sampleD, _ := MakeTorrentFromBytes(3, []byte("asdfasdf"), "readme.md")

	if sampleA.Equals(sampleD) {
		fmt.Println("Equals diff seg length!")
		t.Fail()
	}
}
