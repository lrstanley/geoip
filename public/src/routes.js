import Lookup from './components/Lookup'
import Docs from './components/Docs'
import About from './components/About'

export default [
  { path: '/', component: Lookup, name: 'lookup' },
  { path: '/docs/api', component: Docs, name: 'apidocs' },
  { path: '/about', component: About, name: 'about' }
]
