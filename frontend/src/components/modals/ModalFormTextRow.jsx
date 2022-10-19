import React from 'react';
import PropTypes from 'prop-types';
import { Form } from 'react-bootstrap';
import ModalFormRow from './ModalFormRow';

export default function ModalFormTextRow({
  label,
  placeholder,
  value,
  onChange,
  disabled,
}) {
  return (
    <ModalFormRow label={label}>
      <Form.Control
        type="text"
        placeholder={placeholder}
        value={value}
        onInput={onChange}
        disabled={disabled}
      />
    </ModalFormRow>
  );
}

ModalFormTextRow.propTypes = {
  label: PropTypes.string.isRequired,
  placeholder: PropTypes.string.isRequired,
  value: PropTypes.string.isRequired,
  onChange: PropTypes.func.isRequired,
  disabled: PropTypes.bool,
};

ModalFormTextRow.defaultProps = {
  disabled: false,
};
