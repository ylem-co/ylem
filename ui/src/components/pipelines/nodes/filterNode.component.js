import React from 'react';

import { Handle } from 'react-flow-renderer';

export const FilterNodeComponent = ({ data }) => {
  return (
      <div className="filterNode node">
        <Handle type="source" position="left" />
        <div>{data.name}</div>
        <Handle type="target" position="right" />
      </div>
  );
};
