import * as React from 'react';
import { Card, Elevation, Tag, Intent } from '@blueprintjs/core';

export interface JSONRepSummary {
  ValidTorrFraction: number;
  QualityTorrFraction: number;
  AccurateTorrFraction: number;

  NotReceivedLayerFraction: number;
  ValidLayerFraction: number;
}

interface ReputationProps {
  address: string;
  summary: JSONRepSummary | null;
}

interface ReputationState {

}

export class ReputationDisplay extends React.Component<ReputationProps, ReputationState> {
  constructor(props: ReputationProps) {
    super(props);

  }

  renderTag(ratio: number, str: string) {
    var intent = Intent.DANGER;

    if (ratio > 0.8) {
      intent = Intent.SUCCESS;
    } else if (ratio > 0.5) {
      intent = Intent.WARNING;
    }
    return (
      <Tag large={true} intent={intent} round={true}>
        {ratio * 100}% {str}
      </Tag>
    );
  }

  render() {
    if (this.props.summary) {
      return (
        <Card elevation={Elevation.TWO}>
          <h6>{this.props.address}</h6>
          <p>
            Torrents: 
            {this.renderTag(this.props.summary.ValidTorrFraction, 'Valid')}
            {this.renderTag(this.props.summary.QualityTorrFraction, 'High Quality')}
            {this.renderTag(this.props.summary.AccurateTorrFraction, 'Accurate')}
          </p>
          <p>
            Layers:
            {this.renderTag(1.0 - this.props.summary.NotReceivedLayerFraction, 'Received')}
            {this.renderTag(this.props.summary.ValidLayerFraction, 'Valid')}
          </p>
        </Card>
      );
    } else {
      return (<div/>);
    }
  }
}