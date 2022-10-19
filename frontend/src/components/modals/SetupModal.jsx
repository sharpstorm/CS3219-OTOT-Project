import React, { useState } from 'react';
import {
  Modal,
  Button,
  Alert,
  Spinner,
} from 'react-bootstrap';
import PropTypes from 'prop-types';
import './ModalExtras.css';
import ModalFormTextRow from './ModalFormTextRow';

export default function SetupModal({
  show,
  onSetup,
  isLoadButtonDisabled,
  errorMessage,
}) {
  const [apiKeyInput, setApiKeyInput] = useState('');
  const onSubmit = () => {
    const apiKey = apiKeyInput;
    onSetup(apiKey);
  };

  return (
    <Modal show={show} backdrop="static" dialogClassName="setup-modal">
      <Modal.Header>
        <Modal.Title>API Key Setup</Modal.Title>
      </Modal.Header>
      <Modal.Body>
        <Alert variant="danger">
          { (errorMessage.length > 0) ? errorMessage : 'An API Key is required to continue!' }
        </Alert>
        <ModalFormTextRow
          label="API Key"
          placeholder="Enter API Key"
          value={apiKeyInput}
          onChange={(evt) => setApiKeyInput(evt.target.value)}
        />
      </Modal.Body>
      <Modal.Footer>
        <Button variant="primary" onClick={onSubmit} disabled={isLoadButtonDisabled}>
          {isLoadButtonDisabled && (
            <Spinner
              as="span"
              animation="border"
              size="sm"
              role="status"
              aria-hidden="true"
            />
          )}
          <span style={{ marginLeft: '8px' }}>Load API Key</span>
        </Button>
      </Modal.Footer>
    </Modal>
  );
}

SetupModal.defaultProps = {
  errorMessage: '',
};

SetupModal.propTypes = {
  show: PropTypes.bool.isRequired,
  onSetup: PropTypes.func.isRequired,
  isLoadButtonDisabled: PropTypes.bool.isRequired,
  errorMessage: PropTypes.string,
};
