export interface BroadcastMessage {
  type: string
  playerAnsweredContent: {
    playerID: string
  }
  stateContent: BroadcastState
}

export interface BroadcastState {
  creator: BroadcastParticipant
  players: BroadcastParticipant[]
  currentQuestion: string
  currentDeadline: Date
}

export interface BroadcastParticipant {
  id: string
  nickname: string
  color: string
  backgroundColor: string
}
