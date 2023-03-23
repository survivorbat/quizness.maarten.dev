import BackendSdk from "../logic/sdk";
import {useEffect, useState} from "react";
import Creator from "../models/creator";
import {Quiz} from "../models/quiz";
import {Link} from "react-router-dom";
import {Button} from "@mui/material";
import Game from "../models/game";

interface CreatorPageProps {
  sdk: BackendSdk;
}

const getGameButton = (game: Game, startGame: (id: string) => void) => {
  if (game.code) {
    return <Button key={game.id}>
      <Link to={`/games/${game.id}`}>{game.code}</Link>
    </Button>
  }

  return <Button onClick={() => startGame(game.id)}>Start Game</Button>
}

function CreatorPage({sdk}: CreatorPageProps) {
  const [creator, setCreator] = useState(undefined as Creator | undefined);
  const [quizzes, setQuizzes] = useState(undefined as Quiz[] | undefined)

  const refresh = () => {
    sdk.getCreator().then(setCreator);
    sdk.getQuizzes().then(setQuizzes);
  }

  useEffect(refresh, [sdk]);

  const createGame = (quiz: string) => {
    const playerLimit = prompt('Player limit?')!
    sdk.createGame(quiz, {playerLimit: parseInt(playerLimit, 10)}).then(refresh)
  }

  const startGame = (game: string) => {
    sdk.startGame(game).then(refresh);
  }

  if (!creator || !quizzes) {
    return <span>Loading...</span>
  }

  return <div>
    <h1>Welcome {creator.nickname}</h1>
    <p>Your quizzes:</p>
    <table>
      <thead>
        <tr>
          <th>Name</th>
          <th>Description</th>
          <th>Questions</th>
          <th>Games</th>
        </tr>
      </thead>
      <tbody>
      {quizzes.map((quiz: Quiz) =>
        <tr key={quiz.id}>
          <td>{quiz.name}</td>
          <td>{quiz.description}</td>
          <td>{quiz.multipleChoiceQuestions.length}</td>
          <td>
            <ul>
              {quiz.games.map((game) =>
                <li key={game.id}>{getGameButton(game, startGame)}</li>
              )}
              <li><Button onClick={() => createGame(quiz.id)}>Create</Button></li>
            </ul>
          </td>
        </tr>
      )}
      </tbody>
    </table>
  </div>
}

export default CreatorPage;
