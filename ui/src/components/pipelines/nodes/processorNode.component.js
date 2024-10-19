import React from 'react';

import { Handle } from 'react-flow-renderer';

export const ProcessorNodeComponent = ({ data }) => {
  return (
      <div className="processorNode node">
        <Handle type="source" position="left" />
        <div>{data.name}</div>
        <Handle type="target" position="right" />
      </div>
  );
};
