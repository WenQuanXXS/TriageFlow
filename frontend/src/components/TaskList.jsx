import { useState, useEffect, useCallback } from 'react'
import { Table, Button, Select, Space, message } from 'antd'
import { Link } from 'react-router-dom'
import { EyeOutlined } from '@ant-design/icons'
import { listTasks, toggleTaskStatus } from '../api/tasks'
import { useLocale } from '../locales'
import StatusBadge, { statusLabelKeys } from './shared/StatusBadge'
import PriorityBadge, { priorityLabelKeys } from './shared/PriorityBadge'

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
      render: (name, record) => (
        <Link to={`/console/tasks/${record.id}`} style={{ fontWeight: 500 }}>
          {name}
        </Link>
      ),
    },
    {
      title: t('chiefComplaint'),
      dataIndex: 'chief_complaint',
      key: 'chief_complaint',
      ellipsis: true,
      render: (text) => <span style={{ color: '#5a6a7e' }}>{text}</span>,
    },
    {
      title: t('status'),
      dataIndex: 'status',
      key: 'status',
      width: 130,
      render: (s) => <StatusBadge status={s} label={t(statusLabelKeys[s])} />,
    },
    {
      title: t('finalPriority'),
      key: 'final_priority',
      width: 140,
      render: (_, record) => {
        const pri = record.final_priority || record.priority
        return (
          <Space size={4}>
            <PriorityBadge priority={pri} label={t(priorityLabelKeys[pri])} />
            {record.rule_triggered && <span className="tf-rule-badge">R</span>}
          </Space>
        )
      },
    },
    {
      title: t('finalDepartment'),
      dataIndex: 'final_department',
      key: 'final_department',
      render: (dept) => <span style={{ fontWeight: 500 }}>{tDept(dept)}</span>,
    },
    {
      title: t('created'),
      dataIndex: 'created_at',
      key: 'created_at',
      width: 170,
      render: (v) => (
        <span style={{ color: '#8e9aab', fontSize: 13 }}>
          {new Date(v).toLocaleString()}
        </span>
      ),
    },
    {
      title: t('action'),
      key: 'action',
      width: 180,
      render: (_, record) => (
        <Space size={8}>
          <Button
            size="small"
            onClick={() => handleToggle(record.id)}
            disabled={record.status === 'completed'}
            style={{ borderRadius: 6, fontSize: 13 }}
          >
            {t('toggleStatus')}
          </Button>
          <Link to={`/console/tasks/${record.id}`}>
            <Button size="small" type="text" icon={<EyeOutlined />} style={{ borderRadius: 6, fontSize: 13 }}>
              {t('viewDetail')}
            </Button>
          </Link>
        </Space>
      ),
    },
  ]

  return (
    <div className="tf-table-wrap" style={{ marginTop: 24 }}>
      <div className="tf-table-toolbar">
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
      </div>
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
