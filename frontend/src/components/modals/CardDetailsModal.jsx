import React, { useState, useEffect } from 'react';
import {
  Modal,
  Button,
  Form,
  Alert,
  Spinner,
} from 'react-bootstrap';
import PropTypes from 'prop-types';
import './ModalExtras.css';
import ModalFormTextRow from './ModalFormTextRow';

export default function CardDetailsModal({
  show,
  title,
  onClose,
  isCloseDisabled,
  actions,
  cardUniqueId,
  cardName,
  cardImageUrl,
  children,
  isInputDisabled,
}) {
  const [cardUniqueIdInput, setCardUniqueId] = useState('');
  const [cardNameInput, setCardName] = useState('');
  const [cardImageInput, setCardImage] = useState('');

  const [errorMessage, setErrorMessage] = useState('');

  const onSubmit = (action) => {
    if (cardUniqueIdInput.length === 0
      || cardNameInput.length === 0
      || cardImageInput.length === 0) {
      setErrorMessage('All input fields are required');
      return;
    }
    setErrorMessage('');
    action({
      cardUniqueId: cardUniqueIdInput,
      cardName: cardNameInput,
      cardImageUrl: cardImageInput,
    });
  };

  useEffect(() => {
    if (show) {
      setCardUniqueId(cardUniqueId);
      setCardName(cardName);
      setCardImage(cardImageUrl);
    }
  }, [show, cardUniqueId, cardName, cardImageUrl]);

  return (
    <Modal show={show} backdrop="static" dialogClassName="setup-modal">
      <Modal.Header>
        <Modal.Title>{title}</Modal.Title>
      </Modal.Header>
      <Modal.Body>
        <Form>
          {children}
          <ModalFormTextRow
            label="Card Unique ID"
            placeholder="Enter Card ID"
            value={cardUniqueIdInput}
            onChange={(evt) => setCardUniqueId(evt.target.value)}
            disabled={isInputDisabled}
          />
          <ModalFormTextRow
            label="Card Name"
            placeholder="Enter Card Name"
            value={cardNameInput}
            onChange={(evt) => setCardName(evt.target.value)}
            disabled={isInputDisabled}
          />
          <ModalFormTextRow
            label="Card Image URL"
            placeholder="Enter Image URL"
            value={cardImageInput}
            onChange={(evt) => setCardImage(evt.target.value)}
            disabled={isInputDisabled}
          />
        </Form>
        {errorMessage.length > 0 && (
          <Alert variant="danger">
            { errorMessage }
          </Alert>
        )}
      </Modal.Body>
      <Modal.Footer>
        <Button variant="secondary" onClick={onClose} disabled={isCloseDisabled}>
          Back
        </Button>
        {actions.map((action) => (
          <Button
            key={action.actionName}
            variant={action.buttonVariant}
            onClick={() => onSubmit((data) => action.actionListener(data))}
            disabled={action.isActionLoading}
          >
            {action.isActionLoading && (
              <Spinner
                as="span"
                animation="border"
                size="sm"
                role="status"
                aria-hidden="true"
              />
            )}
            <span style={{ marginLeft: '8px' }}>{action.actionName}</span>
          </Button>
        ))}
      </Modal.Footer>
    </Modal>
  );
}

CardDetailsModal.defaultProps = {
  cardUniqueId: '',
  cardName: '',
  cardImageUrl: '',
  children: [],
  actions: [],
  isInputDisabled: false,
  isCloseDisabled: false,
};

CardDetailsModal.propTypes = {
  show: PropTypes.bool.isRequired,
  title: PropTypes.string.isRequired,
  onClose: PropTypes.func.isRequired,
  isInputDisabled: PropTypes.bool,
  isCloseDisabled: PropTypes.bool,

  cardUniqueId: PropTypes.string,
  cardName: PropTypes.string,
  cardImageUrl: PropTypes.string,

  actions: PropTypes.arrayOf(PropTypes.shape({
    actionName: PropTypes.string.isRequired,
    actionListener: PropTypes.func.isRequired,
    isActionLoading: PropTypes.bool.isRequired,
    buttonVariant: PropTypes.string.isRequired,
  })),

  children: PropTypes.oneOfType([
    PropTypes.arrayOf(PropTypes.node),
    PropTypes.node,
  ]),
};
