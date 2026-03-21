import { useState, useCallback } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { Card, Button, Spin, message } from 'antd'
import {
  ArrowLeftOutlined,
  UserOutlined,
  ClockCircleOutlined,
  CheckCircleOutlined,
  BellOutlined,
} from '@ant-design/icons'
import { listTasks, getTaskDetail } from '../api/tasks'
import { useLocale } from '../locales'
import { usePolling } from '../hooks/usePolling'

const patientStatusConfig = {
  pending: { key: 'statusWaiting', color: '#e67e22', icon: <ClockCircleOutlined /> },
  in_progress: { key: 'yourTurn', color: '#2e9b6e', icon: <BellOutlined /> },
  completed: { key: 'visitComplete', color: '#8e9aab', icon: <CheckCircleOutlined /> },
}

function formatNo(id) {
  return `#${String(id).padStart(3, '0')}`
}

export default function PatientQueue() {
  const { id } = useParams()
  const navigate = useNavigate()
  const { t, tDept } = useLocale()
  const [myTask, setMyTask] = useState(null)
  const [allTasks, setAllTasks] = useState([])
  const [loading, setLoading] = useState(true)

  const fetchData = useCallback(async () => {
    try {
      const [taskRes, listRes] = await Promise.all([
        getTaskDetail(id),
        listTasks(),
      ])
      setMyTask(taskRes.data)
      setAllTasks(listRes.data)
    } catch {
      message.error(t('loadDetailFail'))
    } finally {
      setLoading(false)
    }
  }, [id, t])

  usePolling(fetchData, 10000)

  if (loading) return <Spin style={{ display: 'block', margin: '100px auto' }} />
  if (!myTask) return null

  const statusCfg = patientStatusConfig[myTask.status] || patientStatusConfig.pending

  // Calculate queue position: pending tasks created before this one
  const pendingTasks = allTasks
    .filter((t) => t.status === 'pending')
    .sort((a, b) => new Date(a.created_at) - new Date(b.created_at))
  const myIndex = pendingTasks.findIndex((t) => t.id === myTask.id)
  const ahead = myIndex >= 0 ? myIndex : 0

  // Queue board data
  const nowServing = allTasks
    .filter((t) => t.status === 'in_progress')
    .sort((a, b) => new Date(b.updated_at) - new Date(a.updated_at))
    .slice(0, 3)

  const waitingList = pendingTasks.slice(0, 8)

  return (
    <div style={{ maxWidth: 800, margin: '0 auto' }}>
      <div className="tf-detail-header">
        <Button
          className="tf-back-btn"
          icon={<ArrowLeftOutlined />}
          onClick={() => navigate('/patient')}
        >
          {t('back')}
        </Button>
      </div>

      {/* My Status Card */}
      <Card className="tf-card tf-queue-my-status">
        <div className="tf-queue-my-header">
          <div>
            <div className="tf-queue-label">{t('yourNumber')}</div>
            <div className="tf-queue-number">{formatNo(myTask.id)}</div>
          </div>
          <div className="tf-queue-my-right">
            <div
              className={`tf-queue-status-badge ${myTask.status === 'in_progress' ? 'active' : ''}`}
              style={{ color: statusCfg.color, borderColor: statusCfg.color }}
            >
              {statusCfg.icon}
              <span>{t(statusCfg.key)}</span>
            </div>
          </div>
        </div>

        <div className="tf-queue-info-grid">
          {myTask.final_department && (
            <div className="tf-queue-info-item">
              <div className="tf-queue-info-label">{t('assignedDept')}</div>
              <div className="tf-queue-info-value">{tDept(myTask.final_department)}</div>
            </div>
          )}
          {myTask.status === 'pending' && (
            <div className="tf-queue-info-item">
              <div className="tf-queue-info-label">{t('queuePosition')}</div>
              <div className="tf-queue-info-value">
                {ahead === 0 ? t('noOneAhead') : `${ahead}${t('peopleAhead')}`}
              </div>
            </div>
          )}
        </div>

        {myTask.status === 'pending' && (
          <div className="tf-queue-hint">
            <ClockCircleOutlined /> {t('waitingForCall')}
          </div>
        )}

        {myTask.status === 'in_progress' && (
          <div className="tf-queue-alert">
            <BellOutlined /> {t('yourTurn')}
          </div>
        )}
      </Card>

      {/* Queue Display Board */}
      <Card className="tf-card tf-queue-board-card" title={t('queueBoard')}>
        <div className="tf-queue-board">
          <div className="tf-queue-board-section now-serving">
            <div className="tf-queue-board-title">{t('nowServing')}</div>
            <div className="tf-queue-board-numbers">
              {nowServing.length > 0 ? (
                nowServing.map((task) => (
                  <div
                    key={task.id}
                    className={`tf-queue-board-number serving ${task.id === myTask.id ? 'is-me' : ''}`}
                  >
                    {formatNo(task.id)}
                  </div>
                ))
              ) : (
                <div className="tf-queue-board-empty">—</div>
              )}
            </div>
          </div>
          <div className="tf-queue-board-section waiting">
            <div className="tf-queue-board-title">{t('waitingList')}</div>
            <div className="tf-queue-board-numbers">
              {waitingList.length > 0 ? (
                waitingList.map((task) => (
                  <div
                    key={task.id}
                    className={`tf-queue-board-number ${task.id === myTask.id ? 'is-me' : ''}`}
                  >
                    {formatNo(task.id)}
                  </div>
                ))
              ) : (
                <div className="tf-queue-board-empty">{t('noWaitingPatients')}</div>
              )}
            </div>
          </div>
        </div>
      </Card>
    </div>
  )
}
