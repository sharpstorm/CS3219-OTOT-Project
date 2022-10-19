import React, { useState } from 'react';
import { useNetworkAdapterContext } from './CardDataProvider';
import SetupModal from './modals/SetupModal';

export default function SetupFlow() {
  const [isSetup, setIsSetup] = useState(false);
  const [isSetupLoading, setIsSetupLoading] = useState(false);
  const [setupErrorMessage, setSetupErrorMessage] = useState('');
  const networkAdapter = useNetworkAdapterContext();

  const loadApiKey = (key) => {
    networkAdapter.setApiKey(key);
    setIsSetupLoading(true);
    networkAdapter.netGetCards().then(() => {
      setIsSetup(true);
    }).catch(() => {
      setSetupErrorMessage('Either server cannot be reached, or API key is incorrect');
      setIsSetupLoading(false);
    });
  };

  return (
    <SetupModal
      show={!isSetup}
      onSetup={loadApiKey}
      isLoadButtonDisabled={isSetupLoading}
      errorMessage={setupErrorMessage}
    />
  );
}
