# GoBlockShare
[![GoDoc](https://godoc.org/github.com/denverquane/GoBlockShare?status.png)](https://godoc.org/github.com/denverquane/GoBlockShare)
[![Build Status](https://travis-ci.org/denverquane/GoBlockShare.svg?branch=master)](https://travis-ci.org/denverquane/GoBlockShare)

This app seeks to demonstrate using blockchain tech. for lightweight transfers of a cryptocurrency, but with secure chat
functionality built on top.

This is mainly a personal exploration into the blockchain realm, but may also serve as a convenient way to demonstrate how
cryptocurrency can be a viable payment method when used in conjunction with something like a chat application, for example.
This would allow users to communicate in realtime, while also being able to "inject" block and transaction information
during a discussion.


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