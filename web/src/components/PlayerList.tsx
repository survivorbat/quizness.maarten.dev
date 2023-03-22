import {BroadcastParticipant} from "../models/broadcast-message";
import Player from "../models/player";
import {Fragment} from "react";
import ParticipantDot from "./ParticipantDot";

interface PlayerListProps {
  players: BroadcastParticipant[] | Player[];
}

function PlayerList({players}: PlayerListProps) {
  return <Fragment>
    {players.map((player) => <ParticipantDot key={player.id} participant={player} />)}
  </Fragment>
}

export default PlayerList;
