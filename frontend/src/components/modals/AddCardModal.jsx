import React from 'react';
import PropTypes from 'prop-types';
import CardDetailsModal from './CardDetailsModal';
import CardModel from '../../model/CardModel';

export default function AddCardModal({
  show,
  onSave,
  onClose,
  isSaving,
}) {
  const saveCard = (content) => {
    const { cardUniqueId, cardName, cardImageUrl } = content;
    onSave(new CardModel(0, cardUniqueId, cardName, cardImageUrl));
  };

  const actions = [
    {
      actionName: 'Create Card',
      actionListener: saveCard,
      isActionLoading: isSaving,
      buttonVariant: 'primary',
    },
  ];

  return (
    <CardDetailsModal
      show={show}
      title="Create New Card"
      actions={actions}
      onClose={onClose}
      isCloseDisabled={isSaving}
      isInputDisabled={isSaving}
    />
  );
}

AddCardModal.propTypes = {
  show: PropTypes.bool.isRequired,
  onSave: PropTypes.func.isRequired,
  onClose: PropTypes.func.isRequired,
  isSaving: PropTypes.bool.isRequired,
};
