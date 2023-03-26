import React, { useEffect, useState } from 'react'
import { useLocation, useNavigate } from 'react-router-dom'

function useQuery() {
  const { search } = useLocation()

  return React.useMemo(() => new URLSearchParams(search), [search])
}

interface LoginPageProps {
  authenticateFunction: (code: string) => Promise<string>
  successCallback: (token: string) => void
}

function LoginPage({ authenticateFunction, successCallback }: LoginPageProps) {
  const [invalid, setInvalid] = useState(false)

  const query = useQuery()

  const navigate = useNavigate()
  const code = query.get('code')

  useEffect(() => {
    // Redirect to the home page if there is no code
    if (code == null) {
      setInvalid(true)
      return
    }

    authenticateFunction(code)
      .then(successCallback)
      .then(() => {
        navigate('/creator')
      })
      .catch(() => {
        setInvalid(true)
      })
  }, [successCallback, authenticateFunction, navigate, code])

  if (invalid) {
    return <span>Failed to authenticate you, please try again</span>
  }

  return <span>Authenticating...</span>
}

export default LoginPage
