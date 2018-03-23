import * as React from 'react';

import {
    Overlay,
    Button,
    EditableText,
    Intent,
    Classes,
} from '@blueprintjs/core';

import {
    Panel,
} from 'react-bootstrap';

import { AuthTransaction } from './Transaction';

interface InputProps {
    isOverlayOpen: boolean;
    BACKEND_IP: string;
    onClose: () => void;
}

interface InputState {
    transaction: AuthTransaction;
}

export class InputTransaction extends React.Component<InputProps, InputState> {

    constructor(props: InputProps) {
        super(props);

        this.state = {
            transaction: {
                Username: 'username',
                Password: 'password',
                Channel: 'channel',
                Message: 'message',
                TransactionType: 'ADD_MESSAGE'
            }
        };
    }

    componentWillReceiveProps(newProps: InputProps) {
        this.props = newProps;
    }

    render() {
        return (
            <Overlay isOpen={this.props.isOverlayOpen} >
                <Panel>
                    <div className={Classes.CARD}>

                        <h3>Please provide the details of the transaction you wish to post:</h3>

                        <EditableText
                            placeholder={this.state.transaction.Username}
                            confirmOnEnterKey={true}
                            onConfirm={(val: string) => {
                                this.setState({
                                    transaction: {
                                        ...this.state.transaction,
                                        Username: val
                                    }
                                });
                            }}
                        />
                        <EditableText
                            placeholder={this.state.transaction.Password}
                            confirmOnEnterKey={true}
                            onConfirm={(val: string) => {
                                this.setState({
                                    transaction: {
                                        ...this.state.transaction,
                                        Password: val
                                    }
                                });
                            }}
                        />
                        <EditableText
                            placeholder={this.state.transaction.Message}
                            confirmOnEnterKey={true}
                            onConfirm={(val: string) => {
                                this.setState({
                                    transaction: {
                                        ...this.state.transaction,
                                        Message: val
                                    }
                                });
                            }}
                        />

                        <Button
                            onClick={() => {
                                fetch(this.props.BACKEND_IP + '/postTransaction', {
                                    method: 'POST',
                                    mode: 'no-cors',
                                    headers: {
                                        'Access-Control-Allow-Origin': '*',
                                        'Accept': 'application/json',
                                        'Content-Type': 'application/json',
                                    },
                                    body: JSON.stringify(
                                        this.state.transaction
                                    )
                                })
                                    .then(results => {
                                        return results;
                                    }).then(data => {
                                        // let blocks = data.Blocks.map((block: string) => {
                                        //     return block;
                                        // });
                                        /*tslint:disable*/
                                        console.log({ data });
                                    });
                                    this.props.onClose();
                            }}
                        >
                            Post Transaction
                        </Button>
                        <Button
                            intent={Intent.DANGER}
                            onClick={() => {
                                this.props.onClose();
                            }}
                        >Close
                        </Button>
                    </div>
                </Panel>
            </Overlay>
        );
    }
}