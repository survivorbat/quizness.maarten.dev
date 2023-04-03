import IconButton from '@mui/material/IconButton'
import Avatar from '@mui/material/Avatar'
import Menu from '@mui/material/Menu'
import MenuItem from '@mui/material/MenuItem'
import Button from '@mui/material/Button'
import React, { useState } from 'react'
import { Link } from 'react-router-dom'

const userSettings = [
  {
    name: 'Profile',
    path: 'profile'
  },
  {
    name: 'Logout',
    path: 'logout'
  }
]
function UserButton() {
  const [anchorElUser, setAnchorElUser] = useState<null | HTMLElement>(null)

  const handleOpenUserMenu = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorElUser(event.currentTarget)
  }

  const handleCloseUserMenu = () => {
    setAnchorElUser(null)
  }
  return <>
    <IconButton onClick={handleOpenUserMenu} sx={{ p: 0 }}>
                  <Avatar alt="Remy Sharp" src="/static/images/avatar/2.jpg" />
                </IconButton>
                <Menu
                  sx= {{ my: '45px' }}
                  id="user-menu"
                  anchorEl={anchorElUser}
                  anchorOrigin={{
                    vertical: 'top',
                    horizontal: 'right'
                  }}
                  keepMounted
                  transformOrigin={{
                    vertical: 'top',
                    horizontal: 'right'
                  }}
                  open={Boolean(anchorElUser)}
                  onClose={handleCloseUserMenu}
                >
                  {userSettings.map((item) => (
                    <MenuItem key={item.name} onClick={handleCloseUserMenu}>
                      <Button
                        component={Link}
                        to={`/${item.path}`}>
                          {item.name}
                      </Button>
                    </MenuItem>
                  ))

                  }
                </Menu>
                </>
}

export default UserButton
