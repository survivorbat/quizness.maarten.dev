import React, {useEffect, useState} from "react";
import {Navigate} from "react-router-dom";

interface LogoutPageProps {
  callback: () => void;
}

function LogoutPage({callback}: LogoutPageProps) {
  const [loggedOut, setLoggedOut] = useState(false);

  useEffect(() => {
    callback();
    setLoggedOut(true);
  }, [callback])

  return loggedOut ? <Navigate to="/"/> : <span>Logging out...</span>;
}

export default LogoutPage;