import React from 'react'
import { createEmptyLogin, Login } from './login.vm'

interface Props {
  onLogin: (login: Login) => void
}

export const LoginComponent: React.FC<Props> = (props) => {
  const { onLogin } = props
  const [login, setLogin] = React.useState<Login>(createEmptyLogin())

  const handleNavigation = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    onLogin(login)
  }

  const updateFieldValue = (name: keyof Login) => (e: React.ChangeEvent<HTMLInputElement>) => {
    setLogin({
      ...login,
      [name]: e.target.value,
    })
  }

  return (
    <form onSubmit={handleNavigation}>
      <div className="login-container">
        <input
          placeholder="Username"
          value={login.username}
          onChange={updateFieldValue('username')}
        />
        <input
          placeholder="Password"
          type="password"
          value={login.password}
          onChange={updateFieldValue('password')}
        />
        <button type="submit">login</button>
      </div>
    </form>
  )
}
