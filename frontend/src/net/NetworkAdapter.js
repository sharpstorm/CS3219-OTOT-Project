import CardModel from '../model/CardModel';

let dummyData = [
  new CardModel(1, 'swsh4-23', 'Charmander', 'https://images.pokemontcg.io/swsh4/23_hires.png'),
  new CardModel(2, 'swsh4-24', 'Charmelon', 'https://images.pokemontcg.io/swsh4/24_hires.png'),
  new CardModel(3, 'swsh4-25', 'Charizard', 'https://images.pokemontcg.io/swsh4/25_hires.png'),
  new CardModel(4, 'swsh4-36', 'Galarian Darmantian', 'https://images.pokemontcg.io/swsh4/36_hires.png'),
];
let nextIndex = 5;

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

    this.changeCallback([...dummyData]);
  }

  async netCreateCard(newCardModel) {
    if (!this.changeCallback) {
      return;
    }

    dummyData.push(new CardModel(
      nextIndex,
      newCardModel.cardUniqueId,
      newCardModel.name,
      newCardModel.imageUrl,
    ));
    nextIndex += 1;

    await this.netGetCards();
  }

  async netEditCard(newCardModel) {
    if (!this.changeCallback) {
      return;
    }

    dummyData = dummyData.map(
      (card) => ((card.cardId === newCardModel.cardId) ? newCardModel : card),
    );

    await this.netGetCards();
  }

  async netDeleteCard(cardModel) {
    if (!this.changeCallback) {
      return;
    }

    dummyData = dummyData.filter((card) => (card.cardId !== cardModel.cardId));

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
