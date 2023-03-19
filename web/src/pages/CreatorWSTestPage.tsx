import BackendSdk from "../logic/sdk";
import {useEffect, useState} from "react";
import CreatorGameClient from "../logic/creator-game-client";
import {GameCallbacks} from "../logic/player-game-client";
import {BroadcastParticipant, BroadcastState} from "../models/broadcast-message";

interface CreatorWSTestPageProps {
  sdk: BackendSdk;
}

function CreatorWSTestPage({sdk}: CreatorWSTestPageProps) {
  const [client, setClient] = useState(undefined as CreatorGameClient | undefined);
  const [players, setPlayers] = useState([] as BroadcastParticipant[]);
  const [creator, setCreator] = useState({} as BroadcastParticipant);

  useEffect(() => {
    const callbacks: GameCallbacks = {
      state(state: BroadcastState) {
        setPlayers(state.players);
        setCreator(state.creator);
      },
      close() {
        alert('disconnected');
      },
      error: console.error,
    }

    const creatorClient = sdk.getCreatorClient(prompt('gameID')!, callbacks);
    setClient(creatorClient);
    creatorClient.connect();

    return () => {
      creatorClient.close();
    }
  }, []);

  return <div>
    Players: {players.map((player) => player.nickname)} <br/>
    Creator: {creator.nickname}
  </div>
}

export default CreatorWSTestPage;