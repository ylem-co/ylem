import React, { Component } from 'react'
import { Navigate } from 'react-router-dom';
import { connect } from "react-redux";
import { logout } from "../actions/auth";
import { validatePermissions, PERMISSION_LOGGED_IN }   from "../actions/pipeline";

import { Dropdown } from 'react-bootstrap';
import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';

import Avatar from "../components/avatar.component";

class HeaderDropdown extends Component {
  constructor(props) {
    super(props);
    this.logOut = this.logOut.bind(this);
  }

  logOut() {
    this.props.dispatch(logout());
  }

  render() {
    const {user, isLoggedIn} = this.props;

    if (!validatePermissions(isLoggedIn, user, PERMISSION_LOGGED_IN)) {
      return <Navigate to="/login"/>;
    }

    return (
      <Row>
        <Col sm={9}></Col>
        <Col sm={3} className="profileDropdown">
          <Dropdown
              className="c-header-nav-items mx-5 buttonDropdown"
          >
            <Dropdown.Toggle 
              className="c-header-nav-link p-0" 
              variant="white"
            >
                <Avatar
                    email={user.email}
                    avatar_url={null}
                    size={24}
                />
            </Dropdown.Toggle>
            <Dropdown.Menu className="pt-0 settingsDropdown">
               <Dropdown.Item
                  tag="div"
                  color="light"
                  className="text-center"
              >
                <strong>{user.first_name}</strong>
              </Dropdown.Item>
              <Dropdown.Item href='/settings'>
                Settings
              </Dropdown.Item>
              <Dropdown.Item onClick={this.logOut}>
                Logout
              </Dropdown.Item>
            </Dropdown.Menu>
          </Dropdown>
        </Col>
      </Row>
    )
  }
}

function mapStateToProps(state) {
  const { user } = state.auth;
  const { isLoggedIn } = state.auth;
  return {
    user,
    isLoggedIn,
  };
}

export default connect(mapStateToProps)(HeaderDropdown);
