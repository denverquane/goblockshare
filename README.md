# GoBlockShare
[![GoDoc](https://godoc.org/github.com/denverquane/GoBlockShare?status.png)](https://godoc.org/github.com/denverquane/GoBlockShare)
[![Build Status](https://travis-ci.org/denverquane/GoBlockShare.svg?branch=master)](https://travis-ci.org/denverquane/GoBlockShare)

This application seeks to demonstrate how decentralized file sharing protocols like Bittorrent can be adapted for use with blockchain technology. This would allow users to distribute files and content in a method similar to Bittorrent, but with the motivations and incentives lended by blockchain systems like Bitcoin and Ethereum. Users that publish valid files will be rewarded -on the blockchain- proportional to the amount of users that agree with the content's validity, and these users that cast their "vote" will be similarly rewarded for sharing their opinion and contributing to the uploader's reputation. Potentially a minimum number of nodes would have to cast their votes regarding a content's validity before any are rewarded for siding with the majority

Potentially, some uploaders could restrict their content to nodes that pay them first. This would ensure that 1. Uploaders can restrict access to those with financial incentives to request access, 2. Users that pay to access content are encouraged to vote to recoup some of the cost, 3. Users are encouraged to redistribute the content to similarly recoup losses (need to explore this idea further)

## Goals:
- [ ] Proof of Work for posting messages/blocks (prevent spam/abuse)
  - [X] Basic difficulty/cryptographic proof validation
  - [ ] Scaling difficulty of blocks as the ~~chain~~ userbase grows
  - [ ] Rewards for propagating the chain, even if not posting
    - Reputation as a global "currency"
    - See Bitcoin whitepaper for inspiration: [Bitcoin.org](https://bitcoin.org/bitcoin.pdf)
- [ ] Node discovery
  - [ ] Save chain/channel to disk -> especially for channels with low usercounts w/ possibility of perm. loss
  - [ ] Ability to run app as a node registry/lookup
  - [ ] Active central registry for node lookup
  - [ ] Automatically propagate chain changes to other nodes
  - [ ] Consensus algo. for chain conflicts (explore merging non-conflicting message types/re-generation of transaction)
  - [X] Dockerize app for simplified multi-node testing on a single physical machine
- [X] Author/poster validation (login validation)
  - [X] Basic authentication
  - [X] ECDSA Public/Private key generation
- [X] Basic JS Frontend for viewing the blockchain in realtime
- [ ] JS Frontend for posting, deleting, editing, etc. messages and transactions
- [ ] Electron integration for running the web app as a native desktop app

## Project Structure:
This app is comprised of two distinct parts as of 4/23/18.
These parts serve as the frontend and backend services for the overarching application.

- #### GO Backend
In the aptly-named "Go" directory are all the [GoLang](https://golang.org/) source files, which run the backend
application that handles blockchain operations, including peer discovery/chain propagation, author and user validation,
and rudimentary ["Proof of Work"](https://en.wikipedia.org/wiki/Proof-of-work_system) calculations to ensure users don't
overwhelm blockchain peers with transactions.

- #### Ts/React Frontend
See [ReactBlockShare](https://github.com/denverquane/ReactBlockShare)
