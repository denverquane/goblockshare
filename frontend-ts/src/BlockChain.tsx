import * as React from 'react';
import { Block, BlockDisplay } from './Block';

import {
  ListGroup,
  ListGroupItem,
} from 'react-bootstrap';

// import { Button } from 'react-bootstrap/lib/InputGroup';

interface ChainProps {
  blocks: Block[];
}

interface ChainState {
  isOpen: boolean;
}

export class ChainDisplay extends React.Component<ChainProps, ChainState> {
  constructor(props: ChainProps) {
    super(props);

    this.state = {
      isOpen: true
    };
  }

  componentWillReceiveProps(newProps: ChainProps) {
    if (newProps.blocks.length > this.props.blocks.length) {
      this.props = newProps;
    }
  }

  render() {
    return (
      <div>
        <ListGroup>
          {
            this.props.blocks.map((block: Block) => {
              return (
                (<ListGroupItem key={block.Index}>
                  <BlockDisplay
                    block={block}
                  />
                </ListGroupItem>
                )
              );
            })
          }
        </ListGroup>
      </div>
    );
  }
}