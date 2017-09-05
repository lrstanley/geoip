import Hello from './components/Hello'
import Counter from './components/Counter'

export default [
  {
    path: '/',
    component: Hello,
    name: 'hello'
  }, {
    path: '/counter',
    component: Counter,
    name: 'counter'
  }
]
