import {BroadcastMessage} from "../models/broadcast-message";

export default class CreatorGameClient {
  private socket?: WebSocket;

  constructor(private readonly baseUrl: string, private readonly token: string, private readonly gameID: string) {
  }

  connect(onMessage: (message: BroadcastMessage) => void, onClose: () => void) {
    if (this.socket) {
      throw new Error("already connected");
    }

    this.socket = new WebSocket(`${this.baseUrl}/api/v1/games/${this.gameID}/connection`);
    this.socket.onmessage = ((event: MessageEvent<BroadcastMessage>) => onMessage(event.data))
    this.socket.onclose = onClose
  }

  next() {
    if (!this.socket) {
      throw new Error("not connected");
    }

    this.socket.send(JSON.stringify({action: 'next'}));
  }

  finish() {
    if (!this.socket) {
      throw new Error("not connected");
    }

    this.socket.send(JSON.stringify({action: 'finish'}));
  }
}