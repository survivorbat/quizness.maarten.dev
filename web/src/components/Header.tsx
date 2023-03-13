import React from "react";
import {Grid} from "@mui/material";
import LoginButton from "./LoginButton";
import LogoutButton from "./LogoutButton";
import Navbar from "./Navbar";

interface HeaderProps {
  authenticated: boolean;
}

function Header({authenticated}: HeaderProps) {
  return<>
    <Navbar authenticated={authenticated}/>
    <Grid>
    </Grid>
  </>
}

export default Header;
