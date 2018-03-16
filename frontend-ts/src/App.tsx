import * as React from 'react';
// import * as Blueprint from '@blueprintjs/core';
import {
  ListGroup,
  ListGroupItem,
  Table,
  Alert,
} from 'react-bootstrap';
import './App.css';
import { Transaction } from './Transaction';

// import { EditableText } from 'blueprintjs/core';

const logo = require('./logo.svg');

interface SampleProps {
}

interface Block {
  Index: number;
  Timestamp: string;
  Transactions: Transaction[];
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
            })
          };
        });
        let newState = { blocks: blocks };
        this.setState(newState);
        /*tslint:disable*/
        console.log(data);
      });
  }

  renderTransAsRow(trans: Transaction | undefined) {
    if (trans !== undefined) {
      return <tr>
        <td style={{ width: '10%' }}><Alert>Added</Alert></td>
        <td style={{ width: '90%' }}><div>
          <Table>
            <thead>
              <tr>
                <th>Author</th>
                <th>Channel</th>
                <th>Message</th>
              </tr>
              <tr>
                <td style={{ width: '10%' }}>{trans.Author}</td>
                <td style={{ width: '10%' }}>{trans.Channel}</td>
                <td style={{ width: '80%' }}>{trans.Message}</td>
              </tr>
            </thead>
          </Table>
        </div>
        </td>
      </tr>
    } else {
      return <tr><td style={{ width: '10%' }} />
        <td style={{ width: '90%' }}>Initial Block; No transactions</td>
      </tr>
    }
  }

  render() {
    return (
      <div className="App">
        <header className="App-header">
          <img src={logo} className="App-logo" alt="logo" />
          <h1 className="App-title">Welcome to React</h1>
        </header>
        <ListGroup>
          {
            this.state.blocks.reverse().map((block: Block) => {
              return (
                <ListGroupItem>
                  <div key={block.Index} style={{ display: 'flex', flexAlign: 'center', flexDirection: 'column' }}>
                    <Table>
                      <thead>
                        <tr>
                          <th>{block.Index}</th>
                          <th>{block.Timestamp}</th>
                        </tr>
                        {/* <ListGroup>{block.Transactions.reverse().map((trans: Transaction, index) => {
                              return <ListGroupItem>
                                <div key={index} style={{ display: 'flex', flexAlign: 'left' }}>
                                  <div><b>Author: </b>{trans.Author}</div>
                                  <div><b>Channel: </b>{trans.Channel}</div>
                                  <div><b>Message: </b>{trans.Message}</div>
                                </div>
                              </ListGroupItem>;
                            })}
                            </ListGroup> */}
                        {this.renderTransAsRow(block.Transactions.pop())}
                      </thead>
                    </Table>
                  </div>
                </ListGroupItem>
              );
            })
          }
        </ListGroup>
      </div>
    );
  }
}
