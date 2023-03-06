import React from "react";
import {Grid} from "@mui/material";
import LoginButton from "./LoginButton";

function Header() {
  return <Grid xs={12}>
    <LoginButton/>
  </Grid>
}

export default Header;
