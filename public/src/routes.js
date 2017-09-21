import Lookup from './components/Lookup'
import BulkLookup from './components/BulkLookup'
import Docs from './components/Docs'
import About from './components/About'
import NotFound from './components/NotFound'

export default [
  {
    name: 'lookup',
    path: '/',
    component: Lookup,
    meta: { title: 'Lookup address' },
  },
  {
    name: 'bulkLookup',
    path: '/lookup/bulk',
    component: BulkLookup,
    meta: { title: 'Bulk lookup' },
  },
  {
    name: 'apidocs',
    path: '/lookup/docs',
    component: Docs,
    meta: { title: 'API docs' },
  },
  {
    name: 'about',
    path: '/about',
    component: About,
    meta: { title: 'About this tool' },
  },
  {
    name: 'catchall',
    path: '*',
    redirect: '/404'
  },
  {
    name: '404',
    path: '/404',
    component: NotFound,
    meta: { title: 'Page not found' },
  },
]
