import axios from 'axios'

const api = axios.create({
  baseURL: '/api',
})

export function createTask(data) {
  return api.post('/tasks', data)
}

export function listTasks(params = {}) {
  return api.get('/tasks', { params })
}

export function toggleTaskStatus(id) {
  return api.patch(`/tasks/${id}/status`)
}

export function getDashboard() {
  return api.get('/dashboard')
}
