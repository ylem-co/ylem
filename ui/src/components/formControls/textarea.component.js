import React, { Fragment } from "react";
import PropTypes from "prop-types";
import { control } from "react-validation";
import TextareaAutosize from 'react-textarea-autosize';
 
const Textarea = ({ error, isChanged, isUsed, ...props }) => (
  <Fragment>
    <TextareaAutosize minRows="5" maxRows="30" {...props} />
    {isChanged && isUsed && error}
  </Fragment>
);

Textarea.propTypes = {
  error: PropTypes.oneOfType([PropTypes.node, PropTypes.string])
};

export default control(Textarea);
