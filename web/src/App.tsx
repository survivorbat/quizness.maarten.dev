import React, { useState } from 'react'
import { BrowserRouter, Route, Routes } from 'react-router-dom'
import './App.css'
import FrontPage from './pages/FrontPage'
import CreatorPage from './pages/CreatorPage'
import LoginPage from './pages/LoginPage'
import BackendSdk from './logic/sdk'
import { Container, Grid } from '@mui/material'
import ProtectedRoute from './components/ProtectedRoute'
import LogoutPage from './pages/LogoutPage'
import { createTheme, ThemeProvider } from '@mui/material/styles'
import QuizPage from './pages/QuizPage'
import PlayerGame from './pages/PlayerGame'
import CreatorGame from './pages/CreatorGame'
import Player from './models/player'
import Navbar from './components/Navbar'

const darkTheme = createTheme({
  palette: {
    mode: 'dark',
  },
})

function App() {
  const [token, setToken] = useState(localStorage.getItem('token'))

  const sdk = new BackendSdk(token)

  const loginCallback = (token: string) => {
    localStorage.setItem('token', token)
    setToken(token)
  }
  const logoutCallback = () => {
    localStorage.removeItem('token')
    setToken('')
  }

  const joinGame = async (code: string): Promise<Player | null> => {
    try {
      const { id } = await sdk.getGameByCode(code)
      return await sdk.createPlayer(id)
    } catch {
      return null
    }
  }

  const authenticated = token !== null

  return (
    <ThemeProvider theme={darkTheme}>
      <BrowserRouter>
        <Container>
          <Grid container>
            <Navbar authenticated={authenticated} />
            <Routes>
              <Route path='/' element={<FrontPage codeSubmitCallback={joinGame} />} />
              <Route
                path='/login'
                element={
                  <LoginPage
                    successCallback={loginCallback}
                    authenticateFunction={async (token) => await sdk.authenticate(token)}
                  />
                }
              />
              <Route path='/logout' element={<LogoutPage callback={logoutCallback} />} />
              <Route path='/games/:game/players/:player' element={<PlayerGame sdk={sdk} />} />

              <Route
                path='/games/:game'
                element={
                  <ProtectedRoute authenticated={authenticated}>
                    <CreatorGame sdk={sdk} />
                  </ProtectedRoute>
                }
              />
              <Route
                path='/creator'
                element={
                  <ProtectedRoute authenticated={authenticated}>
                    <CreatorPage sdk={sdk} />
                  </ProtectedRoute>
                }
              />
              <Route
                path='/creator/:quiz'
                element={
                  <ProtectedRoute authenticated={authenticated}>
                    <QuizPage sdk={sdk} />
                  </ProtectedRoute>
                }
              />
            </Routes>
          </Grid>
        </Container>
      </BrowserRouter>
    </ThemeProvider>
  )
}

export default App
