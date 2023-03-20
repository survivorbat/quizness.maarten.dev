import {BroadcastMessage} from "../models/broadcast-message";
import {baseSocketUrl} from "./constants";
import {GameCallbacks} from "./player-game-client";

export default class CreatorGameClient {
  private socket?: WebSocket;

  constructor(private readonly token: string, private readonly gameID: string, private readonly callbacks: GameCallbacks) {
  }

  connect() {
    if (this.socket) {
      throw new Error("already connected");
    }

    this.socket = new WebSocket(`${baseSocketUrl}/api/v1/games/${this.gameID}/connection`, [`Bearer_${this.token}`]);
    this.socket.onmessage = ((event: MessageEvent<string>) => {
      const message: BroadcastMessage = JSON.parse(event.data);

      switch (message.type) {
        case 'state':
          this.callbacks.state(message.stateContent);
      }
    });
    this.socket.onclose = this.callbacks.close
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

  close() {
    this.socket?.close();
  }
}