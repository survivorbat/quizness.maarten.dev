import React, {useState} from 'react';
import {BrowserRouter, Route, Routes} from 'react-router-dom';
import './App.css';
import FrontPage from "./pages/FrontPage";
import CreatorPage from "./pages/CreatorPage";
import LoginPage from "./pages/LoginPage";
import BackendSdk from "./logic/sdk";
import {Grid} from "@mui/material";
import Header from "./components/Header";
import ProtectedRoute from "./components/ProtectedRoute";
import LogoutPage from "./pages/LogoutPage";
import { ThemeProvider, createTheme } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';

const darkTheme = createTheme({
  palette: {
    mode: 'dark',
  },
});

const backendUrl = process.env.REACT_APP_BACKEND_URL as string;

function App() {
  const [token, setToken] = useState(localStorage.getItem('token'));

  const sdk = new BackendSdk(backendUrl, token);

  const loginCallback = (token: string) => {
    localStorage.setItem('token', token);
    setToken(token);
  }
  const logoutCallback = () => {
    localStorage.removeItem('token');
    setToken('');
  }

  return (
    <ThemeProvider theme={darkTheme}>
    <BrowserRouter>
      <Grid container>
        <Header authenticated={!!token}/>
        <Routes>
          <Route path="/" element={<FrontPage/>}/>
          <Route path="/login" element={<LoginPage successCallback={loginCallback}
                                                   authenticateFunction={(token) => sdk.authenticate(token)}/>}/>
          <Route path="/logout" element={<LogoutPage callback={logoutCallback}/>}/>
          <Route path="/creator"
                 element={<ProtectedRoute authenticated={!!token}><CreatorPage sdk={sdk}/></ProtectedRoute>}/>
        </Routes>
      </Grid>
    </BrowserRouter>
    </ThemeProvider>
  );
}

export default App;
