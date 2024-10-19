import React from 'react';

import { Handle } from 'react-flow-renderer';

import Tooltip from '@mui/material/Tooltip';

export const CONDITION_CONNECTOR_TRUE = "condition-connector-true";
export const CONDITION_CONNECTOR_FALSE = "condition-connector-false";

export const ConditionNodeComponent = ({ data }) => {
  return (
    <>
        <div className="conditionNode node">
          <Handle type="source" position="left" />
            <div>{data.name}</div>
            <Tooltip title="true" placement="right">
              <Handle 
                type="target" 
                position="right" 
                id={CONDITION_CONNECTOR_TRUE} 
                style={{ top: '30%' }}
              />
            </Tooltip>
            <Tooltip title="false" placement="right">
              <Handle 
                type="target" 
                position="right" 
                id={CONDITION_CONNECTOR_FALSE} 
                style={{ top: '70%' }}
              />
            </Tooltip>
        </div>
    </>
  );
};
