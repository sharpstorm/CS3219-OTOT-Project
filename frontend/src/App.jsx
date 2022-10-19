import React from 'react';

import { CardDataProvider } from './components/CardDataProvider';
import ContentFlow from './components/ContentFlow';
import SetupFlow from './components/SetupFlow';

function App() {
  return (
    <React.StrictMode>
      <CardDataProvider>
        <SetupFlow />
        <ContentFlow />
      </CardDataProvider>
    </React.StrictMode>
  );
}

export default App;
