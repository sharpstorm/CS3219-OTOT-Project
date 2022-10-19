import React, { useEffect, useState } from 'react';
import PropTypes from 'prop-types';
import CardDetailsModal from './CardDetailsModal';
import CardModel from '../../model/CardModel';

export default function ViewCardModal({
  show,
  onEdit,
  onDelete,
  onClose,
  isSaving,
  cardData,
}) {
  const [isEditMode, setIsEditMode] = useState(false);
  const [isDeleteMode, setIsDeleteMode] = useState(false);

  useEffect(() => {
    setIsEditMode(false);
    setIsDeleteMode(false);
  }, [show]);

  const enterEditMode = () => setIsEditMode(true);
  const saveEdits = (content) => {
    const { cardUniqueId, cardName, cardImageUrl } = content;
    onEdit(new CardModel(cardData.cardId, cardUniqueId, cardName, cardImageUrl));
  };

  const enterDeleteMode = () => setIsDeleteMode(true);
  const confirmDelete = () => {
    onDelete(cardData);
  };

  const closeInterceptor = () => {
    if (isEditMode) {
      setIsEditMode(false);
    } else if (isDeleteMode) {
      setIsDeleteMode(false);
    } else {
      onClose();
    }
  };

  const editAction = {
    actionName: isEditMode ? 'Save Card' : 'Edit Card',
    actionListener: isEditMode ? saveEdits : enterEditMode,
    isActionLoading: isSaving,
    buttonVariant: 'primary',
  };

  const deleteAction = {
    actionName: isDeleteMode ? 'Confirm Delete Card' : 'Delete Card',
    actionListener: isDeleteMode ? confirmDelete : enterDeleteMode,
    isActionLoading: isSaving,
    buttonVariant: 'danger',
  };

  let actions = [editAction, deleteAction];

  if (isEditMode) {
    actions = [editAction];
  } else if (isDeleteMode) {
    actions = [deleteAction];
  }

  return (
    <CardDetailsModal
      show={show}
      title="Card Info"
      onClose={closeInterceptor}
      isCloseDisabled={isSaving}
      isInputDisabled={!isEditMode || isSaving}
      actions={actions}
      cardName={cardData.name}
      cardUniqueId={cardData.cardUniqueId}
      cardImageUrl={cardData.imageUrl}
    />
  );
}

ViewCardModal.defaultProps = {
  cardData: new CardModel(0, '', '', ''),
};

ViewCardModal.propTypes = {
  show: PropTypes.bool.isRequired,
  onEdit: PropTypes.func.isRequired,
  onDelete: PropTypes.func.isRequired,
  onClose: PropTypes.func.isRequired,
  isSaving: PropTypes.bool.isRequired,
  cardData: PropTypes.instanceOf(CardModel),
};
