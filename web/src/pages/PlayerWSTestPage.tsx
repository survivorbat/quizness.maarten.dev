import BackendSdk from "../logic/sdk";
import {useEffect, useState} from "react";
import {BroadcastParticipant, BroadcastState} from "../models/broadcast-message";
import PlayerGameClient, {GameCallbacks} from "../logic/player-game-client";
import {useParams} from "react-router-dom";

interface PlayerWSTestPageProps {
  sdk: BackendSdk;
}

function PlayerWSTestPage({sdk}: PlayerWSTestPageProps) {
  const {player, game} = useParams();

  const [client, setClient] = useState(undefined as PlayerGameClient | undefined);
  const [players, setPlayers] = useState([] as BroadcastParticipant[]);
  const [creator, setCreator] = useState({} as BroadcastParticipant);

  useEffect(() => {
    const callbacks: GameCallbacks = {
      state(state: BroadcastState) {
        setPlayers(state.players);
        setCreator(state.creator || {});
      },
      close() {
        alert('disconnected');
      },
      error: console.error,
    }

    const playerClient = sdk.getPlayerClient(game!, player!, callbacks);
    playerClient.connect();
    setClient(playerClient);

    return () => {
      playerClient.close();
    }
  }, []);

  return <div>
    Players: {players.map((player) => player.nickname)} <br/>
    Creator: {creator.nickname}
  </div>
}

export default PlayerWSTestPage;