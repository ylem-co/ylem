import React from 'react';

import { Handle } from 'react-flow-renderer';

export const PipelineRunNodeComponent = ({ data }) => {
  return (
      <div className="pipelineRunNode node">
        <Handle type="source" position="left" />
        <div>{data.name}</div>
      </div>
  );
};
