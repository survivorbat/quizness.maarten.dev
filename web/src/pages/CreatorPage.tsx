import BackendSdk from "../logic/sdk";
import {useEffect, useState} from "react";
import Creator from "../models/creator";
import {Quiz} from "../models/quiz";
import {Link} from "react-router-dom";

interface CreatorPageProps {
  sdk: BackendSdk;
}

function CreatorPage({sdk}: CreatorPageProps) {
  const [creator, setCreator] = useState(undefined as Creator | undefined);
  const [quizzes, setQuizzes] = useState(undefined as Quiz[] | undefined)

  useEffect(() => {
    sdk.getCreator().then(setCreator);
    sdk.getQuizzes().then(setQuizzes);
  }, [sdk]);

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
          <th>Games</th>
        </tr>
      </thead>
      <tbody>
      {quizzes.map((quiz: Quiz) =>
        <tr key={quiz.id}>
          <td>{quiz.name}</td>
          <td>{quiz.description}</td>
          <td>{quiz.multipleChoiceQuestions.length}</td>
          <td>
            {quiz.games?.map((game) => <Link key={game.id} to={`/games/${game.id}`}>{game.code}</Link>)}
          </td>
        </tr>
      )}
      </tbody>
    </table>
  </div>
}

export default CreatorPage;
