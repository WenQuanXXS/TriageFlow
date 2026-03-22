import axios from 'axios'

const api = axios.create({
  baseURL: '/api',
})

export function createTask(data) {
  return api.post('/tasks', data)
}

export function getTaskDetail(id) {
  return api.get(`/tasks/${id}`)
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

// Queue APIs

export function listQueue(params = {}) {
  return api.get('/queue', { params })
}

export function getQueuePosition(taskId) {
  return api.get(`/queue/${taskId}/position`)
}

export function callQueuePatient(taskId) {
  return api.patch(`/queue/${taskId}/call`)
}

export function completeQueuePatient(taskId) {
  return api.patch(`/queue/${taskId}/complete`)
}
