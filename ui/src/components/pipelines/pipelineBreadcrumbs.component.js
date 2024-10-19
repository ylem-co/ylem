import React from 'react';
import { Link } from 'react-router-dom';

import { PIPELINE_ROOT_FOLDER } from "../../services/pipeline.service"; 

export const PipelineBreadcrumbs = ({ folder, item = null, parentFolder, onDrop, onDragOver, type }) => {
  return (
    <div className="pb-4">
      <span className="bc pipelineBreadcrumb">
        {
          (item !== null || folder !== null)
          && 
          <Link 
            to={"/" + type} 
            className="pipelineBreadcrumb droppable"
            droppable="true"
            onDragOver={(e)=>onDragOver(e)}
            onDrop={(e)=>onDrop(e, PIPELINE_ROOT_FOLDER)}
          >{type.charAt(0).toUpperCase() + type.slice(1)}</Link>
        }
        {
          folder !== null
          &&
          <span> &gt; ... &gt;</span>
        }
        {
          parentFolder !== null 
          && 
          <>
            &nbsp;<Link 
                    to={"/" + type + "/folder/" + parentFolder.uuid} 
                    className="pipelineBreadcrumb droppable"
                    onDrop={(e)=>onDrop(e, parentFolder)}
                    onDragOver={(e)=>onDragOver(e)}
                    droppable="true"
                  >{parentFolder.name}</Link> &gt;
          </>
        }
        {
          folder !== null
          &&
          <span>
            &nbsp;<Link 
              to={"/" + type + "/folder/" + folder.uuid} 
              className="pipelineBreadcrumb"
            >{folder.name}</Link>
          </span>
        }
        {
          item !== null
          && 
          <span>{" > " + item.name}</span>
        }
      </span>
    </div>
  );
};
