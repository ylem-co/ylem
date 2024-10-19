import React, { Fragment } from "react";
import PropTypes from "prop-types";
import { control } from "react-validation";
 
const Input = ({ error, isChanged, isUsed, ...props }) => (
  <Fragment>
    <input {...props} />
    {isChanged && isUsed && error}
  </Fragment>
);

Input.propTypes = {
  error: PropTypes.oneOfType([PropTypes.node, PropTypes.string])
};

export default control(Input);
