import * as React from 'react';

import './App.css';
import { Transaction } from './Transaction';
import { ChainDisplay } from './BlockChain';
import { Button, Toaster, Position, Intent } from '@blueprintjs/core';

const logo = require('./logo.svg');

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
    this.getBlocks = this.getBlocks.bind(this);
  }
  
  componentDidMount() {
    this.getBlocks();
  }

  render() {
    return (
      <div className="App">
        <header className="App-header">
          <img src={logo} className="App-logo" alt="logo" onClick={this.getBlocks} />
          <h1 className="App-title">Welcome to React</h1>
        </header>
        <Button onClick={this.getBlocks}>Update</Button>
        <ChainDisplay blocks={this.state.blocks} />
      </div>
    );
  }

  getBlocks() {
    let handle = MyToaster.show({
      message: 'Fetching data',
      intent: Intent.PRIMARY,
    });
    fetch('http://localhost:8040/chain')
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
    MyToaster.dismiss(handle);
  }
}
