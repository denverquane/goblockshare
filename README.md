# GoBlockShare
[![GoDoc](https://godoc.org/github.com/denverquane/GoBlockShare?status.png)](https://godoc.org/github.com/denverquane/GoBlockShare)
[![codecov](https://codecov.io/gh/denverquane/GoBlockShare/branch/master/graph/badge.svg)](https://codecov.io/gh/denverquane/GoBlockShare)
[![Build Status](https://travis-ci.org/denverquane/GoBlockShare.svg?branch=master)](https://travis-ci.org/denverquane/GoBlockShare)
[![Go Report Card](https://goreportcard.com/badge/github.com/denverquane/GoBlockShare)](https://goreportcard.com/report/github.com/denverquane/GoBlockShare)

#### DISCLAIMER
This project is under constant development, and likely will not be fully (or even partially) functioning until this 
disclaimer has been removed. Clone at your own peril!

## Summary

This application seeks to demonstrate how decentralized file sharing protocols like BitTorrent can be adapted for use with
blockchain technology. This would allow users to distribute files and content in a method similar to the Bittorrent protocol,
but with the motivations and incentives lended by blockchain systems like Bitcoin and Ethereum. Users that publish
files will be rewarded -on the blockchain- by users that validate the file's contents and validity, while users that
cast their "vote" will be rewarded with higher priority and bandwidth for future file requests.

The use of a blockchain solution prevents any single entity from deciding who is a reputable file source, while also providing
a "monetary" incentive for posting valid and high quality content. Similarly, users that validate and submit feedback on
files and content will be rewarded using the blockchain, by allowing all parties to search the chain history to see a user's
reputation, feedback, etc. for files they have downloaded or published. 

This project aims to solve a crucial problem with the BitTorrent protocol, namely the lack of incentivization for "seeders", or those
that actively help distribute content to other users, as opposed to "leechers", who solely download and do not redistribute content.
While there are many BitTorrent communities that restrict access to users with a certain ratio of uploads/downloads to 
navigate this issue, these communities are often hard to find, restrictive to newcomers, subject to "single-party" authoritation, and (rarely) offer any incentivization for content providers.

Blockchain technology offers the ability to search a decentralized and public database to easily determine which users have uploaded
valid and quality content, users that contribute feedback or redistribute content, and which users are incentivizing content
providers via direct or indirect compensation. At the same time, malicious users who submit false feedback or verification
results, solely "leech" content, or who do not contribute *any* feedback or direct compensation 
to content creators, will be quickly vetted or outnumbered by the users who *do* contribute to the community.

This project aims to make many of these steps automatic and intuitive (such as restricting content sharing with unreputable 
users, or verifying content and submitting feedback on download), while providing flexible behavior for users or uploaders
who desire more granularity and control.  
 

## Installation
A valid Go installation is required to be able to install and run this project: https://golang.org/doc/install

Assuming your Go installation is configured correctly, you can then clone this project using 
`git clone https://github.com/denverquane/goblockshare`, and run `go build` on either the `torrentshare/main/torrentServer.go` file (which acts as a server for the actual torrent filesharing), or `blockchain/main/blockchainServer.go`, which functions as the server that hosts the blockchain and processes transactions.

These two servers can be ran concurrently to allow full functionality with file sharing, torrent file broadcasting on the blockchain, checking reputation of users to permit or restrict access, etc.

Happy Sharing!

## Goals:
- [ ] Proof of Work for posting messages/blocks (prevent spam/abuse)
  - [X] Basic difficulty/cryptographic proof validation
  - [ ] Scaling difficulty of blocks as the ~~chain~~ userbase grows
  - [ ] Rewards for propagating files/layers without uploading new content
    - Reputation as a "currency"
    - See Bitcoin whitepaper for inspiration: [Bitcoin.org](https://bitcoin.org/bitcoin.pdf)
- [ ] Node discovery
  - [ ] Ability to run app as a node registry/lookup
  - [ ] Active central registry for node lookup
  - [ ] Automatically propagate chain changes to other nodes
  - [ ] Consensus algo. for chain conflicts (explore merging non-conflicting message types/re-generation of transaction)
  - [X] Dockerize app for simplified multi-node testing on a single physical machine
- [X] Author/poster validation (login validation)
  - [X] Basic authentication
  - [X] ECDSA Public/Private key generation
- [ ] "Torrent" creation
  - [X] Slice files into segments or "layers"
  - [ ] Signage of overall file, and other metadata by original poster
  - [ ] Broadcast local segments available for sharing
- [X] Basic JS Frontend for viewing the blockchain in realtime
- [ ] JS Frontend for posting, deleting, editing, etc. messages and transactions
- [ ] Electron integration for running the web app as a native desktop app

## Project Structure:
This app is comprised of two distinct parts as of 7/1/18.
These parts serve as the frontend and backend services for the overarching application.

- #### Go Backend
This repository contains all the [GoLang](https://golang.org/) source files, which run the backend
application that handles blockchain operations, including peer discovery/chain propagation, author and user validation,
torrent file creation, and rudimentary ["Proof of Work"](https://en.wikipedia.org/wiki/Proof-of-work_system) calculations to ensure users don't
overwhelm blockchain peers with transactions.

- #### TypeScript/React Frontend
See [ReactBlockShare](https://github.com/denverquane/ReactBlockShare)
