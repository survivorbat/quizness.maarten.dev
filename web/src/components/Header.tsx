import React from "react";
import {Button, Grid} from "@mui/material";
import LoginButton from "./LoginButton";
import {Link} from "react-router-dom";

interface HeaderProps {
  authenticated: boolean;
}

function Header({authenticated}: HeaderProps) {
  let menu = [<LoginButton/>]

  if (authenticated) {
    menu = [
      <Button key="quizzes">
        <Link to="/creator">Your Quizzes</Link>
      </Button>,
      <Button key="logout">
        <Link to="/logout">Logout</Link>
      </Button>
    ]
  }

  return <Grid container>
    <Grid>
      <Button><Link to="/">Home</Link></Button>
      <Button><Link to="/about">About</Link></Button>
    </Grid>
    <Grid>
      {menu}
    </Grid>
  </Grid>
}

export default Header;
