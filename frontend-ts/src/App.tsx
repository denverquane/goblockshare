import * as React from 'react';
import * as Blueprint from '@blueprintjs/core';
// import { Panel } from 'react-bootstrap';
import './App.css';

// import { EditableText } from 'blueprintjs/core';

const logo = require('./logo.svg');

interface SampleProps {
}

interface Block {
  Index: number;
  Timestamp: string;
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
          return block;
        });
        let newState = { blocks: blocks };
        this.setState(newState);
        /*tslint:disable*/
        console.log(this.state.blocks);
      });
  }

  render() {
    return (
      <div className="App">
        <header className="App-header">
          <img src={logo} className="App-logo" alt="logo" />
          <h1 className="App-title">Welcome to React</h1>
        </header>
        <div>
          {
            this.state.blocks.map((block: Block) => {
              return (
                <div key={block.Index}>
                  <Blueprint.Card >
                    {block.Timestamp}
                    </Blueprint.Card>
                </div>
              );
            })
          }
        </div>
      </div>
    );
  }
}
