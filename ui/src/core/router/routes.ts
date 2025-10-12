import { generatePath } from 'react-router-dom';

interface SwitchRoutes {
  login: string
  signup: string
  fotos: string
  subir: string
  albumes: string
  favoritos: string
  foto: string
}

export const switchRoutes: SwitchRoutes = {
  login: "/",
  signup: "/signup",
  fotos: "/fotos",
  subir: "/subir",
  albumes: "/albumes",
  favoritos: "/favoritos",
  foto: "/foto/:id",
}

interface Routes extends Omit<SwitchRoutes, 'foto'> {
  foto: (id: string) => string
}

export const routes: Routes = {
  ...switchRoutes,
  foto: (id: string) => generatePath(switchRoutes.foto, { id })
}

