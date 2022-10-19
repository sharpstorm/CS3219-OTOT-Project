import React from 'react';
import PropTypes from 'prop-types';
import {
  Navbar,
  Nav,
  Container,
} from 'react-bootstrap';

export default function HeaderBar({
  addListener,
}) {
  return (
    <Navbar bg="dark" variant="dark" expand="md" sticky="top">
      <Container>
        <Navbar.Brand href="#">Pokemon TCG DB</Navbar.Brand>
        <Navbar.Toggle aria-controls="basic-navbar-nav" />
        <Navbar.Collapse id="basic-navbar-nav">
          <Nav className="me-auto">
            <Nav.Link onClick={addListener}>Add Card</Nav.Link>
          </Nav>
        </Navbar.Collapse>
      </Container>
    </Navbar>
  );
}

HeaderBar.propTypes = {
  addListener: PropTypes.func.isRequired,
};
