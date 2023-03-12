import BackendSdk from "../logic/sdk";
import {useEffect, useState} from "react";
import Creator from "../models/creator";
import {Quiz} from "../models/quiz";

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

  return <div>
    <h1>Welcome {creator.nickname}</h1>
    <p>Your quizzes:</p>
    <table>
      <thead>
      <tr>
        <th>Name</th>
        <th>Description</th>
        <th>Questions</th>
      </tr>
      </thead>
      <tbody>
      {quizzes.map((quiz: Quiz) => <tr>
        <td>{quiz.id}</td>
        <td>{quiz.description}</td>
        <td>{quiz.multipleChoiceQuestions.length}</td>
      </tr>)}
      </tbody>
    </table>
  </div>
}

export default CreatorPage;
