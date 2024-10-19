import React from 'react';

import { Handle } from 'react-flow-renderer';

export const TransformerNodeComponent = ({ data }) => {
  return (
      <div className="transformerNode node">
        <Handle type="source" position="left" />
        <div>{data.name}</div>
        <Handle type="target" position="right" />
      </div>
  );
};
