import PropTypes from 'prop-types';
import React, { useContext, useState, useMemo } from 'react';
import { ToastContainer, Toast } from 'react-bootstrap';

const ToastProviderContext = React.createContext(() => {});

function ToastProvider({ children }) {
  const [notifications, setNotifications] = useState({});

  const pushNotification = useMemo(() => ((notification) => {
    const newNotifs = { ...notifications };
    newNotifs[new Date().getTime()] = notification;
    setNotifications(newNotifs);
  }), []);

  const removeNotification = (key) => {
    setNotifications((curNotifs) => {
      const copy = { ...curNotifs };
      delete copy[key];
      return copy;
    });
  };

  return (
    <>
      <ToastContainer className="p-3" position="top-center" containerPosition="fixed">
        {Array.from(Object.keys(notifications)).filter((x) => x !== undefined).map((key) => (
          <Toast
            key={key}
            onClose={() => removeNotification(key)}
            show
            delay={2000}
            autohide
          >
            <Toast.Header closeButton={false}>Notification</Toast.Header>
            <Toast.Body>{notifications[key]}</Toast.Body>
          </Toast>
        ))}
      </ToastContainer>
      <ToastProviderContext.Provider value={pushNotification}>
        {children}
      </ToastProviderContext.Provider>
    </>
  );
}

ToastProvider.defaultProps = {
  children: [],
};

ToastProvider.propTypes = {
  children: PropTypes.oneOfType([
    PropTypes.arrayOf(PropTypes.node),
    PropTypes.node,
  ]),
};

function useToast() {
  return useContext(ToastProviderContext);
}

export {
  ToastProvider,
  useToast,
};
