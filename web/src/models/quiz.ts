import Game from './game'

export interface Quiz {
  id: string
  name: string
  description: string
  multipleChoiceQuestions: MultipleChoiceQuestion[]
  games: Game[]
}

export interface MultipleChoiceQuestion {
  id: string
  title: string
  description: string
  durationInSeconds: number
  category: string
  order: number
  answerID: string
  options: QuestionOption[]
}

export interface QuestionOption {
  id: string
  textOption: string
}
