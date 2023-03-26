import Player from './player'
import Answer from './answer'

interface Game {
  id: string
  quizID: string

  code: string
  playerLimit: string

  currentQuestion: string
  currentDeadline: string

  players: Player[]
  answers: Answer[]

  startTime: string
  finishTime: string
}

export default Game
