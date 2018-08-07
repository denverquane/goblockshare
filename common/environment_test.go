package common

import "testing"

func TestLoadEnvFromFile(t *testing.T) {
	env := LoadEnvFromFile("")

	if env.BlockchainHost == "" || env.BlockchainPort == "" || env.TorrentPort == "" {
		t.Fail()
	}
}
