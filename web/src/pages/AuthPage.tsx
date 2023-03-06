import React, {useEffect} from "react";
import {useLocation, useNavigate} from "react-router-dom";
import BackendSdk from "../logic/sdk";

function useQuery() {
  const {search} = useLocation();

  return React.useMemo(() => new URLSearchParams(search), [search]);
}

interface AuthPageProps {
  sdk: BackendSdk;
}

function AuthPage({sdk}: AuthPageProps) {
  const query = useQuery();

  const navigate = useNavigate();
  const code = query.get('code');

  useEffect(() => {
    if (!code) {
      navigate("/")
      return
    }

    sdk.authenticate(code)
        .then(() => navigate('/'))
        .catch(console.error)
  })

  return <span>Authenticating...</span>
}

export default AuthPage;