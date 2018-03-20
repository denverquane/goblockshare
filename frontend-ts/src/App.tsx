import * as React from 'react';

import './App.css';
import { Transaction } from './Transaction';
import { ChainExplorer } from './ChainExplorer';

const logo = require('./logo.svg');

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
}

export default class App extends React.Component<SampleProps, SampleState> {

  constructor(props: SampleProps) {
    super(props);

    this.state = {
      blocks: [],
    };
  }

  componentDidMount() {
    fetch('http://localhost:8040')
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
        let newState = { blocks: blocks };
        this.setState(newState);
      });
  }

  render() {
    return (
      <div className="App">
        <header className="App-header">
          <img src={logo} className="App-logo" alt="logo" />
          <h1 className="App-title">Welcome to React</h1>
        </header>
        <ChainExplorer blocks={this.state.blocks}/>
      </div>
    );
  }
}
