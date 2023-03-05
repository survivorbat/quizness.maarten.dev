import React, {useEffect} from "react";
import {useLocation, useNavigate} from "react-router-dom";

function useQuery() {
  const { search } = useLocation();

  return React.useMemo(() => new URLSearchParams(search), [search]);
}

function AuthPage() {
  const query = useQuery();

  const navigate = useNavigate();
  const code = query.get('code');

  useEffect(() => {
    if (!code) {
      navigate("/")
      return
    }

    fetch(`http://localhost:8000/api/v1/tokens`, {method: 'post', body: JSON.stringify({code})}).then((result) => {
      const token = result.headers.get('token')

      if (token) {
        localStorage.setItem('token', token)
        navigate("/")
      }
    }).catch(console.error);
  })

  return <span>Authenticating...</span>
}

export default AuthPage;