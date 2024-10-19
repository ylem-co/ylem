import React from 'react';

import { Handle } from 'react-flow-renderer';

export const NotificationNodeComponent = ({ data }) => {
  return (
      <div className="notificationNode node">
        <Handle type="source" position="left" />
        <div>{data.name}</div>
      </div>
  );
};
