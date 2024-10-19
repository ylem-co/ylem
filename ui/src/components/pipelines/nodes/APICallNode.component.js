import React from 'react';

import { Handle } from 'react-flow-renderer';

export const APICallNodeComponent = ({ data }) => {
  return (
      <div className="api_callNode node">
        <Handle type="source" position="left" />
        <div>{data.name}</div>
        <Handle type="target" position="right" />
      </div>
  );
};
