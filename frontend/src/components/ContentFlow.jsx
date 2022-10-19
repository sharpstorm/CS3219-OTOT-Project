import React, { useState } from 'react';
import HeaderBar from './HeaderBar';
import ContentPane from './ContentPane';
import AddCardModal from './modals/AddCardModal';
import ViewCardModal from './modals/ViewCardModal';
import { useNetworkAdapterContext } from './CardDataProvider';
import { useToast } from './ToastProvider';

const MODAL_NONE = 0;
const MODAL_CREATE = 1;
const MODAL_VIEW = 2;

export default function ContentFlow() {
  const [activeModal, setActiveModal] = useState(MODAL_NONE);
  const [isSaving, setIsSaving] = useState(false);
  const [cardData, setCardData] = useState(undefined);
  const networkAdapter = useNetworkAdapterContext();
  const pushToast = useToast();

  const viewCard = (cardModel) => {
    setCardData(cardModel);
    setActiveModal(MODAL_VIEW);
  };

  const checkPriceForCard = (cardModel) => {
    networkAdapter.netCheckPrice(cardModel).catch(() => {
      pushToast('Failed to get card price');
    });
  };

  const saveCreateCard = (newCard) => {
    setIsSaving(true);
    networkAdapter.netCreateCard(newCard).then(() => {
      setIsSaving(false);
      setActiveModal(MODAL_NONE);
    }).catch((err) => {
      pushToast(`Failed to save card, error: ${err}`);
    });
  };

  const saveEditCard = (newCard) => {
    setIsSaving(true);
    networkAdapter.netEditCard(newCard).then(() => {
      setIsSaving(false);
      setActiveModal(MODAL_NONE);
    }).catch((err) => {
      pushToast(`Failed to save card changes, error: ${err}`);
    });
  };

  const deleteCard = (card) => {
    setIsSaving(true);
    networkAdapter.netDeleteCard(card).then(() => {
      setIsSaving(false);
      setActiveModal(MODAL_NONE);
    }).catch((err) => {
      pushToast(`Failed to delete card, error: ${err}`);
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
        onPriceCheck={checkPriceForCard}
      />
    </>
  );
}
