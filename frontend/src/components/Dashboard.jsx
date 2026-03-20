import { useState, useEffect } from 'react'
import { Card, Col, Row, Statistic, message } from 'antd'
import { getDashboard } from '../api/tasks'
import { useLocale } from '../locales'

export default function Dashboard() {
  const [data, setData] = useState(null)
  const { t } = useLocale()

  useEffect(() => {
    getDashboard()
      .then((res) => setData(res.data))
      .catch(() => message.error(t('loadDashboardFail')))
  }, [t])

  if (!data) return null

  const statusMap = {}
  for (const s of data.by_status || []) {
    statusMap[s.status] = s.count
  }

  return (
    <Row gutter={[16, 16]}>
      <Col span={4}>
        <Card>
          <Statistic title={t('total')} value={data.total} />
        </Card>
      </Col>
      <Col span={4}>
        <Card>
          <Statistic title={t('pending')} value={statusMap.pending || 0} valueStyle={{ color: '#fa8c16' }} />
        </Card>
      </Col>
      <Col span={4}>
        <Card>
          <Statistic title={t('inProgress')} value={statusMap.in_progress || 0} valueStyle={{ color: '#1890ff' }} />
        </Card>
      </Col>
      <Col span={4}>
        <Card>
          <Statistic title={t('completed')} value={statusMap.completed || 0} valueStyle={{ color: '#52c41a' }} />
        </Card>
      </Col>
      <Col span={4}>
        <Card>
          <Statistic title={t('triaged')} value={data.triage_count || 0} valueStyle={{ color: '#722ed1' }} />
        </Card>
      </Col>
      <Col span={4}>
        <Card>
          <Statistic title={t('ruleOverrides')} value={data.rule_overrides || 0} valueStyle={{ color: '#eb2f96' }} />
        </Card>
      </Col>
    </Row>
  )
}
