import * as React from 'react';
import * as ReactDOM from 'react-dom/client';
import App from './App';
import reportWebVitals from './reportWebVitals';
import 'bootstrap/dist/css/bootstrap.css';
import './index.css';

const container = document.getElementById('root');
if (!container) {
  throw new Error('Failed to find the root element');
}
const root = ReactDOM.createRoot(container);
const app = React.createElement(App);

root.render(app);

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();
