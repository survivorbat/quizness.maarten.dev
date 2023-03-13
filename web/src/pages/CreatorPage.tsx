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

function CreatorPage({sdk}: CreatorPageProps) {
  const [creator, setCreator] = useState(undefined as Creator | undefined);
  const [quizzes, setQuizzes] = useState(undefined as Quiz[] | undefined)

  useEffect(() => {
    if (!creator || !quizzes) {
      sdk.getCreator().then(setCreator);
      sdk.getQuizzes().then(setQuizzes);
    }
  });

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
