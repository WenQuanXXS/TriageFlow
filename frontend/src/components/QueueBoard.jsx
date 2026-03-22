import { useState, useCallback } from 'react'
import { Card, Spin, message } from 'antd'
import { listQueue } from '../api/tasks'
import { useLocale } from '../locales'
import { usePolling } from '../hooks/usePolling'

function formatNo(id) {
  return `#${String(id).padStart(3, '0')}`
}

export default function QueueBoard() {
  const { t, tDept } = useLocale()
  const [entries, setEntries] = useState(null)
  const [loading, setLoading] = useState(true)

  const fetchData = useCallback(async () => {
    try {
      const { data } = await listQueue()
      setEntries(data)
    } catch {
      message.error(t('loadDetailFail'))
    } finally {
      setLoading(false)
    }
  }, [t])

  usePolling(fetchData, 5000)

  if (loading) return <Spin style={{ display: 'block', margin: '100px auto' }} />
  if (!entries) return null

  const nowServing = entries.filter((e) => e.queue_status === 'called')
  const waitingList = entries.filter((e) => e.queue_status === 'waiting')

  return (
    <div style={{ maxWidth: 800, margin: '0 auto' }}>
      <Card className="tf-card tf-queue-board-card" title={t('queueBoard')}>
        <div className="tf-queue-board">
          <div className="tf-queue-board-section now-serving">
            <div className="tf-queue-board-title">{t('nowServing')}</div>
            <div className="tf-queue-board-numbers">
              {nowServing.length > 0 ? (
                nowServing.map((entry) => (
                  <div key={entry.id} className="tf-queue-board-number serving">
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
                  <div key={entry.id} className="tf-queue-board-number">
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
