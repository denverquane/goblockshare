import * as React from 'react';
import { Transaction, TransactionDisplay, renderSimpleTransaction } from './Transaction';
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
  Difficulty: number;
  Nonce: string;
}

export function renderSimpleBlock(block: Block) {
  return (
    <div style={{ display: 'flex', flexDirection: 'column', width: '100%' }}>
      Most recent block:
      <div style={{ display: 'flex', flexDirection: 'column', width: '80%' }}>
        <div style={{ display: 'flex', flexDirection: 'row' }}>
          <Callout icon={null} intent={Intent.PRIMARY} style={{ width: '10%' }}>
            {block.Index}
          </Callout>
          <div>
            <Callout
              style={{ width: '90%' }}
              icon={null}
              title={block.Index !== 0
                ? 'Added on ' + block.Timestamp
                : 'Chain created on ' + block.Timestamp}
              intent={Intent.PRIMARY}
            />
            <div>
              {block.Transactions ? renderSimpleTransaction(block.Transactions[0])
                 : <div />
              }
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

export function renderBlock(block: Block) {
  return (
    <div style={{ display: 'flex', flexDirection: 'column', maxWidth: '100%' }}>
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
                title={block.Transactions && block.Transactions[0]
                  ? 'Added on ' + block.Timestamp
                  : 'Chain created on ' + block.Timestamp}
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

                    <div style={{ width: '100%', maxWidth: '100%' }}>
                      {block.Transactions ? <TransactionDisplay
                        transaction={block.Transactions.pop()}
                      /> : <div />
                      }
                    </div>
                  </div>
                </ListGroupItem>
              </ListGroup>
            </td>
          </tr>
        </thead>
      </Table>
    </div>);
}