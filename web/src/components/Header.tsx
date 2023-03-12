import React from "react";
import {Grid} from "@mui/material";
import LoginButton from "./LoginButton";
import LogoutButton from "./LogoutButton";

interface HeaderProps {
  authenticated: boolean;
}

function Header({authenticated}: HeaderProps) {
  return <Grid>
    {authenticated ? <LogoutButton/> : <LoginButton/>}
  </Grid>
}

export default Header;
