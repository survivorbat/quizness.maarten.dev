import {BroadcastMessage} from "../models/broadcast-message";

export default class PlayerGameClient {
  private socket?: WebSocket;

  constructor(private readonly baseUrl: string, private readonly gameID: string, private readonly playerID: string) {
  }

  connect(onMessage: (message: BroadcastMessage) => void, onClose: () => void) {
    if (this.socket) {
      throw new Error("already connected");
    }

    this.socket = new WebSocket(`${this.baseUrl}/api/v1/games/${this.gameID}/players/${this.playerID}/connection`);
    this.socket.onmessage = ((event: MessageEvent<BroadcastMessage>) => onMessage(event.data))
    this.socket.onclose = onClose
  }
}