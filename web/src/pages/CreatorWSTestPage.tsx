import BackendSdk from "../logic/sdk";
import {useEffect, useState} from "react";
import CreatorGameClient from "../logic/creator-game-client";
import {GameCallbacks} from "../logic/player-game-client";
import {BroadcastParticipant, BroadcastState} from "../models/broadcast-message";
import {useParams} from "react-router-dom";
import {Quiz} from "../models/quiz";
import PlayerList from "../components/PlayerList";
import ParticipantDot from "../components/ParticipantDot";

interface CreatorWSTestPageProps {
  sdk: BackendSdk;
}

function CreatorWSTestPage({sdk}: CreatorWSTestPageProps) {
  const {game} = useParams();

  const [client, setClient] = useState(undefined as CreatorGameClient | undefined);
  const [players, setPlayers] = useState([] as BroadcastParticipant[]);
  const [creator, setCreator] = useState({} as BroadcastParticipant);
  const [quiz, setQuiz] = useState({} as Quiz);
  const [currentQuestion, setCurrentQuestion] = useState('00000000-0000-0000-0000-000000000000');

  useEffect(() => {
    const callbacks: GameCallbacks = {
      state(state: BroadcastState) {
        setPlayers(state.players);
        setCreator(state.creator);
        setCurrentQuestion(state.currentQuestion);
      },
      close() {
        alert('disconnected');
      },
      error: console.error,
    }

    const creatorClient = sdk.getCreatorClient(game!, callbacks);
    setClient(creatorClient);
    creatorClient.connect();

    sdk.getQuizByGame(game!).then(setQuiz);

    return () => {
      creatorClient.close();
    }
  }, [game, sdk]);

  const question = quiz?.multipleChoiceQuestions?.find((q) => q.id === currentQuestion);

  return <div>
    <p>Quiz: {quiz.name}</p>
    <PlayerList players={players}/>
    <ParticipantDot participant={creator}/>
    <p>Current question: {question?.title}</p>
    <button onClick={() => client?.next()}>Next Question</button>
  </div>
}

export default CreatorWSTestPage;