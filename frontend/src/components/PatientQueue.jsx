import { useState, useCallback } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { Card, Button, Spin, message } from 'antd'
import {
  ArrowLeftOutlined,
  ClockCircleOutlined,
  CheckCircleOutlined,
  BellOutlined,
} from '@ant-design/icons'
import { getQueuePosition } from '../api/tasks'
import { useLocale } from '../locales'
import { usePolling } from '../hooks/usePolling'

const patientStatusConfig = {
  waiting: { key: 'statusWaiting', color: '#e67e22', icon: <ClockCircleOutlined /> },
  called: { key: 'yourTurn', color: '#2e9b6e', icon: <BellOutlined /> },
  completed: { key: 'visitComplete', color: '#8e9aab', icon: <CheckCircleOutlined /> },
}

function formatNo(id) {
  return `#${String(id).padStart(3, '0')}`
}

export default function PatientQueue() {
  const { id } = useParams()
  const navigate = useNavigate()
  const { t, tDept } = useLocale()
  const [position, setPosition] = useState(null)
  const [loading, setLoading] = useState(true)

  const fetchData = useCallback(async () => {
    try {
      const { data } = await getQueuePosition(id)
      setPosition(data)
    } catch {
      message.error(t('loadDetailFail'))
    } finally {
      setLoading(false)
    }
  }, [id, t])

  usePolling(fetchData, 10000)

  if (loading) return <Spin style={{ display: 'block', margin: '100px auto' }} />
  if (!position) return null

  const statusCfg = patientStatusConfig[position.queue_status] || patientStatusConfig.waiting
  const nowServing = position.now_serving || []
  const waitingList = position.waiting_list || []

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
            <div className="tf-queue-number">{formatNo(position.task_id)}</div>
          </div>
          <div className="tf-queue-my-right">
            <div
              className={`tf-queue-status-badge ${position.queue_status === 'called' ? 'active' : ''}`}
              style={{ color: statusCfg.color, borderColor: statusCfg.color }}
            >
              {statusCfg.icon}
              <span>{t(statusCfg.key)}</span>
            </div>
          </div>
        </div>

        <div className="tf-queue-info-grid">
          {position.department && (
            <div className="tf-queue-info-item">
              <div className="tf-queue-info-label">{t('assignedDept')}</div>
              <div className="tf-queue-info-value">{tDept(position.department)}</div>
            </div>
          )}
          {position.queue_status === 'waiting' && (
            <div className="tf-queue-info-item">
              <div className="tf-queue-info-label">{t('queuePosition')}</div>
              <div className="tf-queue-info-value">
                {position.ahead === 0 ? t('noOneAhead') : `${position.ahead}${t('peopleAhead')}`}
              </div>
            </div>
          )}
        </div>

        {position.queue_status === 'waiting' && (
          <div className="tf-queue-hint">
            <ClockCircleOutlined /> {t('waitingForCall')}
          </div>
        )}

        {position.queue_status === 'called' && (
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
                nowServing.map((entry) => (
                  <div
                    key={entry.id}
                    className={`tf-queue-board-number serving ${entry.task_id === Number(id) ? 'is-me' : ''}`}
                  >
                    {t('pleaseGoTo', { no: formatNo(entry.task_id), dept: tDept(entry.department) })}
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
                waitingList.map((entry) => (
                  <div
                    key={entry.id}
                    className={`tf-queue-board-number ${entry.task_id === Number(id) ? 'is-me' : ''}`}
                  >
                    {formatNo(entry.task_id)}
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
