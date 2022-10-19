import React from 'react';
import PropTypes from 'prop-types';
import { Row } from 'react-bootstrap';
import CardTile from './CardTile';
import CardModel from '../model/CardModel';
import 'bootstrap/dist/css/bootstrap.css';
import { useCardPriceCacheContext } from './CardDataProvider';

export default function CardPane({ cards, onCardSelected, onPriceCheck }) {
  const priceCache = useCardPriceCacheContext();

  return (
    <Row xs={1} md={2} xl={4}>
      {cards.map((cardModel) => {
        let priceData;
        if (cardModel.cardUniqueId in priceCache) {
          const cachedData = priceCache[cardModel.cardUniqueId];
          const priceArray = Array.from(Object.keys(cachedData.prices))
            .map((key) => ({
              name: key,
              ...cachedData.prices[key],
            }));
          priceData = {
            updatedAt: cachedData.updatedAt,
            data: priceArray,
          };
        }

        return (
          <CardTile
            key={cardModel.cardId.toString()}
            title={cardModel.name}
            imageUrl={cardModel.imageUrl}
            onView={() => onCardSelected(cardModel)}
            onPriceCheck={() => onPriceCheck(cardModel)}
            priceData={priceData}
          />
        );
      })}
    </Row>
  );
}

CardPane.propTypes = {
  cards: PropTypes.arrayOf(PropTypes.instanceOf(CardModel)).isRequired,
  onCardSelected: PropTypes.func.isRequired,
  onPriceCheck: PropTypes.func.isRequired,
};
