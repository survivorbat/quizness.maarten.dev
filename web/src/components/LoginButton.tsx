import Avatar from '@mui/material/Avatar'
import IconButton from '@mui/material/IconButton'

export const getGoogleUrl = () => {
  const rootUrl = 'https://accounts.google.com/o/oauth2/v2/auth'

  const options = {
    redirect_uri: process.env.REACT_APP_AUTH_REDIRECT as string,
    client_id: process.env.REACT_APP_AUTH_CLIENT_ID as string,
    access_type: 'offline',
    response_type: 'code',
    prompt: 'consent',
    scope: 'openid',
    state: 'a'
  }

  const qs = new URLSearchParams(options)

  return `${rootUrl}?${qs.toString()}`
}

function LoginButton() {
  return (
    <IconButton href={getGoogleUrl()} sx={{ p: 0 }}>
      <Avatar alt="Remy Sharp" src="/static/images/avatar/2.jpg" />
    </IconButton>
  )
}

export default LoginButton
