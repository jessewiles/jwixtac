import PublicIndex from './views/public/index.svelte'
import PublicLayout from './views/public/layout.svelte'

const routes = [
  {
    name: '/',
    component: PublicIndex,
    layout: PublicLayout,
  },
  {
      name: 'story',
      layout: PublicLayout  // TBD
  },
]

export { routes }
