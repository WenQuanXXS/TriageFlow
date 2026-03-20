import { useState, useEffect, useCallback } from 'react'
import { Table, Tag, Button, Select, Space, message } from 'antd'
import { Link } from 'react-router-dom'
import { listTasks, toggleTaskStatus } from '../api/tasks'
import { useLocale } from '../locales'

const statusColors = {
  pending: 'orange',
  in_progress: 'blue',
  completed: 'green',
}

const priorityColors = {
  urgent: 'red',
  high: 'volcano',
  normal: 'blue',
  low: 'default',
}

const statusLabelKeys = {
  pending: 'statusPending',
  in_progress: 'statusInProgress',
  completed: 'statusCompleted',
}

const priorityLabelKeys = {
  urgent: 'priorityUrgent',
  high: 'priorityHigh',
  normal: 'priorityNormal',
  low: 'priorityLow',
}

export default function TaskList() {
  const [tasks, setTasks] = useState([])
  const [loading, setLoading] = useState(false)
  const [statusFilter, setStatusFilter] = useState('')
  const [priorityFilter, setPriorityFilter] = useState('')
  const { t, tDept } = useLocale()

  const statusOptions = [
    { value: '', label: t('allStatus') },
    { value: 'pending', label: t('pending') },
    { value: 'in_progress', label: t('inProgress') },
    { value: 'completed', label: t('completed') },
  ]

  const priorityOptions = [
    { value: '', label: t('allPriority') },
    { value: 'urgent', label: t('urgent') },
    { value: 'high', label: t('high') },
    { value: 'normal', label: t('normal') },
    { value: 'low', label: t('low') },
  ]

  const fetchTasks = useCallback(async () => {
    setLoading(true)
    try {
      const params = {}
      if (statusFilter) params.status = statusFilter
      if (priorityFilter) params.priority = priorityFilter
      const { data } = await listTasks(params)
      setTasks(data)
    } catch {
      message.error(t('loadTasksFail'))
    } finally {
      setLoading(false)
    }
  }, [statusFilter, priorityFilter, t])

  useEffect(() => {
    fetchTasks()
  }, [fetchTasks])

  const handleToggle = async (id) => {
    try {
      await toggleTaskStatus(id)
      fetchTasks()
    } catch {
      message.error(t('updateStatusFail'))
    }
  }

  const columns = [
    {
      title: t('patient'),
      dataIndex: 'patient_name',
      key: 'patient_name',
      render: (name, record) => <Link to={`/tasks/${record.id}`}>{name}</Link>,
    },
    { title: t('chiefComplaint'), dataIndex: 'chief_complaint', key: 'chief_complaint', ellipsis: true },
    {
      title: t('status'),
      dataIndex: 'status',
      key: 'status',
      render: (s) => <Tag color={statusColors[s]}>{t(statusLabelKeys[s])}</Tag>,
    },
    {
      title: t('finalPriority'),
      key: 'final_priority',
      render: (_, record) => {
        const pri = record.final_priority || record.priority
        return (
          <Space size={4}>
            <Tag color={priorityColors[pri]}>{t(priorityLabelKeys[pri])}</Tag>
            {record.rule_triggered && <Tag color="warning">R</Tag>}
          </Space>
        )
      },
    },
    { title: t('finalDepartment'), dataIndex: 'final_department', key: 'final_department', render: (dept) => tDept(dept) },
    {
      title: t('created'),
      dataIndex: 'created_at',
      key: 'created_at',
      render: (v) => new Date(v).toLocaleString(),
    },
    {
      title: t('action'),
      key: 'action',
      render: (_, record) => (
        <Space>
          <Button size="small" onClick={() => handleToggle(record.id)} disabled={record.status === 'completed'}>
            {t('toggleStatus')}
          </Button>
          <Link to={`/tasks/${record.id}`}>
            <Button size="small" type="link">{t('viewDetail')}</Button>
          </Link>
        </Space>
      ),
    },
  ]

  return (
    <div style={{ marginTop: 24 }}>
      <Space style={{ marginBottom: 16 }}>
        <Select
          value={statusFilter}
          onChange={setStatusFilter}
          options={statusOptions}
          style={{ width: 150 }}
        />
        <Select
          value={priorityFilter}
          onChange={setPriorityFilter}
          options={priorityOptions}
          style={{ width: 150 }}
        />
      </Space>
      <Table
        dataSource={tasks}
        columns={columns}
        rowKey="id"
        loading={loading}
        pagination={{ pageSize: 10 }}
      />
    </div>
  )
}
