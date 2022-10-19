import React, {
  useContext,
  useMemo,
  useState,
} from 'react';
import PropTypes from 'prop-types';
import NetworkAdapter from '../net/NetworkAdapter';

const CardCacheContext = React.createContext([]);
const NetworkAdapterContext = React.createContext(undefined);

function CardDataProvider({
  children,
}) {
  const [cardCache, setCardCache] = useState([]);
  const networkAdapter = useMemo(() => new NetworkAdapter(
    (data) => setCardCache(data),
  ), []);
  return (
    <CardCacheContext.Provider value={cardCache}>
      <NetworkAdapterContext.Provider value={networkAdapter}>
        {children}
      </NetworkAdapterContext.Provider>
    </CardCacheContext.Provider>
  );
}

CardDataProvider.defaultProps = {
  children: [],
};

CardDataProvider.propTypes = {
  children: PropTypes.oneOfType([
    PropTypes.arrayOf(PropTypes.node),
    PropTypes.node,
  ]),
};

function useCardCacheContext() {
  return useContext(CardCacheContext);
}

function useNetworkAdapterContext() {
  return useContext(NetworkAdapterContext);
}
export {
  CardDataProvider,
  CardCacheContext,
  NetworkAdapterContext,
  useCardCacheContext,
  useNetworkAdapterContext,
};
