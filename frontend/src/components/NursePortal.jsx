import { useState, useCallback } from 'react'
import { Table, Button, Select, Space, Tag, message } from 'antd'
import { ReloadOutlined, PhoneOutlined, CheckOutlined } from '@ant-design/icons'
import { listTasks, toggleTaskStatus } from '../api/tasks'
import { useLocale } from '../locales'
import PriorityBadge, { priorityLabelKeys } from './shared/PriorityBadge'
import { usePolling } from '../hooks/usePolling'

const nurseStatusColors = {
  pending: '#e67e22',
  in_progress: '#3498db',
  completed: '#27ae60',
}

const nurseStatusKeys = {
  pending: 'statusWaiting',
  in_progress: 'statusInConsultation',
  completed: 'statusCompleted',
}

export default function NursePortal() {
  const [tasks, setTasks] = useState([])
  const [loading, setLoading] = useState(false)
  const [statusFilter, setStatusFilter] = useState('')
  const [lastUpdated, setLastUpdated] = useState(null)
  const { t, tDept } = useLocale()

  const fetchTasks = useCallback(async () => {
    try {
      const { data } = await listTasks()
      // Sort: in_progress first, then pending (by created_at asc), then completed
      const sorted = [...data].sort((a, b) => {
        const order = { in_progress: 0, pending: 1, completed: 2 }
        if (order[a.status] !== order[b.status]) return order[a.status] - order[b.status]
        return new Date(a.created_at) - new Date(b.created_at)
      })
      setTasks(sorted)
      setLastUpdated(new Date())
    } catch {
      message.error(t('loadTasksFail'))
    }
  }, [t])

  usePolling(fetchTasks, 5000)

  const handleCall = async (id) => {
    setLoading(true)
    try {
      await toggleTaskStatus(id)
      message.success(`#${String(id).padStart(3, '0')} ${t('callPatient')}`)
      fetchTasks()
    } catch {
      message.error(t('updateStatusFail'))
    } finally {
      setLoading(false)
    }
  }

  const handleComplete = async (id) => {
    setLoading(true)
    try {
      await toggleTaskStatus(id)
      fetchTasks()
    } catch {
      message.error(t('updateStatusFail'))
    } finally {
      setLoading(false)
    }
  }

  const filtered = statusFilter ? tasks.filter((t) => t.status === statusFilter) : tasks

  const statusOptions = [
    { value: '', label: t('allStatus') },
    { value: 'pending', label: t('statusWaiting') },
    { value: 'in_progress', label: t('statusInConsultation') },
    { value: 'completed', label: t('completed') },
  ]

  const columns = [
    {
      title: t('patientNo'),
      dataIndex: 'id',
      key: 'id',
      width: 100,
      render: (id) => (
        <span style={{ fontWeight: 700, fontSize: 15, color: '#1a2332' }}>
          #{String(id).padStart(3, '0')}
        </span>
      ),
    },
    {
      title: t('patient'),
      dataIndex: 'patient_name',
      key: 'patient_name',
      render: (name) => <span style={{ fontWeight: 500 }}>{name}</span>,
    },
    {
      title: t('chiefComplaint'),
      dataIndex: 'chief_complaint',
      key: 'chief_complaint',
      ellipsis: true,
      render: (text) => <span style={{ color: '#5a6a7e' }}>{text}</span>,
    },
    {
      title: t('finalPriority'),
      key: 'final_priority',
      width: 120,
      render: (_, record) => {
        const pri = record.final_priority || record.priority
        return <PriorityBadge priority={pri} label={t(priorityLabelKeys[pri])} />
      },
    },
    {
      title: t('finalDepartment'),
      dataIndex: 'final_department',
      key: 'final_department',
      width: 130,
      render: (dept) => <span style={{ fontWeight: 500 }}>{tDept(dept)}</span>,
    },
    {
      title: t('status'),
      dataIndex: 'status',
      key: 'status',
      width: 120,
      render: (s) => (
        <Tag color={nurseStatusColors[s]} style={{ borderRadius: 6, fontWeight: 500 }}>
          {t(nurseStatusKeys[s])}
        </Tag>
      ),
    },
    {
      title: t('action'),
      key: 'action',
      width: 120,
      render: (_, record) => {
        if (record.status === 'pending') {
          return (
            <Button
              type="primary"
              size="small"
              icon={<PhoneOutlined />}
              onClick={() => handleCall(record.id)}
              loading={loading}
              className="tf-nurse-call-btn"
            >
              {t('callPatient')}
            </Button>
          )
        }
        if (record.status === 'in_progress') {
          return (
            <Button
              size="small"
              icon={<CheckOutlined />}
              onClick={() => handleComplete(record.id)}
              loading={loading}
            >
              {t('completeVisit')}
            </Button>
          )
        }
        return <CheckOutlined style={{ color: '#27ae60' }} />
      },
    },
  ]

  return (
    <div>
      <div className="tf-nurse-toolbar">
        <h2 className="tf-page-title">{t('nurseQueue')}</h2>
        <Space>
          <Select
            value={statusFilter}
            onChange={setStatusFilter}
            options={statusOptions}
            style={{ width: 150 }}
          />
          <Button icon={<ReloadOutlined />} onClick={fetchTasks}>
            {t('refresh')}
          </Button>
          {lastUpdated && (
            <span className="tf-last-updated">
              {t('lastUpdated')}: {lastUpdated.toLocaleTimeString()}
            </span>
          )}
        </Space>
      </div>
      <div className="tf-table-wrap">
        <Table
          dataSource={filtered}
          columns={columns}
          rowKey="id"
          pagination={{ pageSize: 15 }}
          rowClassName={(record) =>
            record.status === 'in_progress' ? 'tf-nurse-row-active' : ''
          }
        />
      </div>
    </div>
  )
}
