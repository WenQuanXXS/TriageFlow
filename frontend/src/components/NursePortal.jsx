import { useState, useCallback } from 'react'
import { Table, Button, Select, Space, Tag, message } from 'antd'
import { ReloadOutlined, PhoneOutlined, CheckOutlined } from '@ant-design/icons'
import { listQueue, callQueuePatient, completeQueuePatient } from '../api/tasks'
import { useLocale } from '../locales'
import PriorityBadge, { priorityLabelKeys } from './shared/PriorityBadge'
import { usePolling } from '../hooks/usePolling'

const queueStatusColors = {
  waiting: '#e67e22',
  called: '#3498db',
  completed: '#27ae60',
}

const queueStatusKeys = {
  waiting: 'statusWaiting',
  called: 'statusCalled',
  completed: 'statusCompleted',
}

export default function NursePortal() {
  const [entries, setEntries] = useState([])
  const [loading, setLoading] = useState(false)
  const [statusFilter, setStatusFilter] = useState('')
  const [lastUpdated, setLastUpdated] = useState(null)
  const { t, tDept } = useLocale()

  const fetchQueue = useCallback(async () => {
    try {
      const params = {}
      if (statusFilter) params.queue_status = statusFilter
      const { data } = await listQueue(params)
      setEntries(data)
      setLastUpdated(new Date())
    } catch {
      message.error(t('loadTasksFail'))
    }
  }, [statusFilter, t])

  usePolling(fetchQueue, 5000)

  const handleCall = async (taskId) => {
    setLoading(true)
    try {
      await callQueuePatient(taskId)
      message.success(`#${String(taskId).padStart(3, '0')} ${t('callPatient')}`)
      fetchQueue()
    } catch {
      message.error(t('updateStatusFail'))
    } finally {
      setLoading(false)
    }
  }

  const handleComplete = async (taskId) => {
    setLoading(true)
    try {
      await completeQueuePatient(taskId)
      fetchQueue()
    } catch {
      message.error(t('updateStatusFail'))
    } finally {
      setLoading(false)
    }
  }

  const statusOptions = [
    { value: '', label: t('allStatus') },
    { value: 'waiting', label: t('statusWaiting') },
    { value: 'called', label: t('statusCalled') },
    { value: 'completed', label: t('completed') },
  ]

  const columns = [
    {
      title: t('patientNo'),
      dataIndex: 'queue_number',
      key: 'queue_number',
      width: 100,
      render: (num, record) => (
        <span style={{ fontWeight: 700, fontSize: 15, color: '#1a2332' }}>
          #{String(record.task_id).padStart(3, '0')}
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
      dataIndex: 'priority',
      key: 'priority',
      width: 120,
      render: (pri) => <PriorityBadge priority={pri} label={t(priorityLabelKeys[pri])} />,
    },
    {
      title: t('finalDepartment'),
      dataIndex: 'department',
      key: 'department',
      width: 130,
      render: (dept) => <span style={{ fontWeight: 500 }}>{tDept(dept)}</span>,
    },
    {
      title: t('status'),
      dataIndex: 'queue_status',
      key: 'queue_status',
      width: 120,
      render: (s) => (
        <Tag color={queueStatusColors[s]} style={{ borderRadius: 6, fontWeight: 500 }}>
          {t(queueStatusKeys[s])}
        </Tag>
      ),
    },
    {
      title: t('action'),
      key: 'action',
      width: 120,
      render: (_, record) => {
        if (record.queue_status === 'waiting') {
          return (
            <Button
              type="primary"
              size="small"
              icon={<PhoneOutlined />}
              onClick={() => handleCall(record.task_id)}
              loading={loading}
              className="tf-nurse-call-btn"
            >
              {t('callPatient')}
            </Button>
          )
        }
        if (record.queue_status === 'called') {
          return (
            <Button
              size="small"
              icon={<CheckOutlined />}
              onClick={() => handleComplete(record.task_id)}
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
          <Button icon={<ReloadOutlined />} onClick={fetchQueue}>
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
          dataSource={entries}
          columns={columns}
          rowKey="id"
          pagination={{ pageSize: 15 }}
          rowClassName={(record) =>
            record.queue_status === 'called' ? 'tf-nurse-row-active' : ''
          }
        />
      </div>
    </div>
  )
}
