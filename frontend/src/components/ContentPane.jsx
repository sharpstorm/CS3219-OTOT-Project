import React, { useMemo, useState } from 'react';
import PropTypes from 'prop-types';
import { Container } from 'react-bootstrap';
import { useCardCacheContext } from './CardDataProvider';
import CardPane from './CardPane';
import SearchBar from './SearchBar';

export default function ContentPane({ onCardSelected, onPriceCheck }) {
  const globalCards = useCardCacheContext();
  const [filter, setFilter] = useState('');

  const cachedVisible = useMemo(() => globalCards.filter((card) => (
    card && (
      `${card.cardId}`.toLowerCase().includes(filter)
      || card.name.toLowerCase().includes(filter)
      || card.cardUniqueId.toLowerCase().includes(filter)
    ))), [globalCards, filter]);

  const searchCallback = (searchTerm) => {
    // eslint-disable-next-line no-console
    console.log(`Searching ${searchTerm}`);
    setFilter(searchTerm.toLowerCase());
  };

  return (
    <Container className="pt-4">
      <SearchBar searchCallback={searchCallback} />
      <CardPane cards={cachedVisible} onCardSelected={onCardSelected} onPriceCheck={onPriceCheck} />
    </Container>
  );
}

ContentPane.propTypes = {
  onCardSelected: PropTypes.func.isRequired,
  onPriceCheck: PropTypes.func.isRequired,
};
