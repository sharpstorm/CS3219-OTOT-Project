import React from 'react';
import PropTypes from 'prop-types';
import { Row } from 'react-bootstrap';
import CardTile from './CardTile';
import CardModel from '../model/CardModel';
import 'bootstrap/dist/css/bootstrap.css';

export default function CardPane({ cards, onCardSelected }) {
  return (
    <Row xs={1} md={2} xl={4}>
      {cards.map((cardModel) => (
        <CardTile
          key={cardModel.cardId.toString()}
          title={cardModel.name}
          imageUrl={cardModel.imageUrl}
          onClick={() => onCardSelected(cardModel)}
        />
      ))}
    </Row>
  );
}

CardPane.propTypes = {
  cards: PropTypes.arrayOf(PropTypes.instanceOf(CardModel)).isRequired,
  onCardSelected: PropTypes.func.isRequired,
};
