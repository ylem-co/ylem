import React, { Component } from 'react';

import { nodeColors } from './pipeline.component';

import Tooltip from '@mui/material/Tooltip';

class PipelineSidebar extends Component {
  constructor(props) {
      super(props);
      this.onDragStart = this.onDragStart.bind(this);
  }

  onDragStart = (e, nodeType) => {
    e.dataTransfer.setData('application/reactflow', nodeType);
    e.dataTransfer.effectAllowed = 'move';
  };

  render() {
    let nodesIterator = 1

    return (
      <aside className="pipelineSidebar">
        {
          Object.entries(nodeColors).map(([key, value]) => (
            <>
              <Tooltip key={key} title={key.charAt(0).toUpperCase() + key.slice(1)} placement="top">
                <div 
                  className={key + "Node draggableNode"} 
                  onDragStart={(event) => this.onDragStart(event, key)} 
                  draggable={this.props.isDraggingAllowed}
                >
                  {key[0].charAt(0).toUpperCase() + key[1].charAt(0)}
                </div>
              </Tooltip>
              {
                (nodesIterator++ && nodesIterator % 2 === 0)
                ? <div className="pipelineSideBarBreak"></div>
                : <><div className="pipelineSideBarBreak"></div><div className="clearfix"></div></>
              }
            </>
          ))
        }
      </aside>
    );
  }
};

export default PipelineSidebar;
