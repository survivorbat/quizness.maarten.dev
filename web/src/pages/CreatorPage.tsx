import BackendSdk from '../logic/sdk'
import { useEffect, useState } from 'react'
import type Creator from '../models/creator'
import { Quiz } from '../models/quiz'
import QuizComponent from '../components/QuizComponent'
import { Typography } from '@mui/material'
import Grid from '@mui/material/Grid'
import '../styles/CreatorPage.css'

interface CreatorPageProps {
  sdk: BackendSdk
}

function CreatorPage({ sdk }: CreatorPageProps) {
  const [creator, setCreator] = useState(undefined as Creator | undefined)
  const [quizzes, setQuizzes] = useState(undefined as Quiz[] | undefined)

  const refresh = () => {
    sdk.getCreator().then(setCreator).catch(console.error)
    sdk.getQuizzes().then(setQuizzes).catch(console.error)
  }

  useEffect(refresh, [sdk])

  if (creator == null || quizzes == null) {
    return <span>Loading...</span>
  }

  return (
    <div className='creator'>
      <h1>Welcome {creator.nickname}</h1>
      <Typography variant='h4'> Your Quizzes</Typography>

      <Grid container spacing={2}>
        {quizzes.map((quiz: Quiz) => (
          <Grid key={quiz.id} item xs={12} sm={6} md={6}>
            <QuizComponent {...quiz} />
          </Grid>
        ))}
      </Grid>
    </div>
  )
}

export default CreatorPage
