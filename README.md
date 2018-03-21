# GoBlockChat
[![GoDoc](https://godoc.org/github.com/denverquane/GoBlockChat?status.png)](https://godoc.org/github.com/denverquane/GoBlockChat)

This app seeks to demonstrate using blockchain tech. for a chat application like Slack or Discord.

By using the blockchain, no user can delete, modify, or falsify chat records and exchanges without other users
being aware of the change. This not ensures integrity of the chat, but allows for interesting functionality 
regarding "rewinding" or rollback of chat engagements.

## Goals:
- [ ] Proof of Work for posting messages/blocks (prevent spam/abuse)
  - [X] Basic difficulty/cryptographic proof validation
  - [ ] Scaling difficulty of blocks as the chain grows
  - [ ] Rewards for propagating the chain (?)
- [ ] Automatic discovery of other Nodes, and auto propagation of the blockchain itself
- [X] Author/poster validation (login validation)
- [X] Basic JS Frontend for viewing the blockchain in realtime
- [ ] JS Frontend for posting, deleting, editing, etc. messages and transactions
- [ ] GO app to interact with the chain, without the Webapp (?)

## Project Structure:
This app is comprised of two distinct parts as of 3/21/18.
These parts serve as the frontend and backend services for the overarching application.

- #### GO Backend
In the aptly-named "Go" directory are all the [GoLang](https://golang.org/) source files, which run the backend application that handles blockchain operations, including peer discovery/chain propagation, author and user validation, and rudimentary ["Proof of Work"](https://en.wikipedia.org/wiki/Proof-of-work_system) calculations to ensure users don't overwhelm blockchain peers with transactions.

- #### - Ts/React Frontend
In the "frontend-ts" directory are the TypeScript/Javascript/React files that comprise the frontend Webapp. This application currently only allows for a user-friendly view of the underlying blockchain structure (assuming the backend application is also running), but an interactive interface for posting messages and adding transactions to the blockchain is planned for future development.
