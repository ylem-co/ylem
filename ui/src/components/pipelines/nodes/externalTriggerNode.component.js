import React from 'react';

import { Handle } from 'react-flow-renderer';

export const ExternalTriggerNodeComponent = ({ data }) => {
  return (
      <div className="external_triggerNode node">
        <Handle type="source" position="left" />
        <div>{data.name}</div>
        <Handle type="target" position="right" />
      </div>
  );
};
