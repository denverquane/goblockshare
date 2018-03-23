import * as React from 'react';
import { Table } from 'react-bootstrap';
import { Callout, IconName, Intent } from '@blueprintjs/core';

export interface Transaction {
  Username: string;
  Channel: string;
  Message: string;
  TransactionType: string;
}

export interface AuthTransaction {
  Username: string;
  Password: string;
  Channel: string;
  Message: string;
  TransactionType: string;
}

interface TransactionProps {
  transaction: Transaction | undefined;
}

interface TransactionState {

}

export class TransactionDisplay extends React.Component<TransactionProps, TransactionState> {
  constructor(props: TransactionProps) {
    super(props);
  }

  render() {
    if (this.props.transaction !== undefined) {
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
                <td style={{ width: '20%' }}>{this.props.transaction.Channel}</td>
                <td style={{ width: '20%' }}>{this.renderTransType(this.props.transaction.TransactionType)}</td>
                <td style={{ width: '60%' }}>{this.props.transaction.Message}</td>
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

  renderTransType(type: string) {
    var iconn: IconName;
    var text: string;
    var intent: Intent;
    switch (type) {
      case 'ADD_MESSAGE':
        iconn = 'add';
        text = 'Add Message';
        intent = Intent.SUCCESS;
        break;
      case 'DELETE_MESSAGE':
        iconn = 'trash';
        text = 'Delete Message';
        intent = Intent.DANGER;
        break;
      case 'ADD_USER':
        iconn = 'new-person';
        text = 'Add User';
        intent = Intent.SUCCESS;
        break;
      default:
        iconn = 'cross';
        text = 'INVALID';
        intent = Intent.DANGER;
        break;
    }
    return (
      <Callout icon={iconn} intent={intent}>{text}</Callout>
    );
  }
}