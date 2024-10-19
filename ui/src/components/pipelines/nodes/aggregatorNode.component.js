import React from 'react';

import { Handle } from 'react-flow-renderer';

export const AggregatorNodeComponent = ({ data }) => {
  return (
      <div className="aggregatorNode node">
        <Handle type="source" position="left" />
        <div>{data.name}</div>
        <Handle type="target" position="right" />
      </div>
  );
};
