import CardModel from '../model/CardModel';

const priceCheckUrl = 'https://asia-southeast1-cs3219-otot-b-363213.cloudfunctions.net/cs3219-otot-b-serverless/GetPrice';

class NetworkAdapter {
  apiKey;

  changeCallback;

  priceCacheSetter;

  constructor(changeCallback, priceCacheSetter) {
    this.apiKey = '';
    this.changeCallback = changeCallback;
    this.priceCacheSetter = priceCacheSetter;
  }

  isReady() {
    return this.apiKey !== '';
  }

  setApiKey(apiKey) {
    // eslint-disable-next-line no-console
    console.log('API Key was set');
    this.apiKey = apiKey;
  }

  async netGetCards() {
    if (!this.changeCallback) {
      return;
    }

    const response = await fetch('/api/card', {
      headers: {
        Authorization: `Bearer ${this.apiKey}`,
      },
    });
    const data = await response.json();
    this.changeCallback(data.map((x) => new CardModel(x.id, x.uniqueId, x.pokemon, x.imageUrl)));
  }

  async netCreateCard(newCardModel) {
    if (!this.changeCallback) {
      return;
    }

    await fetch('/api/card', {
      method: 'POST',
      headers: {
        Authorization: `Bearer ${this.apiKey}`,
      },
      body: JSON.stringify({
        uniqueId: newCardModel.cardUniqueId,
        pokemon: newCardModel.name,
        imageUrl: newCardModel.imageUrl,
      }),
    });

    await this.netGetCards();
  }

  async netEditCard(newCardModel) {
    if (!this.changeCallback) {
      return;
    }

    await fetch(`/api/card/${newCardModel.cardId}`, {
      method: 'PUT',
      headers: {
        Authorization: `Bearer ${this.apiKey}`,
      },
      body: JSON.stringify({
        id: newCardModel.cardId,
        uniqueId: newCardModel.cardUniqueId,
        pokemon: newCardModel.name,
        imageUrl: newCardModel.imageUrl,
      }),
    });

    await this.netGetCards();
  }

  async netDeleteCard(cardModel) {
    if (!this.changeCallback) {
      return;
    }

    await fetch(`/api/card/${cardModel.cardId}`, {
      method: 'DELETE',
      headers: {
        Authorization: `Bearer ${this.apiKey}`,
      },
    });

    await this.netGetCards();
  }

  // eslint-disable-next-line class-methods-use-this
  async netCheckPrice(cardModel) {
    if (this.priceCacheSetter === undefined) {
      return;
    }

    const result = await fetch(priceCheckUrl, {
      method: 'POST',
      body: JSON.stringify({
        cardUniqueId: cardModel.cardUniqueId,
      }),
    });
    const response = await result.json();
    this.priceCacheSetter(cardModel.cardUniqueId, response.data);
  }
}

export default NetworkAdapter;
