import React from "react";
import { Button, Card, CardActions, CardContent, CardMedia, Typography } from "@mui/material";
import { MultipleChoiceQuestion, Quiz } from "../models/quiz";
import { Link } from "react-router-dom";

function OverviewQuestionComponent(multipleChoiceQuestion : MultipleChoiceQuestion) {

    return <Typography gutterBottom variant="h4">{multipleChoiceQuestion.title}</Typography>

}

export default OverviewQuestionComponent