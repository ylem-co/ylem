import React from 'react';

import { Handle } from 'react-flow-renderer';

export const GptNodeComponent = ({ data }) => {
  return (
      <div className="gptNode node">
        <Handle type="source" position="left" />
        <div>{data.name}</div>
        <Handle type="target" position="right" />
      </div>
  );
};
