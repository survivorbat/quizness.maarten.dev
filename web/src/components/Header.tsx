import React from "react";
import {Button, Grid} from "@mui/material";
import LoginButton from "./LoginButton";
import {Link} from "react-router-dom";
import Navbar from "./Navbar";

interface HeaderProps {
  authenticated: boolean;
}

function Header({authenticated}: HeaderProps) {
  return<>
    <Navbar authenticated={authenticated}/>
  </>
}

export default Header;
