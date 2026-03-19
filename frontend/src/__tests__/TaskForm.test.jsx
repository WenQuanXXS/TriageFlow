import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { BrowserRouter } from 'react-router-dom'
import { LocaleProvider } from '../locales'
import TaskForm from '../components/TaskForm'

vi.mock('../api/tasks', () => ({
  createTask: vi.fn(() => Promise.resolve({ data: {} })),
}))

import { createTask } from '../api/tasks'

function renderForm() {
  return render(
    <BrowserRouter>
      <LocaleProvider>
        <TaskForm />
      </LocaleProvider>
    </BrowserRouter>
  )
}

describe('TaskForm', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('submits form with valid data', async () => {
    const user = userEvent.setup()
    renderForm()

    // Default locale is zh, so use Chinese placeholders
    await user.type(screen.getByPlaceholderText('请输入患者姓名'), 'John Doe')
    await user.type(screen.getByPlaceholderText('请描述主诉症状'), 'Severe headache')

    await user.click(screen.getByRole('button', { name: /提\s*交/i }))

    await waitFor(() => {
      expect(createTask).toHaveBeenCalledWith({
        patient_name: 'John Doe',
        chief_complaint: 'Severe headache',
        priority: 'normal',
      })
    })
  })
})
