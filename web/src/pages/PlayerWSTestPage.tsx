import BackendSdk from "../logic/sdk";
import {useEffect, useState} from "react";
import {BroadcastParticipant, BroadcastState} from "../models/broadcast-message";
import PlayerGameClient, {GameCallbacks} from "../logic/player-game-client";
import {useParams} from "react-router-dom";
import {Quiz} from "../models/quiz";
import PlayerList from "../components/PlayerList";
import ParticipantDot from "../components/ParticipantDot";

interface PlayerWSTestPageProps {
  sdk: BackendSdk;
}

function PlayerWSTestPage({sdk}: PlayerWSTestPageProps) {
  const {player, game} = useParams();

  const [client, setClient] = useState(undefined as PlayerGameClient | undefined);
  const [players, setPlayers] = useState([] as BroadcastParticipant[]);
  const [creator, setCreator] = useState({} as BroadcastParticipant);
  const [quiz, setQuiz] = useState({} as Quiz);
  const [currentQuestion, setCurrentQuestion] = useState('00000000-0000-0000-0000-000000000000');

  useEffect(() => {
    const callbacks: GameCallbacks = {
      state(state: BroadcastState) {
        setPlayers(state.players);
        setCreator(state.creator || {});
        setCurrentQuestion(state.currentQuestion);
      },
      close() {
        alert('disconnected');
      },
      error: console.error,
    }

    const playerClient = sdk.getPlayerClient(game!, player!, callbacks);
    playerClient.connect();
    setClient(playerClient);

    sdk.getQuizByGame(game!).then(setQuiz);

    return () => {
      playerClient.close();
    }
  }, [player, game, sdk]);

  const pickAnswer = (id: string) => {
    client?.answer(id);
  }

  const question = quiz?.multipleChoiceQuestions?.find((q) => q.id === currentQuestion);

  return <div>
    <p>Quiz: {quiz.name}</p>
    <PlayerList players={players}/>
    <ParticipantDot participant={creator}/>
    <p>Current question: {question?.title}</p>
    {question?.options.map((o) => <button key={o.id} onClick={() => pickAnswer(o.id)}>{o.textOption}</button>)}
  </div>
}

export default PlayerWSTestPage;