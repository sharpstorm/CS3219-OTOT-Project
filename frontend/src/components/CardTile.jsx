import React from 'react';
import {
  Button,
  Card,
  Col,
  ListGroup,
} from 'react-bootstrap';
import PropTypes from 'prop-types';

export default function CardTile({
  title,
  imageUrl,
  priceData,
  onView,
  onPriceCheck,
}) {
  const descriptionComponent = (
    priceData === undefined
    || priceData.data === undefined
    || priceData.data.length <= 0
  ) ? undefined : (
      priceData.data.map((x) => (
        <React.Fragment key={x.name}>
          <Card.Text><strong>{`Price for ${x.name}`}</strong></Card.Text>
          <ListGroup variant="flush">
            <ListGroup.Item key="market">{`Current Market: $${x.market}`}</ListGroup.Item>
            <ListGroup.Item key="high">{`Highest: $${x.high}`}</ListGroup.Item>
            <ListGroup.Item key="mid">{`Average: $${x.mid}`}</ListGroup.Item>
            <ListGroup.Item key="low">{`Lowest: $${x.low}`}</ListGroup.Item>
          </ListGroup>
          <hr />
        </React.Fragment>
      ))
    );

  const footerComponent = (priceData === undefined
    || priceData.updatedAt === undefined
  ) ? undefined : (
    <Card.Footer>
      {`Prices updated at ${priceData.updatedAt}`}
    </Card.Footer>
    );

  return (
    <Col>
      <Card border="secondary" style={{ marginTop: 8 }}>
        <Card.Img variant="top" src={imageUrl} />
        <Card.Body>
          <Card.Title>
            {title}
          </Card.Title>
          <hr />
          {descriptionComponent}
          <div className="text-center">
            <Button variant="primary" onClick={onView} style={{ marginRight: 8 }}>View</Button>
            <Button variant="secondary" onClick={onPriceCheck}>Check Price</Button>
          </div>
        </Card.Body>
        {footerComponent}
      </Card>
    </Col>
  );
}

CardTile.defaultProps = {
  priceData: {
    updatedAt: undefined,
    data: [],
  },
};

CardTile.propTypes = {
  title: PropTypes.string.isRequired,
  imageUrl: PropTypes.string.isRequired,
  priceData: PropTypes.shape({
    updatedAt: PropTypes.string,
    data: PropTypes.arrayOf(PropTypes.shape({
      name: PropTypes.string.isRequired,
      high: PropTypes.number.isRequired,
      mid: PropTypes.number.isRequired,
      low: PropTypes.number.isRequired,
      market: PropTypes.number.isRequired,
    })),
  }),
  onView: PropTypes.func.isRequired,
  onPriceCheck: PropTypes.func.isRequired,
};
