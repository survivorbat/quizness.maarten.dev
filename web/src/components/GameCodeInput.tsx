import {Button, InputLabel, TextField} from "@mui/material";
import React, {FormEvent, useState} from "react";

interface GameCodeInputProps {
  callback: (code: string) => Promise<boolean>;
}

function GameCodeInput({callback}: GameCodeInputProps) {
  const [code, setCode] = useState('');
  const [error, setError] = useState(false);

  const handleSubmit = async (event: FormEvent) => {
    setError(false);
    event.preventDefault();
    setError(await callback(code));
  }

  return <form onSubmit={handleSubmit}>
    <InputLabel htmlFor="join-code">Game Code</InputLabel>
    <TextField id="join-code" value={code} error={error} onChange={(e) => setCode(e.target.value.toUpperCase().substr(0, 6))}/>
    <Button type="submit">Join</Button>
  </form>
}


export default GameCodeInput