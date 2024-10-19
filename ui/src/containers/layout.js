import React from 'react'
import Container from 'react-bootstrap/Container';
import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';

import {
  Header,
  Sidebar,
  Content
} from './index';

const Layout = () => {

  return (
      <>
          <Container fluid className="d-inline-block mainContainer">
            <Row className="h-100">
              <Col className="px-0 mainSidebar">
                <Sidebar/>
              </Col>
              <Col className="px-0">
                  <Header/>
                  <Content/>
              </Col>
            </Row>
          </Container>
      </>
  )
}

export default Layout
