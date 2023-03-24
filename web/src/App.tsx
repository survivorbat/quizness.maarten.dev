import React, {useState} from 'react';
import {BrowserRouter, Route, Routes} from 'react-router-dom';
import './App.css';
import FrontPage from "./pages/FrontPage";
import CreatorPage from "./pages/CreatorPage";
import LoginPage from "./pages/LoginPage";
import BackendSdk from "./logic/sdk";
import {Container, Grid} from "@mui/material";
import Header from "./components/Header";
import ProtectedRoute from "./components/ProtectedRoute";
import LogoutPage from "./pages/LogoutPage";
import { ThemeProvider, createTheme } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
import QuizPage from './pages/QuizPage';

const darkTheme = createTheme({
  palette: {
    mode: 'dark',
  },
});
import PlayerGame from "./pages/PlayerGame";
import CreatorGame from "./pages/CreatorGame";
import Player from "./models/player";

function App() {
  const [token, setToken] = useState(localStorage.getItem('token'));

  const sdk = new BackendSdk(token);

  const loginCallback = (token: string) => {
    localStorage.setItem('token', token);
    setToken(token);
  }
  const logoutCallback = () => {
    localStorage.removeItem('token');
    setToken('');
  }

  const joinGame = async (code: string): Promise<Player | null> => {
    try {
      const {id} = await sdk.getGameByCode(code);
      return await sdk.createPlayer(id);
    } catch {
      return null
    }
  }

  return (
    <ThemeProvider theme={darkTheme}>
    <BrowserRouter>
      <Container>
        <Grid container>
          <Header authenticated={!!token}/>
          <Routes>
            <Route path="/" element={<FrontPage codeSubmitCallback={joinGame}/>}/>
            <Route path="/login" element={<LoginPage successCallback={loginCallback}
                                                     authenticateFunction={(token) => sdk.authenticate(token)}/>}/>
            <Route path="/logout" element={<LogoutPage callback={logoutCallback}/>}/>
            <Route path="/games/:game/players/:player" element={<PlayerGame sdk={sdk}/>}/>

            <Route path="/creator"
                   element={<ProtectedRoute authenticated={!!token}><CreatorPage sdk={sdk}/></ProtectedRoute>}/>

            <Route path="/games/:game"
                   element={<ProtectedRoute authenticated={!!token}><CreatorGame sdk={sdk}/></ProtectedRoute>}/>
            <Route path="/creator/quiz" element={<QuizPage sdk = {sdk}/>}></Route>
        </Routes>
        </Grid>
      </Container>
    </BrowserRouter>
    </ThemeProvider>
  );
}

export default App;
