import React from 'react'
import { RouterComponent } from '@/core'
import { ProfileProvider } from '@/core/profile'

export const App = () => {
  return (
    <ProfileProvider>
      <RouterComponent />
    </ProfileProvider>
  )
}
