import * as React from 'react';
import { Table } from 'react-bootstrap';
import { JSONRepSummary, ReputationDisplay } from './Reputation';
import * as TxTypes from './TxTypes';
import { Icon } from '@blueprintjs/core';
import { BLOCKCHAIN_IP } from './App';

export interface Transaction {
  Origin: {
    PubKeyX: string,
    PubKeyY: string,
    Address: string
  };
  Transaction: Object;
  TransactionType: string;
  R: string;
  S: string;
  TxID: string;
}

export interface AuthTransaction {
  OriginPubKeyX: string;
  OriginPubKeyY: string;
  OriginAddress: string;
  SignedMsg: string;
  TxRef: string[];
  R: string;
  S: string;
  DestAddr: string;
}

interface TransactionProps {
  transaction: Transaction | undefined;
}

interface TransactionState {
  alias: string;
  summary: JSONRepSummary | null;
  recipSummary: JSONRepSummary | null;
}

interface Alias {
  Data: string;
}

export function renderSimpleTransaction(trans: Transaction | undefined) {
  return (
    <div>
      {trans ? trans.Origin.Address : ''}
       
      {trans ? getIcon(trans.TransactionType) : ''}
       
      {trans ? trans.TransactionType : ''}
    </div>
  );
}

function getIcon(param: string) {
  switch (param) {
    case 'SET_ALIAS':
      return <Icon icon="tag" />;
    case 'SHARED_LAYER':
      return <Icon icon="social-media" />;
    case 'PUBLISH_TORRENT':
      return <Icon icon="document-share" />;
    case 'TORRENT_REP':
      return <Icon icon="chat" />;
    case 'LAYER_REP':
      return <Icon icon="chat" />;
    default:
      return <Icon icon="blank" />;
  }
}

function renderTransaction(param: Transaction) {
  switch (param.TransactionType) {
    case 'SET_ALIAS':
      return <div>{(param.Transaction as TxTypes.SetAliasTrans).Alias}</div>;
    case 'SHARED_LAYER':
      return <div>{(param.Transaction as TxTypes.SharedLayerTrans).SharedLayerHash.substr(0, 10)}...</div>;
    case 'TORRENT_REP':
      var rep = (param.Transaction as TxTypes.TorrentRepTrans);
      /*tslint:disable*/
      return (
        <div>
          <p>Accurate Name: {rep.RepMessage.AccurateName ? <Icon icon="tick" /> : <Icon icon="cross" />}</p>
          <p>High Quality: {rep.RepMessage.HighQuality ? <Icon icon="tick" /> : <Icon icon="cross" />}</p>
          <p>Was Valid: {rep.RepMessage.WasValid ? <Icon icon="tick" /> : <Icon icon="cross" />}</p>
        </div>
      );
    case 'PUBLISH_TORRENT':
      var torr = (param.Transaction as TxTypes.PublishTorrentTrans).Torrent;

      return (
        <div>
          <p>Name: '{torr.Name}'</p>
          <p>Size: {torr.TotalByteSize} bytes</p>
          <p>Layer Size: {torr.LayerByteSize} bytes</p>
        </div>
      );
    default:
      return <div />;
  }
}

export class TransactionDisplay extends React.Component<TransactionProps, TransactionState> {
  constructor(props: TransactionProps) {
    super(props);
    this.getAlias = this.getAlias.bind(this);
    this.state = {
      alias: 'NULL',
      summary: null,
      recipSummary: null,
    };
  }

  getAlias = (addr: string) => {
    fetch('http://' + BLOCKCHAIN_IP + '/api/alias/' + addr)
      .then(results => {
        return results.json();
      }).then(data => {
        let aliases = data.map((nalias: Alias) => {
          return nalias.Data;
        });
        this.setState({ alias: aliases[0] });
      });
  };

  getReputation = (addr: string, recip: boolean) => {
    fetch('http://' + BLOCKCHAIN_IP + '/api/reputation/' + addr)
      .then(results => {
        return results.json();
      }).then(data => {
        let reps = data.map((summary: JSONRepSummary) => {
          return summary;
        });
        if (!recip) {
          this.setState({ summary: reps[0] });
        } else {
          this.setState({ recipSummary: reps[0] });
        }
      });
  };

  componentDidMount() {
    if (this.props.transaction) {
      this.getAlias(this.props.transaction.Origin.Address);
      this.getReputation(this.props.transaction.Origin.Address, false);
      if (this.props.transaction.TransactionType === 'SHARED_LAYER') {
        this.getReputation((this.props.transaction.Transaction as TxTypes.SharedLayerTrans).Recipient, true);
      }
    }
  }

  renderDestination(param: Transaction) {
    switch (param.TransactionType) {
      case 'SET_ALIAS':
        return <div>{this.state.alias}</div>;
      case 'SHARED_LAYER':
        return (
          <div>
            <ReputationDisplay address={(param.Transaction as TxTypes.SharedLayerTrans).Recipient} summary={this.state.recipSummary} />
          </div>
        );
      case 'TORRENT_REP':
        return <div>{(param.Transaction as TxTypes.TorrentRepTrans).TxID.substr(0, 10)}...</div>;
      default:
        return <div />;
    }
  }

  render() {
    if (this.props.transaction !== undefined) {
      var addr = this.state.alias === 'NULL' || this.props.transaction.TransactionType === 'SET_ALIAS'
      ? this.props.transaction.Origin.Address : this.state.alias
      return (
        <div>
          <h5>TxID: {this.props.transaction.TxID}</h5>
          <Table condensed={true}>
            <thead>
              <tr>
                <th>Origin</th>
                <th>
                  <div>
                    {getIcon(this.props.transaction.TransactionType)} {this.props.transaction.TransactionType}
                  </div>
                </th>
                <th>{this.props.transaction.TransactionType === 'TORRENT_REP'
                  || this.props.transaction.TransactionType === 'LAYER_REP' ? 'Tx Reference' : 'Recipient'}</th>
              </tr>
              <tr>
                <td style={{ width: '20%' }}>
                  {<ReputationDisplay
                    address={addr}
                    summary={this.state.summary}
                  />}
                </td>
                <td>
                  {
                    renderTransaction(this.props.transaction)
                  }
                </td>
                <td>
                  {this.renderDestination(this.props.transaction)}
                </td>
                {/* <td style={{ width: '10%' }}>{this.props.transaction.TransactionType}</td> */}
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
