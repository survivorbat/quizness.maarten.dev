import Player from './player'
import Game from './game'

interface Answer {
  playerID: string
  player: Player

  gameID: string
  game: Game

  questionID: string
  optionID: string
}

export default Answer
