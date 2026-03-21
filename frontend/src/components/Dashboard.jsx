import { useState, useEffect } from 'react'
import { Card, Col, Row, Statistic, message } from 'antd'
import {
  TeamOutlined,
  ClockCircleOutlined,
  SyncOutlined,
  CheckCircleOutlined,
  RobotOutlined,
  SafetyCertificateOutlined,
} from '@ant-design/icons'
import { getDashboard } from '../api/tasks'
import { useLocale } from '../locales'

const statConfig = [
  { key: 'total', icon: TeamOutlined, color: '#2e9b6e', bg: '#e8f5ef' },
  { key: 'pending', icon: ClockCircleOutlined, color: '#e67e22', bg: '#fef3e6' },
  { key: 'inProgress', icon: SyncOutlined, color: '#3498db', bg: '#ebf5fb' },
  { key: 'completed', icon: CheckCircleOutlined, color: '#27ae60', bg: '#e8f8f0' },
  { key: 'triaged', icon: RobotOutlined, color: '#8e44ad', bg: '#f4ecf7' },
  { key: 'ruleOverrides', icon: SafetyCertificateOutlined, color: '#e74c3c', bg: '#fdecea' },
]

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

  const values = {
    total: data.total,
    pending: statusMap.pending || 0,
    inProgress: statusMap.in_progress || 0,
    completed: statusMap.completed || 0,
    triaged: data.triage_count || 0,
    ruleOverrides: data.rule_overrides || 0,
  }

  return (
    <Row gutter={[16, 16]} style={{ marginBottom: 8 }}>
      {statConfig.map(({ key, icon: Icon, color, bg }) => (
        <Col xs={12} sm={8} md={4} key={key}>
          <Card className="tf-stat-card" bordered={false}>
            <div className="tf-stat-icon" style={{ background: bg, color }}>
              <Icon />
            </div>
            <Statistic
              title={t(key)}
              value={values[key]}
              valueStyle={{ color }}
            />
          </Card>
        </Col>
      ))}
    </Row>
  )
}
