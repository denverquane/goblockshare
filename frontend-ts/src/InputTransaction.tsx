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

interface InputProps {
    isOverlayOpen: boolean;
}

interface InputState {
    isOverlayOpen: boolean;
    username: string;
    password: string;
    message: string;
}

export class InputTransaction extends React.Component<InputProps, InputState> {

    constructor(props: InputProps) {
        super(props);

        this.state = {
            isOverlayOpen: this.props.isOverlayOpen,
            username: 'username',
            password: 'password',
            message: 'message'
        };
    }

    componentWillReceiveProps(newProps: InputProps) {
        this.setState(newProps);
    }

    render() {
        return (
            <Overlay isOpen={this.state.isOverlayOpen} >
                <Panel>
                    <div className={Classes.CARD}>

                        <h3>Please provide the details of the transaction you wish to post:</h3>

                        <EditableText
                            placeholder={this.state.username}
                            confirmOnEnterKey={true}
                            onConfirm={(val: string) => {
                                this.setState({ username: val });
                            }}
                        />
                        <EditableText
                            placeholder={this.state.password}
                            confirmOnEnterKey={true}
                            onConfirm={(val: string) => {
                                this.setState({ password: val });
                            }}
                        />
                        <EditableText
                            placeholder={this.state.message}
                            confirmOnEnterKey={true}
                            onConfirm={(val: string) => {
                                this.setState({ message: val });
                            }}
                        />

                        <Button
                            onClick={() => {
                                /*tslint:disable*/
                                console.log(this.state.username + " " + this.state.password + " " + this.state.message);
                                /*tslint:enable*/
                            }}
                        >
                            Post Transaction
                        </Button>
                        <Button
                            intent={Intent.DANGER}
                            onClick={() => {
                                this.setState({ isOverlayOpen: !this.state.isOverlayOpen });
                            }}
                        >
                            Close
                        </Button>
                    </div>

                </Panel>
            </Overlay>
        );
    }
}