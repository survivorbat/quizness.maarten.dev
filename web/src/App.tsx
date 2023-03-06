import React, {useState} from 'react';
import {BrowserRouter, Route, Routes} from 'react-router-dom';
import './App.css';
import FrontPage from "./pages/FrontPage";
import CreatorPage from "./pages/CreatorPage";
import Login from "./components/Login";
import AuthPage from "./pages/AuthPage";
import BackendSdk from "./logic/sdk";

function App() {
    const [sdk] = useState(new BackendSdk(process.env.REACT_APP_BACKEND_URL as string));

  return (
    <BrowserRouter>
      <Login/>
      <Routes>
        <Route path="/" element={<FrontPage sdk={sdk}/>}/>
        <Route path="/auth" element={<AuthPage sdk={sdk}/>}/>
        <Route path="/creator" element={<CreatorPage sdk={sdk}/>}/>
      </Routes>
    </BrowserRouter>
  );
}

export default App;
