import React from 'react';
import PropTypes from 'prop-types';
import { Form } from 'react-bootstrap';

export default function ModalFormRow({
  label,
  children,
}) {
  return (
    <Form.Group className="mb-3">
      <Form.Label>{label}</Form.Label>
      {children}
    </Form.Group>
  );
}

ModalFormRow.propTypes = {
  label: PropTypes.string.isRequired,
  children: PropTypes.oneOfType([
    PropTypes.arrayOf(PropTypes.node),
    PropTypes.node,
  ]),
};

ModalFormRow.defaultProps = {
  children: [],
};
