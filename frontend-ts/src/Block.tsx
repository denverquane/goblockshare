import * as React from 'react';
import { Transaction } from './Transaction';
import { Callout, Intent } from '@blueprintjs/core';
import {
  ListGroup,
  ListGroupItem,
  Table,
} from 'react-bootstrap';

export interface Block {
  Index: number;
  Timestamp: string;
  Transactions: Transaction[];
  Hash: string;
  PrevHash: string;
}

interface BlockProps {
  block: Block;
}

interface BlockState {
  isOpen: boolean;
}

export class BlockDisplay extends React.Component<BlockProps, BlockState> {
  render() {
    return (
    <div style={{ display: 'flex', flexAlign: 'center', flexDirection: 'column' }}>
      <Table>
        <thead>
          <tr>
            <th style={{ width: '10%' }}>
              <Callout icon={null} intent={Intent.PRIMARY}>
                {this.props.block.Index}
              </Callout></th>
            <th style={{ width: '90%' }}>
              <Callout
                icon={null}
                title={this.props.block.Transactions[0]
                  ? '\'' + this.props.block.Transactions[0].Username + '\' added on ' + this.props.block.Timestamp
                  : 'Chain created on ' + this.props.block.Timestamp}
                intent={Intent.PRIMARY}
              />
            </th>
          </tr>

          <tr>
            <td><Callout icon="new-object" intent={Intent.SUCCESS}>New</Callout></td>
            <td>
              <ListGroup>
                <ListGroupItem>
                  <div style={{ display: 'flex' }}>

                    <div style={{ width: '90%' }}>{this.renderTransAsRow(this.props.block.Transactions.pop())}</div>
                  </div>

                </ListGroupItem>
              </ListGroup>
            </td>
          </tr>
          {(this.props.block.Transactions.length) > 0 ?
            <tr>
              <td><Callout icon="history" intent={Intent.WARNING}>Old</Callout>
              </td>
              <td>
                <ListGroup>
                  {this.props.block.Transactions.reverse().map((trans: Transaction, index) => {
                    return (
                      <ListGroupItem key={index}>
                        <div style={{ display: 'flex' }}>
                          <div style={{ width: '90%' }}>
                            {this.renderTransAsRow(trans)}
                          </div>
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
    </div>);
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