import React from 'react';

import { Handle } from 'react-flow-renderer';

export const MergeNodeComponent = ({ data }) => {
  return (
      <div className="mergeNode node">
        <Handle type="source" position="left" />
        <div>{data.name}</div>
        <Handle type="target" position="right" />
      </div>
  );
};
