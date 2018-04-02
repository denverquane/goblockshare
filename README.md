# GoBlockChat
[![GoDoc](https://godoc.org/github.com/denverquane/GoBlockChat?status.png)](https://godoc.org/github.com/denverquane/GoBlockChat)
[![Build Status](https://travis-ci.org/denverquane/GoBlockChat.svg?branch=master)](https://travis-ci.org/denverquane/GoBlockChat)

This app seeks to demonstrate using blockchain tech. for a chat application like Slack or Discord.

By using the blockchain, no user can delete, modify, or falsify chat records and exchanges without other users
being aware of the change. This not ensures integrity of the chat, but allows for interesting functionality 
regarding "rewinding" or rollback of chat engagements.

## Goals:
- [ ] Proof of Work for posting messages/blocks (prevent spam/abuse)
  - [X] Basic difficulty/cryptographic proof validation
  - [X] Scaling difficulty of blocks as the ~~chain~~userbase grows
  - [ ] Rewards for propagating the chain (?)
- [ ] Node discovery
  - [ ] Ability to run app as a node registry/lookup
  - [ ] Active central registry for node lookup
  - [ ] Automatically propagate chain changes to other nodes
  - [ ] Consensus algo. for chain conflicts (explore merging non-conflicting message types/re-generation of transaction)
- [ ] Author/poster validation (login validation)
  - [X] Basic authentication
  - [ ] Secure authentication (explore security/abuse vulnerabilities)
  - [ ] Permission tiers?
- [ ] Ensure users are running the same program version
  - [X] Hash source code to ensure no modifications to versions
  - [ ] Only accept/transmit to nodes with the same version
- [X] Basic JS Frontend for viewing the blockchain in realtime
- [ ] JS Frontend for posting, deleting, editing, etc. messages and transactions
- [ ] GO app to interact with the chain, without the Webapp (?)

## Project Structure:
This app is comprised of two distinct parts as of 3/30/18.
These parts serve as the frontend and backend services for the overarching application.

- #### GO Backend
In the aptly-named "Go" directory are all the [GoLang](https://golang.org/) source files, which run the backend application that handles blockchain operations, including peer discovery/chain propagation, author and user validation, and rudimentary ["Proof of Work"](https://en.wikipedia.org/wiki/Proof-of-work_system) calculations to ensure users don't overwhelm blockchain peers with transactions.

- #### Ts/React Frontend
See [ReactBlockChat](https://github.com/denverquane/ReactBlockChat)

- #### GO GUI/"Frontend"
Planned for development, see Goals
