import CardModel from '../model/CardModel';

let dummyData = [
  new CardModel(1, 'swsh4-23', 'Charmander', 'https://images.pokemontcg.io/swsh4/23_hires.png'),
  new CardModel(2, 'swsh4-24', 'Charmelon', 'https://images.pokemontcg.io/swsh4/24_hires.png'),
  new CardModel(3, 'swsh4-25', 'Charizard', 'https://images.pokemontcg.io/swsh4/25_hires.png'),
  new CardModel(4, 'swsh4-36', 'Galarian Darmantian', 'https://images.pokemontcg.io/swsh4/36_hires.png'),
];
let nextIndex = 5;

class NetworkAdapter {
  apiKey;

  changeCallback;

  constructor(changeCallback) {
    this.apiKey = '';
    this.changeCallback = changeCallback;
  }

  isReady() {
    return this.apiKey !== '';
  }

  setApiKey(apiKey) {
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
}

export default NetworkAdapter;
