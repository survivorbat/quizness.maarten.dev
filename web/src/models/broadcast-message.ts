export interface BroadcastMessage {
  type: string;
  nextQuestionContent: {
    questionID: string;
  };
  playerAnsweredContent: {
    playerID: string;
  };
  stateContent: BroadcastState;
}

export interface BroadcastState {
  creator: BroadcastParticipant;
  players: BroadcastParticipant[];
  currentQuestion: string;
  currentDeadline: Date;
}

export interface BroadcastParticipant {
  id: string;
  nickname: string;
}