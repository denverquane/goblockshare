import * as React from 'react';

import './App.css';
// import { Transaction } from './Transaction';
import {
  Button,
  Toaster, Position,
  Intent,
  Tabs, Tab, TabId, Callout, Spinner
} from '@blueprintjs/core';
import { renderBlock, renderSimpleBlock, Block } from './Block';

import Sockette from 'sockette';

// const logo = require('./logo.svg');
export const BLOCKCHAIN_IP = 'localhost:5000';

const MyToaster = Toaster.create({
  className: 'my-toaster',
  position: Position.BOTTOM
});

interface SampleProps {
}

interface SampleState {
  connectionStatus: string;
  recentBlock: Block | null;
  blocks: Block[];
  openOverlay: boolean;
  currentTab: TabId;
}

function renderSpinner(status: string) {
  var intent: Intent;
  var value = 0;
  if (status === 'Disconnected') {
    intent = Intent.DANGER;
  } else if (status === 'Not Connecting') {
    intent = Intent.DANGER;
    value = 1;
  } else if (status === 'Connected') {
    intent = Intent.SUCCESS;
    value = 1;
  } else if (status === 'Reconnecting') {
    intent = Intent.WARNING;
  } else {
    intent = Intent.NONE;
  }
  if (value === 0) {
    return (
      <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center' }}>
        <Spinner intent={intent} large={true} />{status}</div>
    );
  } else {
    return (
      <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center' }}>
        <Spinner intent={intent} value={value} large={true} />{status}</div>
    );
  }
}

function renderTitle(status: string, block: Block) {
  return (
    <header className="App-header">
      <div style={{ display: 'flex', flexDirection: 'row' }}>
        {renderSpinner(status)}
        {block
          ? renderSimpleBlock(block)
          : <div />}
      </div>
    </header>
  );
}

export default class App extends React.Component<SampleProps, SampleState> {
  /*tslint:disable*/

  ws = new Sockette('ws://' + BLOCKCHAIN_IP + '/ws', {
    timeout: 5e3,
    maxAttempts: 10,
    onopen: (e: any) => {
      MyToaster.clear();
      this.setState({ ...this.state, connectionStatus: 'Connected' })
      MyToaster.show({
        message: 'Connected to Blockchain Server!',
        intent: Intent.SUCCESS,
      });
    },
    onmessage: (e: any) => {
      console.log('Received:', e.data);
      var recent: Block = e.data
      this.setState({ ...this.state, recentBlock: recent });
      console.log(e.data.Nonce)
    },
    onreconnect: (e: any) => {
      this.setState({ ...this.state, connectionStatus: 'Reconnecting' });
      MyToaster.show({
        message: 'Reconnecting to Blockchain Server...',
        intent: Intent.WARNING,
      })
    },
    onmaximum: (e: any) => {
      MyToaster.clear();
      this.setState({ ...this.state, connectionStatus: 'Not Connecting' })
      MyToaster.show({
        message: 'Maximum reconnect attempts to Blockchain Server; please refresh the page',
        intent: Intent.DANGER,
        timeout: 10000
      });
    },
    // onclose: (e: any) => MyToaster.show({
    //   message: 'Can\'t reach Blockchain server... Is it running?',
    //   intent: Intent.DANGER,
    // }),
    onerror: (e: any) => {
      MyToaster.clear();
      this.setState({ ...this.state, connectionStatus: 'Disconnected' })
      MyToaster.show({
        message: 'Error connecting to Blockchain Server!',
        intent: Intent.DANGER,
      });
    }
  });
  /*tslint:enable*/
  constructor(props: SampleProps) {
    super(props);

    this.state = {
      connectionStatus: 'Disconnected',
      recentBlock: null,
      blocks: [],
      openOverlay: false,
      currentTab: 'Home',
    };
    this.getBlocks = this.getBlocks.bind(this);
  }

  componentWillUnmount() {
    this.ws.close();
  }

  componentDidMount() {
    this.getBlocks();
  }

  render() {
    /*tslint:disable*/
    console.log('Rerendered')
    /*tslint:enable*/
    return (
      <div>
        {renderTitle(this.state.connectionStatus, this.state.blocks[0])}
        <h1 className="App-title" style={{ display: 'flex' }}>GoBlockShare!</h1>
        <Tabs id="MainPageTabs" onChange={this.handleTabChange} selectedTabId={this.state.currentTab}>
          <Tab
            id="Home"
            title="Home"
            panel={
              <Callout
                title="Welcome to GoBlockShare!"
              >
                <p>
                  This is the frontend web interface for GoBlockShare, and allows users to view the blockchain,
                  see user reputations, and even view available torrents (if the torrent backend service is running)
                </p>
              </Callout>
            }
          />
          <Tab
            id="Blockchain"
            title="Blockchain"
            panel={
              <div className="App">

                <Button
                  intent={Intent.SUCCESS}
                  onClick={() => {
                    this.getBlocks();
                  }}
                >
                  Update
                </Button>
                <div style={{ display: 'flex', flexDirection: 'row' }}>
                  <div style={{ width: '100%' }}>
                    {this.state.blocks.map((block: Block, index: number) => {
                      return (
                        (
                          renderBlock(block)
                        )
                      );
                    })}
                  </div>
                </div>

              </div>
            }
          />
        </Tabs>
      </div>
    );
  }

  getBlocks() {
    fetch('http://' + BLOCKCHAIN_IP + '/api/blockchain')
      .then(results => {
        return results.json();
      }).then(data => {
        let blocks = data.Blocks.map((block: Block) => {
          return block;
        });
        if (blocks !== this.state.blocks) {
          blocks = blocks.reverse();
          /*tslint:disable*/
          console.log(blocks[0]);
          this.setState({ ...this.state, blocks: blocks });
        }
      });
  }
  private handleTabChange = (newTab: TabId) => {
    if (newTab !== this.state.currentTab) {
      this.setState({ ...this.state, currentTab: newTab });
      if (newTab.toString() === 'Blockchain') {
        this.getBlocks();
      }
    }
  }
}
