export default class CardModel {
  cardId;

  cardUniqueId;

  name;

  imageUrl;

  constructor(cardId, cardUniqueId, name, imageUrl) {
    this.cardId = cardId;
    this.cardUniqueId = cardUniqueId;
    this.name = name;
    this.imageUrl = imageUrl;
  }
}
