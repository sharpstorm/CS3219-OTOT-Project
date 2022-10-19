import React, { useState } from 'react';
import { InputGroup, Form, Button } from 'react-bootstrap';
import PropTypes from 'prop-types';

export default function SearchBar({ searchCallback }) {
  const [searchBarInput, setSearchBarInput] = useState('');

  const formSubmitCallback = (evt) => {
    evt.preventDefault();
    searchCallback(searchBarInput);
  };

  return (
    <Form onSubmit={formSubmitCallback}>
      <InputGroup className="mb-3">
        <Form.Control
          placeholder="Search Cards..."
          onInput={(evt) => setSearchBarInput(evt.target.value)}
        />
        <Button variant="primary" id="btn-search" type="submit">
          Search
        </Button>
      </InputGroup>
    </Form>
  );
}

SearchBar.propTypes = {
  searchCallback: PropTypes.func.isRequired,
};
