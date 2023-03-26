import GameCodeInput from '../components/GameCodeInput'
import Player from '../models/player'
import { useNavigate } from 'react-router-dom'

interface FrontPageProps {
  codeSubmitCallback: (code: string) => Promise<Player | null>
}

function FrontPage({ codeSubmitCallback }: FrontPageProps) {
  const navigate = useNavigate()

  const joinGame = async (code: string): Promise<boolean> => {
    const result = await codeSubmitCallback(code)
    if (result != null) {
      navigate(`/games/${result.gameID}/players/${result.id}`)
    }

    return true
  }

  return (
    <div>
      <GameCodeInput callback={joinGame} />
    </div>
  )
}

export default FrontPage
