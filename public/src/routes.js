import Lookup from './components/Lookup'
import Docs from './components/Docs'
import About from './components/About'

export default [
  {
    path: '/',
    component: Lookup,
    name: 'lookup',
    meta: { title: 'Lookup Address' },
  },
  {
    path: '/docs/api',
    component: Docs,
    name: 'apidocs',
    meta: { title: 'API Docs' },
  },
  {
    path: '/about',
    component: About,
    name: 'about',
    meta: { title: 'About This Tool' },
  }
]
