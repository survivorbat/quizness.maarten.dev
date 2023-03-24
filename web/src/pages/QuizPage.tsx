import React, {useEffect, useState} from "react";
import {Routes, Route, useParams} from 'react-router-dom';
import BackendSdk from "../logic/sdk";
import { MultipleChoiceQuestion, QuestionOption, Quiz } from "../models/quiz";
import { useLocation } from "react-router-dom";
import { Grid, Typography } from "@mui/material";
import OverviewQuestionComponent from "../components/OverviewQuestionComponent";

interface QuizPageProps{
    sdk : BackendSdk
}

function QuizPage({sdk} : QuizPageProps) {
    const location = useLocation();
    const {quiz} = location.state;
    const questions = quiz.multipleChoiceQuestions
    console.log(questions)

    return <Grid item xs={12}>
            <div className="Quiz">

            <Typography gutterBottom variant="h4">{quiz.name}</Typography>
            <Typography gutterBottom variant="h6">{quiz.description}</Typography>
            {quiz.multipleChoiceQuestions.map((multipleChoiceQuestion : MultipleChoiceQuestion) => {
                <Typography gutterBottom variant="h6">{multipleChoiceQuestion.title}</Typography>
            })}
            

            {/* {questions.map((multipleChoiceQuestions : MultipleChoiceQuestion) => {
                <Grid item xs={12}>
                    <OverviewQuestionComponent {...multipleChoiceQuestions}/>
                </Grid>
            })} */}
            </div>
            </Grid>

}

export default QuizPage