import React, { useState } from 'react';
import HeaderBar from './HeaderBar';
import ContentPane from './ContentPane';
import AddCardModal from './modals/AddCardModal';
import ViewCardModal from './modals/ViewCardModal';
import { useNetworkAdapterContext } from './CardDataProvider';

const MODAL_NONE = 0;
const MODAL_CREATE = 1;
const MODAL_VIEW = 2;

export default function ContentFlow() {
  const [activeModal, setActiveModal] = useState(MODAL_NONE);
  const [isSaving, setIsSaving] = useState(false);
  const [cardData, setCardData] = useState(undefined);
  const networkAdapter = useNetworkAdapterContext();

  const viewCard = (cardModel) => {
    setCardData(cardModel);
    setActiveModal(MODAL_VIEW);
  };

  const saveCreateCard = (newCard) => {
    setIsSaving(true);
    networkAdapter.netCreateCard(newCard).then(() => {
      setIsSaving(false);
      setActiveModal(MODAL_NONE);
    });
  };

  const saveEditCard = (newCard) => {
    setIsSaving(true);
    networkAdapter.netEditCard(newCard).then(() => {
      setIsSaving(false);
      setActiveModal(MODAL_NONE);
    });
  };

  const deleteCard = (card) => {
    setIsSaving(true);
    networkAdapter.netDeleteCard(card).then(() => {
      setIsSaving(false);
      setActiveModal(MODAL_NONE);
    });
  };

  return (
    <>
      <AddCardModal
        show={activeModal === MODAL_CREATE}
        onSave={saveCreateCard}
        onClose={() => setActiveModal(MODAL_NONE)}
        isSaving={isSaving}
      />
      <ViewCardModal
        show={activeModal === MODAL_VIEW}
        onSave={saveEditCard}
        onClose={() => setActiveModal(MODAL_NONE)}
        isSaving={isSaving}
        cardData={cardData}
        onEdit={saveEditCard}
        onDelete={deleteCard}

      />
      <HeaderBar
        addListener={() => setActiveModal(MODAL_CREATE)}
      />
      <ContentPane
        onCardSelected={viewCard}
      />
    </>
  );
}
