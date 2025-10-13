import React from 'react'
import { useNavigate } from 'react-router-dom'
import { routes } from '@/core'
import { ProfileContext } from '@/core/profile'
import { LoginComponent } from './login.component'
import { doLogin } from './login.api'
import { Login } from './login.vm'

const useLoginHook = () => {
  const navigate = useNavigate()
  const { setUserProfile } = React.useContext(ProfileContext)

  const loginSuccededAction = (username) => {
    setUserProfile({ userName: username })
    navigate(routes.list)
  }

  const loginFailedAction = () => {
    alert('Credentials not correct')
  }

  const handleLogin = (login: Login) => {
      const {username, password} = login
      doLogin(username, password).then((result) => {
        if (result) loginSuccededAction(result)
        else loginFailedAction()
      })
    }

  return { handleLogin }
}

export const LoginContainer: React.FC = () => {

  const { handleLogin } = useLoginHook()

  return <LoginComponent onLogin={handleLogin} />
}
