import * as React from 'react';
import { Transaction } from './Transaction';

import {
  ListGroup,
  ListGroupItem,
  Table,
} from 'react-bootstrap';
import { Callout, Intent } from '@blueprintjs/core';

interface SampleProps {
  blocks: Block[];
}

interface Block {
  Index: number;
  Timestamp: string;
  Transactions: Transaction[];
  Hash: string;
  PrevHash: string;
}

interface SampleState {
  isOpen: boolean;
}

export class ChainExplorer extends React.Component<SampleProps, SampleState> {
  constructor(props: SampleProps) {
    super(props);

    this.state = {
      isOpen: false
    };
  }

  // componentDidMount() {
  //   fetch('http://localhost:8040')
  //     .then(results => {
  //       return results.json();
  //     }).then(data => {
  //       let blocks = data.Blocks.map((block: Block) => {
  //         return {
  //           Index: block.Index,
  //           Timestamp: block.Timestamp,
  //           Transactions: block.Transactions.map((trans: Transaction) => {
  //             return trans;
  //           }),
  //           Hash: block.Hash,
  //           PrevHash: block.PrevHash
  //         };
  //       });
  //       let newState = { blocks: blocks };
  //       this.setState(newState);
  //     });

  componentWillReceiveProps(newProps: SampleProps) {
    if (newProps.blocks.length > this.props.blocks.length) {
      this.props = newProps;
    }
  }

  render() {
    return (
      <ListGroup>
        {
          this.props.blocks.reverse().map((block: Block) => {
            return (
              (<ListGroupItem>
                <div key={block.Index} style={{ display: 'flex', flexAlign: 'center', flexDirection: 'column' }}>
                  <Table>
                    <thead>
                      <tr>
                        <th style={{ width: '10%' }}>
                          <Callout icon={null} intent={Intent.PRIMARY}>
                            {block.Index}
                          </Callout></th>
                        <th style={{ width: '90%' }}>
                          <Callout
                            icon={null}
                            title={block.Transactions[0]
                              ? '\'' + block.Transactions[0].Username + '\' added on ' + block.Timestamp
                              : 'Chain created on' + block.Timestamp}
                            intent={Intent.PRIMARY}
                          >
                            {/* <h5 class="pt-callout-title">Callout Heading</h5>
                              Lorem ipsum dolor sit amet, consectetur adipisicing elit. Ex, delectus! */}
                          </Callout>
                        </th>
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

                      <tr>
                        <td><Callout icon="new-object" intent={Intent.SUCCESS}>New</Callout></td>
                        <td>
                          <ListGroup>
                            <ListGroupItem>
                              <div style={{ display: 'flex' }}>

                                <div style={{ width: '90%' }}>{this.renderTransAsRow(block.Transactions.pop())}</div>
                              </div>

                            </ListGroupItem>
                          </ListGroup>
                        </td>
                      </tr>
                      {(block.Transactions.length) > 0 ?
                        <tr>
                          <td><Callout icon="history" intent={Intent.WARNING}>Old</Callout></td>
                          <td>
                            <ListGroup>
                              {block.Transactions.reverse().map((trans: Transaction, index) => {
                                return (
                                  <ListGroupItem key={index}>
                                    <div style={{ display: 'flex' }}>
                                      <div style={{ width: '90%' }}>{this.renderTransAsRow(trans)}</div>
                                    </div>
                                  </ListGroupItem>
                                );
                              }
                              )}
                            </ListGroup>
                          </td>
                        </tr> : <tr />}
                    </thead>
                  </Table>
                </div>
                {/* <div style={{ display: 'flex' }}>
                    <div style={{ width: '50%' }}><b>Hash: </b>{block.Hash} </div>
                    <div style={{ width: '50%' }}><b>PrevHash: </b>{block.PrevHash}</div>
                  </div> */}
              </ListGroupItem>
              )
            );
          })
        }
      </ListGroup>
    );
  }
  renderTransAsRow(trans: Transaction | undefined) {
    if (trans !== undefined) {
      return (
        <div>
          <Table condensed={true}>
            <thead>
              <tr>
                <th>Channel</th>
                <th>Type</th>
                <th>Message</th>
              </tr>
              <tr>
                <td style={{ width: '15%' }}>{trans.Channel}</td>
                <td style={{ width: '15%' }}>{trans.TransactionType}</td>
                <td style={{ width: '70%' }}>{trans.Message}</td>
              </tr>
            </thead>
          </Table>
        </div>
      );
    } else {
      return (
        <div>
          <Table condensed={true}>
            <thead>
              <tr><td style={{ width: '10%' }} />
                <td style={{ width: '90%' }}>Initial Block; No transactions</td>
              </tr>
            </thead>
          </Table>
        </div>
      );
    }
  }
}