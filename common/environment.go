package common

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

const BlockchainHostDefault = "http://localhost"
const BlockchainPortDefault = "5000"
const TorrentPortDefault = "8000"

type EnvVars struct {
	BlockchainHost string
	BlockchainPort string
	TorrentPort    string
}

func LoadEnvFromFile(topDir string) EnvVars {
	err := godotenv.Load(".env")
	if err != nil {
		err = godotenv.Load(topDir + "/.env")
		if err != nil {
			log.Println("Can't load env file; checking system env vars")
		} else {
			log.Println("Using " + topDir + " local env file")
		}
	} else {
		log.Println("Using " + topDir + " local env file")
	}
	return getVarsFromEnvOrDefaults()
}

func getVarsFromEnvOrDefaults() EnvVars {
	vars := EnvVars{}
	vars.BlockchainHost = os.Getenv("BLOCKCHAIN_HOST")
	vars.BlockchainPort = os.Getenv("BLOCKCHAIN_PORT")
	vars.TorrentPort = os.Getenv("TORRENT_PORT")

	if vars.BlockchainHost == "" {
		vars.BlockchainHost = BlockchainHostDefault
		log.Println("Couldn't find BLOCKCHAIN_HOST; using default of \"" + BlockchainHostDefault + "\"")

	}
	if vars.BlockchainPort == "" {
		vars.BlockchainPort = BlockchainPortDefault
		log.Println("Couldn't find BLOCKCHAIN_PORT; using default of " + BlockchainPortDefault)
	}
	if vars.TorrentPort == "" {
		vars.TorrentPort = TorrentPortDefault
		log.Println("Couldn't find TORRENT_PORT; using default of " + TorrentPortDefault)
	}
	return vars
}
