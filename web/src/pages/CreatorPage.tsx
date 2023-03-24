import BackendSdk from "../logic/sdk";
import {useEffect, useState} from "react";
import Creator from "../models/creator";
import {Quiz} from "../models/quiz";
import QuizComponent from "../components/QuizComponent";
import { Typography } from "@mui/material";
import Grid from "@mui/material/Grid";

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

  return <div className="Creator">
    <h1>Welcome {creator.nickname}</h1>
    <Typography variant="h4"> Your Quizzes</Typography>

      <Grid container spacing={2}>
      {quizzes.map((quiz: Quiz) => 
        <Grid item xs={12} sm={6} md={6}>
          <QuizComponent {...quiz}/>
        </Grid>
      )}
      </Grid>
  </div>
}

export default CreatorPage;
