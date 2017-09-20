import Lookup from './components/Lookup'
import BulkLookup from './components/BulkLookup'
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
    path: '/lookup/bulk',
    component: BulkLookup,
    name: 'bulkLookup',
    meta: { title: 'Bulk Lookup' },
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
