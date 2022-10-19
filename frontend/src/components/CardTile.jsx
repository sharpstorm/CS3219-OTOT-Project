import React from 'react';
import { Button, Card, Col } from 'react-bootstrap';
import PropTypes from 'prop-types';

export default function CardTile({
  title,
  imageUrl,
  description,
  footer,
  onClick,
}) {
  const descriptionComponent = (description === '') ? undefined : (
    <Card.Text>
      {description}
    </Card.Text>
  );

  const footerComponent = (footer === '') ? undefined : (
    <Card.Footer>
      {footer}
    </Card.Footer>
  );

  return (
    <Col>
      <Card border="secondary">
        <Card.Img variant="top" src={imageUrl} />
        <Card.Body>
          <Card.Title className="text-center" style={{ display: 'flex' }}>
            <Button variant="primary" onClick={onClick}>View</Button>
            <div style={{
              flex: '1 1 0',
              justifyContent: 'center',
              display: 'flex',
              flexDirection: 'column',
            }}
            >
              <span>{title}</span>
            </div>
          </Card.Title>
          {descriptionComponent}
        </Card.Body>
        {footerComponent}
      </Card>
    </Col>
  );
}

CardTile.defaultProps = {
  description: '',
  footer: '',
};

CardTile.propTypes = {
  title: PropTypes.string.isRequired,
  imageUrl: PropTypes.string.isRequired,
  description: PropTypes.string,
  footer: PropTypes.string,
  onClick: PropTypes.func.isRequired,
};
