import BackendSdk from '../logic/sdk'
import { useEffect, useState } from 'react'
import { BroadcastParticipant, BroadcastState } from '../models/broadcast-message'
import PlayerGameClient, { GameCallbacks } from '../logic/player-game-client'
import { useParams } from 'react-router-dom'
import { Quiz } from '../models/quiz'
import PlayerList from '../components/PlayerList'
import ParticipantDot from '../components/ParticipantDot'

interface PlayerWSTestPageProps {
  sdk: BackendSdk
}

function PlayerGame({ sdk }: PlayerWSTestPageProps) {
  const { player, game } = useParams()

  const [client, setClient] = useState<PlayerGameClient>()
  const [players, setPlayers] = useState<BroadcastParticipant[]>([])
  const [creator, setCreator] = useState<BroadcastParticipant>()
  const [quiz, setQuiz] = useState<Quiz>()
  const [currentQuestion, setCurrentQuestion] = useState('00000000-0000-0000-0000-000000000000')

  useEffect(() => {
    if (player === undefined || game === undefined) {
      return
    }

    const callbacks: GameCallbacks = {
      state(state: BroadcastState) {
        setPlayers(state.players)
        setCreator(state.creator)
        setCurrentQuestion(state.currentQuestion)
      },
      close() {},
      error: console.error
    }

    const playerClient = sdk.getPlayerClient(game, player, callbacks)
    playerClient.connect()
    setClient(playerClient)

    sdk.getQuizByGame(game).then(setQuiz).catch(console.error)

    return () => {
      playerClient.close()
    }
  }, [player, game, sdk])

  const pickAnswer = (id: string) => {
    client?.answer(id)
  }

  const question = quiz?.multipleChoiceQuestions?.find((q) => q.id === currentQuestion)

  if (creator == null) {
    return <span>Loading...</span>
  }

  return (
    <div>
      <p>Quiz: {quiz?.name}</p>
      <p>
        Players: <PlayerList players={players} />
      </p>
      <p>
        Creator: <ParticipantDot participant={creator} />
      </p>
      <p>Current question: {question?.title}</p>
      {question?.options.map((o) => (
        <button
          key={o.id}
          onClick={() => {
            pickAnswer(o.id)
          }}
        >
          {o.textOption}
        </button>
      ))}
    </div>
  )
}

export default PlayerGame
