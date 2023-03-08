import React, {useState} from 'react';
import {BrowserRouter, Route, Routes} from 'react-router-dom';
import './App.css';
import FrontPage from "./pages/FrontPage";
import CreatorPage from "./pages/CreatorPage";
import AuthPage from "./pages/AuthPage";
import BackendSdk from "./logic/sdk";
import {Grid} from "@mui/material";
import Header from "./components/Header";
import Creator from "./models/creator";

const backendUrl = process.env.REACT_APP_BACKEND_URL as string;

function App() {
  const [token, setToken] = useState(localStorage.getItem('token'));

  const sdk = new BackendSdk(backendUrl, token);

  const authCallback = (token: string) => {
    localStorage.setItem('token', token);
    setToken(token);
  }

  return (
    <BrowserRouter>
      <Grid container>
        <Header/>
        <Routes>
          <Route path="/" element={<FrontPage/>}/>
          <Route path="/auth" element={<AuthPage successCallback={authCallback} authenticateFunction={(token) => sdk.authenticate(token)}/>}/>
          <Route path="/creator" element={<CreatorPage creator={{} as Creator}/>}/>
        </Routes>
      </Grid>
    </BrowserRouter>
  );
}

export default App;
