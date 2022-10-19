import React from 'react';

import { CardDataProvider } from './components/CardDataProvider';
import ContentFlow from './components/ContentFlow';
import SetupFlow from './components/SetupFlow';
import { ToastProvider } from './components/ToastProvider';

function App() {
  return (
    <React.StrictMode>
      <CardDataProvider>
        <ToastProvider>
          <SetupFlow />
          <ContentFlow />
        </ToastProvider>
      </CardDataProvider>
    </React.StrictMode>
  );
}

export default App;
