import React from 'react';
import {BrowserRouter, Route, Routes} from 'react-router-dom';
import './App.css';
import FrontPage from "./pages/FrontPage";
import CreatorPage from "./pages/CreatorPage";
import Login from "./components/Login";
import AuthPage from "./pages/AuthPage";

function App() {
  return (
    <BrowserRouter>
      <Login/>
      <Routes>
        <Route path="/" element={<FrontPage/>}/>
        <Route path="/auth" element={<AuthPage/>}/>
        <Route path="/creator" element={<CreatorPage/>}/>
      </Routes>
    </BrowserRouter>
  );
}

export default App;
