import React from 'react'
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom'
import { switchRoutes } from './routes'
import { LoginPage, SignupPage, FotosPage, SubirPage, AlbumesPage, FavoritosPage, FotoPage } from '@/scenes'

export const RouterComponent: React.FC = () => {
  return (
    <Router>
      <Routes>
        <Route path={switchRoutes.login} element={<LoginPage />} />
        <Route path={switchRoutes.signup} element={<SignupPage />} />
        <Route path={switchRoutes.fotos} element={<FotosPage />} />
        <Route path={switchRoutes.subir} element={<SubirPage />} />
        <Route path={switchRoutes.albumes} element={<AlbumesPage />} />
        <Route path={switchRoutes.favoritos} element={<FavoritosPage />} />
        <Route path={switchRoutes.foto} element={<FotoPage />} />
      </Routes>
    </Router>
  )
}
}