import {Badge} from "@mui/material";
import {BroadcastParticipant} from "../models/broadcast-message";
import Creator from "../models/creator";
import Player from "../models/player";

interface PlayerDotProps {
  participant:BroadcastParticipant | Player | Creator;
}

function ParticipantDot({participant}: PlayerDotProps) {
	return <Badge overlap="circular" variant="dot">
    {participant.color}
  </Badge>
}

export default ParticipantDot;
