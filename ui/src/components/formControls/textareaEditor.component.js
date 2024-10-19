import React, { Fragment } from "react";

import Tooltip from '@mui/material/Tooltip';

import { PIPELINE_TYPE_GENERIC, PIPELINE_TYPE_METRIC} from "../../services/pipeline.service";
 
const addText = (e, txtId, callback, txt = null, brackets = false) => {
  let clickedElement = e.target;
  let text = 
    txt === null 
        ? clickedElement.innerHTML + "()" 
        : txt;

  if (brackets === true) {
    text = "{{ " + text + " }}";
  }
  
  let txtElement = document.getElementById(txtId);

  insertAtCursor(txtElement, callback, text);
}

const insertAtCursor = (myField, callback, myValue) => {
    //IE support
    if (document.selection) {
        myField.focus();
        var sel = document.selection.createRange();
        sel.text = myValue;
    }
    //MOZILLA and others
    else if (myField.selectionStart || myField.selectionStart === '0') {
        var startPos = myField.selectionStart;
        var endPos = myField.selectionEnd;
        myField.value = myField.value.substring(0, startPos)
            + myValue
            + myField.value.substring(endPos, myField.value.length);
    } else {
        myField.value += myValue;
    }

    callback(myField);
}

export const TextareaEditor = ({ txtId, callback, brackets = false, pipelineType = PIPELINE_TYPE_GENERIC }) => (
  <Fragment>
    <Tooltip title="AVG(field_name) - average value in the input data set" placement="bottom">
        <div className="textareaEditorButton" onClick={(e) => addText(e, txtId, callback, "AVG(field_name)", brackets)}>AVG</div>
    </Tooltip>
    <Tooltip title="SUM(field_name) - sum of values in the input data set" placement="bottom">
        <div className="textareaEditorButton" onClick={(e) => addText(e, txtId, callback, "SUM(field_name)", brackets)}>SUM</div>
    </Tooltip>
    <Tooltip title="COUNT() - number of rows in the input data set" placement="bottom">
        <div className="textareaEditorButton" onClick={(e) => addText(e, txtId, callback, "COUNT()", brackets)}>COUNT</div>
    </Tooltip>
    <Tooltip title="FIRST(field_name) - first value in the input data set" placement="bottom">
        <div className="textareaEditorButton" onClick={(e) => addText(e, txtId, callback, "FIRST(field_name)", brackets)}>FIRST</div>
    </Tooltip>
    <Tooltip title="LAST(field_name) - last value in the input data set" placement="bottom">
        <div className="textareaEditorButton" onClick={(e) => addText(e, txtId, callback, "LAST(field_name)", brackets)}>LAST</div>
    </Tooltip>
    <Tooltip title="MAX(field_name) - max value in the input data set" placement="bottom">
        <div className="textareaEditorButton" onClick={(e) => addText(e, txtId, callback, "MAX(field_name)", brackets)}>MAX</div>
    </Tooltip>
    <Tooltip title="MIN(field_name) - min value in the input data set" placement="bottom">
        <div className="textareaEditorButton" onClick={(e) => addText(e, txtId, callback, "MIN(field_name)", brackets)}>MIN</div>
    </Tooltip>
    <div className="note float-left px-2 pt-3">{"Array functions"}</div>
    <div className="clearfix">
        <Tooltip title='ROUND(field_name, precision, "floor|ceil") - round a number up or down with a certain precision' placement="bottom">
            <div className="textareaEditorButton" onClick={(e) => addText(e, txtId, callback, 'ROUND(number, 2, "floor")', brackets)}>ROUND</div>
        </Tooltip>
        <Tooltip title='ABS(field_name) - the absolute value of a number' placement="bottom">
            <div className="textareaEditorButton" onClick={(e) => addText(e, txtId, callback, 'ABS(number)', brackets)}>ABS</div>
        </Tooltip>
        <Tooltip title='NEG(field_name) - inverts the sign of a number' placement="bottom">
            <div className="textareaEditorButton" onClick={(e) => addText(e, txtId, callback, 'NEG(number)', brackets)}>NEG</div>
        </Tooltip>
        <Tooltip title='SIGN(field_name) - returns +1, 0 or -1 depending on if the number is positive, zero or negative' placement="bottom">
            <div className="textareaEditorButton" onClick={(e) => addText(e, txtId, callback, 'SIGN(number)', brackets)}>SIGN</div>
        </Tooltip>
        <Tooltip title='STRING(field_name) - converts number to string' placement="bottom">
            <div className="textareaEditorButton" onClick={(e) => addText(e, txtId, callback, 'STRING(number)', brackets)}>STRING</div>
        </Tooltip>
        <Tooltip title='INT(field_name) - returns an integer part of a decimal number' placement="bottom">
            <div className="textareaEditorButton" onClick={(e) => addText(e, txtId, callback, 'INT(number)', brackets)}>INT</div>
        </Tooltip>
        <div className="note float-left px-2 pt-3">{"Single number functions"}</div>
    </div>
    { pipelineType === PIPELINE_TYPE_METRIC &&
        <div className="clearfix">
            <Tooltip title='METRIC_AVG("period", duration) - average value of this metric within a selected period of time: second, minute, hour, day, week, month or year' placement="bottom">
                <div className="textareaEditorButton" onClick={(e) => addText(e, txtId, callback, 'METRIC_AVG("day", 5)', brackets)}>METRIC_AVG</div>
            </Tooltip>
            <Tooltip title='METRIC_QUANTILE(level, "period", duration) - metric quantile within a selected period of time: second, minute, hour, day, week, month or year' placement="bottom">
                <div className="textareaEditorButton" onClick={(e) => addText(e, txtId, callback, 'METRIC_QUANTILE(0.25, "day", 7)', brackets)}>METRIC_QUANTILE</div>
            </Tooltip>
            <Tooltip title='METRIC_MEDIAN("period", duration) - median value of this metric within a selected period of time: second, minute, hour, day, week, month or year' placement="bottom">
                <div className="textareaEditorButton" onClick={(e) => addText(e, txtId, callback, 'METRIC_MEDIAN("day", 5)', brackets)}>METRIC_MEDIAN</div>
            </Tooltip>
            <div className="note float-left px-2 pt-3">{"Metric functions"}</div>
        </div>
    }
    <div className="clearfix">
        <Tooltip title="ENV_variable_name - use environment variable" placement="bottom">
            <div className="textareaEditorButton" onClick={(e) => addText(e, txtId, callback, "ENV_", brackets)}>ENV</div>
        </Tooltip>
        <Tooltip title="NOW() - current timestamp" placement="bottom">
            <div className="textareaEditorButton" onClick={(e) => addText(e, txtId, callback, "NOW()", brackets)}>NOW</div>
        </Tooltip>
        <Tooltip title="INPUT() - the entire input data set as JSON string" placement="bottom">
            <div className="textareaEditorButton" onClick={(e) => addText(e, txtId, callback, "INPUT()", brackets)}>INPUT</div>
        </Tooltip>
        <div className="note float-left px-2 pt-3">{"Other functions"}</div>
    </div>
  </Fragment>
);
