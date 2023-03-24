import {Chip} from "@mui/material";
import {BroadcastParticipant} from "../models/broadcast-message";
import Creator from "../models/creator";
import Player from "../models/player";
import React, {useState} from "react";

interface PlayerDotProps {
  participant:BroadcastParticipant | Player | Creator;
}

function ParticipantDot({participant}: PlayerDotProps) {
  const [expanded, setExpanded] = useState(false);

  const initials = participant?.nickname?.split(' ').map((l) => l[0].toUpperCase()).join('');

  const style = {
    backgroundColor: participant.backgroundColor,
    color: participant.color,
  }

  return <Chip
    onMouseEnter={() => setExpanded(true)}
    onMouseLeave={() => setExpanded(false)}
    label={expanded ? participant.nickname : initials}
    style={style}
  />
}

export default ParticipantDot;
