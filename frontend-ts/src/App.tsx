import * as React from 'react';

import './App.css';
import { Transaction } from './Transaction';
import { ChainDisplay } from './BlockChain';
import { Button, Toaster, Position, Intent, Callout } from '@blueprintjs/core';

const logo = require('./logo.svg');
const IP = 'http://localhost:8040';

const MyToaster = Toaster.create({
  className: 'my-toaster',
  position: Position.TOP
});

interface SampleProps {
}

interface Block {
  Index: number;
  Timestamp: string;
  Transactions: Transaction[];
  Hash: string;
  PrevHash: string;
}

interface SampleState {
  blocks: Block[];
  users: string[];
}

export default class App extends React.Component<SampleProps, SampleState> {

  constructor(props: SampleProps) {
    super(props);

    this.state = {
      blocks: [],
      users: [],
    };
    this.getBlocksAndUsers = this.getBlocksAndUsers.bind(this);
  }

  componentDidMount() {
    this.getBlocksAndUsers();
  }

  render() {
    return (
      <div className="App">
        <header className="App-header">
          <img src={logo} className="App-logo" alt="logo" onClick={this.getBlocksAndUsers} />
        </header>
        <h1 className="App-title">Welcome to the GoBlockChat!</h1>
        <Button onClick={this.getBlocksAndUsers}>Update</Button>
        <div>{this.renderUsers()}</div>
        <ChainDisplay blocks={this.state.blocks} />
      </div>
    );
  }

  getBlocksAndUsers() {
    let handle = MyToaster.show({
      message: 'Fetching data',
      intent: Intent.PRIMARY,
    });
    fetch(IP + '/chain')
      .then(results => {
        return results.json();
      }).then(data => {
        let blocks = data.Blocks.map((block: Block) => {
          return {
            Index: block.Index,
            Timestamp: block.Timestamp,
            Transactions: block.Transactions.map((trans: Transaction) => {
              return trans;
            }),
            Hash: block.Hash,
            PrevHash: block.PrevHash
          };
        });
        let newState = { ...this.state, blocks: blocks };
        this.setState(newState);
      });
    fetch(IP + '/users')
      .then(results => {
        return results.json();
      }).then(data => {
        let users = data.map((user: string) => {
          return user.split(':')[0];
        });
        let newState = { ...this.state, users: users };
        this.setState(newState);
      });

    MyToaster.dismiss(handle);
  }

  renderUsers() {
    return (
      <div style={{ display: 'flex', flexDirection: 'row' }}>
        <Callout icon={null} intent={Intent.PRIMARY} style={{ marginRight: '1%', width: '20%' }}>
          Authorized Users:
        </Callout>
        {this.state.users.map((user: string, index: number) => {
          return (
            <Callout key={index} style={{ width: '10%', marginRight: '1%' }}>{user}</Callout>
          );
        })}
      </div >
    );
  }
}
