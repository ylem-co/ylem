import React from 'react';

import { Handle } from 'react-flow-renderer';

export const QueryNodeComponent = ({ data }) => {
  return (
      <div className="queryNode node">
        <Handle type="source" position="left" />
        <div>{data.name}</div>
        <Handle type="target" position="right" />
      </div>
  );
};
