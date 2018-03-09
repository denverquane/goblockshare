import * as React from "react";
import { addTransaction } from './actions'

import {
  EditableText,
} from '@blueprintjs/core';

interface TransactionState {
  author: string;
  channel: string;
  message: string;
}

interface TransactionProps {

}

export default class App extends React.Component<TransactionProps, TransactionState> {
  constructor(props: TransactionProps) {
    super(props);
    //this.store = this.props.store;

    this.state = {
      author: 'Test',
      channel: 'Channel',
      message: 'hello'
    }
  }

  handleAddTransaction = () => {
    let trans = {
      author: this.state.author,
      channel: this.state.channel,
      message: this.state.message
    }
    //this.store.dispatch(addTransaction(trans))
  }

  render() {
    // return (
    //   <div className='App'>
    //     <div>
    //       <EditableText
    //         defaultValue={this.state.author}
    //         confirmOnEnterKey={true}
    //         onConfirm={(author) => {
    //           this.setState({author})
    //         }}
    //       />
    //     </div>
    //     <div>
    //       <EditableText
    //         placeholder={this.state.channel}
    //         onConfirm={(channel) => {
    //           this.setState({channel})
    //         }}
    //       />
    //     </div>
    //     <div>
    //       <EditableText
    //         placeholder={this.state.message}
    //         onConfirm={(message) => {
    //           this.setState({message})
    //         }}
    //       />
    //     </div>
    //     <div className="col-xs-2">
    //       <img src='logo.svg' height="96" alt='React' onClick={this.handleAddTransaction}></img>
    //     </div>
    //   </div>
    // );
    return (<div/>);
  }
}
