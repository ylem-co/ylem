import React from 'react'

import Navbar from 'react-bootstrap/Navbar';
import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';

import Breadcrumbs from "../components/breadcrumbs.component";

import {HeaderDropdown, HeaderSearch} from './index'

const Header = () => {
  return (
    <>
      <Navbar className="navbarTop">
          <Row className="w-100">
            <Col className="headerSearch">
              <HeaderSearch/>
            </Col>
            <Col className="profileCol">
              <HeaderDropdown/>
            </Col>
          </Row>
      </Navbar>
      <div className="pt-3 px-4">
        <Breadcrumbs/>
      </div>
    </>
  )
}

export default Header
