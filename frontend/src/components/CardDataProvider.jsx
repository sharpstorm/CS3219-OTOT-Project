import React, {
  useContext,
  useMemo,
  useState,
} from 'react';
import PropTypes from 'prop-types';
import NetworkAdapter from '../net/NetworkAdapter';

const CardCacheContext = React.createContext([]);
const CardPriceCacheContext = React.createContext({});
const NetworkAdapterContext = React.createContext(undefined);

function CardDataProvider({
  children,
}) {
  const [cardCache, setCardCache] = useState([]);
  const [priceCache, setPriceCache] = useState({});

  const updatePriceIntoCache = (key, data) => {
    setPriceCache((prevState) => {
      const copy = { ...prevState };
      copy[key] = data;
      return copy;
    });
  };

  const networkAdapter = useMemo(() => new NetworkAdapter(
    (data) => setCardCache(data),
    (key, data) => updatePriceIntoCache(key, data),
  ), []);
  return (
    <CardCacheContext.Provider value={cardCache}>
      <NetworkAdapterContext.Provider value={networkAdapter}>
        <CardPriceCacheContext.Provider value={priceCache}>
          {children}
        </CardPriceCacheContext.Provider>
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

function useCardPriceCacheContext() {
  return useContext(CardPriceCacheContext);
}

export {
  CardDataProvider,
  useCardCacheContext,
  useNetworkAdapterContext,
  useCardPriceCacheContext,
};
