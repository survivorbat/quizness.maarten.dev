import React from 'react'
import { Typography } from '@mui/material'
import { MultipleChoiceQuestion } from '../models/quiz'

function OverviewQuestionComponent(multipleChoiceQuestion: MultipleChoiceQuestion) {
  return (
    <Typography gutterBottom variant="h4">
      {multipleChoiceQuestion.title}
    </Typography>
  )
}

export default OverviewQuestionComponent
