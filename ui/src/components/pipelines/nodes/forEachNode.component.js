import React from 'react';

import { Handle } from 'react-flow-renderer';

export const ForEachNodeComponent = ({ data }) => {
  return (
      <div className="for_eachNode node">
        <Handle type="source" position="left" />
        <div>{data.name}</div>
        <Handle type="target" position="right" />
      </div>
  );
};
