import React, { useEffect, useState } from 'react'
import { useParams } from 'react-router-dom'
import BackendSdk from '../logic/sdk'
import { MultipleChoiceQuestion, Quiz } from '../models/quiz'
import { Grid, Typography } from '@mui/material'
import '../styles/QuizPage.css'

interface QuizPageProps {
  sdk: BackendSdk
}

function QuizPage({ sdk }: QuizPageProps) {
  const { quiz: id } = useParams()

  const [quiz, setQuiz] = useState<Quiz | undefined>(undefined)

  useEffect(() => {
    if (id !== undefined) {
      sdk.getQuizById(id).then(setQuiz).catch(console.error)
    }
  }, [sdk, id])

  if (quiz == null) {
    return <span>Loading...</span>
  }

  return (
    <Grid item xs={12}>
      <div className='quiz'>
        <Typography gutterBottom variant='h4'>
          {quiz.name}
        </Typography>
        <Typography gutterBottom variant='h6'>
          {quiz.description}
        </Typography>
        {quiz.multipleChoiceQuestions.map((multipleChoiceQuestion: MultipleChoiceQuestion) => {
          return (
            <Typography key={multipleChoiceQuestion.id} gutterBottom variant='h6'>
              {multipleChoiceQuestion.title}
            </Typography>
          )
        })}
      </div>
    </Grid>
  )
}

export default QuizPage
