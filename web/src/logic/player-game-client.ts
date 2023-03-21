import {BroadcastMessage, BroadcastState} from "../models/broadcast-message";
import {baseSocketUrl} from "./constants";

export interface GameCallbacks {
  state: (state: BroadcastState) => void;
  close: () => void;
  error: (ev: Event) => void;
}

export default class PlayerGameClient {
  private socket?: WebSocket;

  constructor(private readonly gameID: string, private readonly playerID: string, private readonly callbacks: GameCallbacks) {
  }

  connect() {
    if (this.socket) {
      throw new Error("already connected");
    }

    this.socket = new WebSocket(`${baseSocketUrl}/api/v1/games/${this.gameID}/players/${this.playerID}/connection`);

    this.socket.onerror = this.callbacks.error
    this.socket.onmessage = ((event: MessageEvent<string>) => {
      const message: BroadcastMessage = JSON.parse(event.data);

      switch (message.type) {
        case 'state':
          this.callbacks.state(message.stateContent);
      }
    });

    this.socket.onclose = this.callbacks.close
  }

  answer(option: string) {
    this.socket!.send(JSON.stringify({action: 'answer', content: {optionID: option}}))
  }

  close() {
    this.socket?.close();
  }
}
