import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { BrowserRouter } from 'react-router-dom'
import { LocaleProvider } from '../locales'
import TaskForm from '../components/TaskForm'

vi.mock('../api/tasks', () => ({
  createTask: vi.fn(() => Promise.resolve({ data: { id: 1 } })),
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

  it('submits form with patient name and chief complaint only', async () => {
    const user = userEvent.setup()
    renderForm()

    await user.type(screen.getByPlaceholderText('请输入患者姓名'), '张三')
    await user.type(screen.getByPlaceholderText('请详细描述主诉症状'), '头痛三天')

    await user.click(screen.getByRole('button', { name: /提\s*交/i }))

    await waitFor(() => {
      expect(createTask).toHaveBeenCalledWith(
        expect.objectContaining({
          patient_name: '张三',
          chief_complaint: '头痛三天',
        })
      )
    })
  })

  it('renders all patient info fields', () => {
    renderForm()

    // Ant Design renders duplicate DOM elements in test environment, use getAllBy
    expect(screen.getAllByPlaceholderText('请输入年龄').length).toBeGreaterThan(0)
    expect(screen.getAllByText('性别').length).toBeGreaterThan(0)
    expect(screen.getAllByPlaceholderText('请输入体温（°C）').length).toBeGreaterThan(0)
    expect(screen.getAllByText('疼痛等级').length).toBeGreaterThan(0)
    expect(screen.getAllByPlaceholderText('如怀孕、过敏史等').length).toBeGreaterThan(0)
  })
})
