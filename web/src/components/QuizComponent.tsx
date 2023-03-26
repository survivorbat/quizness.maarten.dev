import React from 'react'
import { Button, Card, CardActions, CardContent, Typography } from '@mui/material'
import { Quiz } from '../models/quiz'
import { Link } from 'react-router-dom'

function QuizComponent(quiz: Quiz) {
  return (
    <Card>
      <CardContent>
        <Typography gutterBottom variant="h5" component="div">
          {quiz.name}
        </Typography>
        <Typography variant="body2" color="text.secondary">
          {quiz.description}
        </Typography>
      </CardContent>
      <CardActions>
        <Button size="small">
          {' '}
          <Link to={'/creator/quiz'} state={{ quiz }}>
            {' '}
            play{' '}
          </Link>
        </Button>
        <Button size="small">View</Button>
      </CardActions>
    </Card>
  )
}

export default QuizComponent
