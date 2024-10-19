import React from 'react';

import { Handle } from 'react-flow-renderer';

export const CodeNodeComponent = ({ data }) => {
  return (
      <div className="codeNode node">
        <Handle type="source" position="left" />
        <div>{data.name}</div>
        <Handle type="target" position="right" />
      </div>
  );
};
