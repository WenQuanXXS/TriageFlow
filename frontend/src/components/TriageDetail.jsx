import { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { Card, Tag, Alert, Button, Spin, Space, message } from 'antd'
import {
  ArrowLeftOutlined,
  UserOutlined,
  RobotOutlined,
  SafetyCertificateOutlined,
  WarningOutlined,
} from '@ant-design/icons'
import { getTaskDetail } from '../api/tasks'
import { useLocale } from '../locales'
import StatusBadge from './shared/StatusBadge'

const priorityColors = {
  urgent: '#e74c3c',
  high: '#e67e22',
  normal: '#2e9b6e',
  low: '#8e9aab',
}

function safeParse(jsonStr) {
  if (!jsonStr) return []
  try {
    const result = JSON.parse(jsonStr)
    return Array.isArray(result) ? result : []
  } catch {
    return []
  }
}

function DetailField({ label, value, full, children }) {
  return (
    <div className={`tf-detail-field${full ? ' full' : ''}`}>
      <div className="tf-detail-field-label">{label}</div>
      <div className="tf-detail-field-value">{children || value || '—'}</div>
    </div>
  )
}


export default function TriageDetail() {
  const { id } = useParams()
  const navigate = useNavigate()
  const { t, tDept } = useLocale()
  const [task, setTask] = useState(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    getTaskDetail(id)
      .then((res) => setTask(res.data))
      .catch(() => message.error(t('loadDetailFail')))
      .finally(() => setLoading(false))
  }, [id, t])

  if (loading) return <Spin style={{ display: 'block', margin: '100px auto' }} />
  if (!task) return null

  const symptoms = safeParse(task.symptoms)
  const riskSignals = safeParse(task.risk_signals)
  const candidateDepts = safeParse(task.candidate_depts)
  const hasRuleOverride = !!task.rule_triggered
  const finalPri = task.final_priority || 'normal'
  const genderKey = task.gender === 'male' ? 'male' : task.gender === 'female' ? 'female' : 'otherGender'

  return (
    <div style={{ maxWidth: 900, margin: '0 auto' }}>
      {/* Header */}
      <div className="tf-detail-header">
        <Button
          className="tf-back-btn"
          icon={<ArrowLeftOutlined />}
          onClick={() => navigate('/console')}
        >
          {t('back')}
        </Button>
      </div>

      {/* Rule Override Alert */}
      {hasRuleOverride && (
        <Alert
          className="tf-rule-alert"
          type="warning"
          showIcon
          icon={<WarningOutlined />}
          message={`${t('ruleOverrideAlert')}：${task.rule_triggered}`}
          description={task.rule_reason}
        />
      )}

      {/* Priority Strip */}
      <div className={`tf-detail-priority-strip ${finalPri}`}>
        <div>
          <div className="tf-detail-priority-label">{t('finalPriority')}</div>
          <div className="tf-detail-priority-value" style={{ color: priorityColors[finalPri] }}>
            {finalPri.toUpperCase()}
            {hasRuleOverride && task.ai_suggested_priority !== task.final_priority && (
              <span style={{ fontSize: 13, fontWeight: 400, marginLeft: 12, opacity: 0.7 }}>
                {task.ai_suggested_priority} → {task.final_priority}
              </span>
            )}
          </div>
        </div>
        <div style={{ textAlign: 'right' }}>
          <div className="tf-detail-priority-label">{t('finalDepartment')}</div>
          <div style={{ fontSize: 18, fontWeight: 600, color: '#1a2332' }}>
            {tDept(task.final_department)}
          </div>
        </div>
      </div>

      {/* Patient Info Card */}
      <Card
        className="tf-card tf-detail-section"
        title={
          <>
            <span className="tf-detail-section-icon" style={{ background: '#e8f5ef', color: '#2e9b6e' }}>
              <UserOutlined />
            </span>
            {t('patientName')}
          </>
        }
      >
        <div className="tf-detail-grid">
          <DetailField label={t('patientName')} value={task.patient_name} />
          <DetailField label={t('status')}>
            <StatusBadge status={task.status} />
          </DetailField>
          <DetailField label={t('chiefComplaint')} full value={task.chief_complaint} />
          {task.age > 0 && <DetailField label={t('age')} value={task.age} />}
          {task.gender && <DetailField label={t('gender')} value={t(genderKey)} />}
          {task.temperature > 0 && <DetailField label={t('temperature')} value={`${task.temperature}°C`} />}
          {task.pain_level > 0 && <DetailField label={t('painLevel')} value={`${task.pain_level}/10`} />}
          {task.special_condition && (
            <DetailField label={t('specialCondition')} full value={task.special_condition} />
          )}
          <DetailField label={t('created')} full value={new Date(task.created_at).toLocaleString()} />
        </div>
      </Card>

      {/* AI Analysis Card */}
      <Card
        className="tf-card tf-detail-section"
        title={
          <>
            <span className="tf-detail-section-icon" style={{ background: '#f4ecf7', color: '#8e44ad' }}>
              <RobotOutlined />
            </span>
            {t('triageDetail')}
          </>
        }
      >
        <div className="tf-detail-grid">
          <DetailField label={t('symptoms')} full>
            <Space size={[6, 6]} wrap>
              {symptoms.map((s, i) => (
                <Tag key={i} style={{ borderRadius: 6, border: '1px solid #d4eddf', background: '#e8f5ef', color: '#237a56' }}>
                  {s}
                </Tag>
              ))}
            </Space>
          </DetailField>

          <DetailField label={t('riskSignals')} full>
            {riskSignals.length > 0 ? (
              <Space size={[6, 6]} wrap>
                {riskSignals.map((s, i) => (
                  <Tag key={i} style={{ borderRadius: 6, border: '1px solid #f5c6cb', background: '#fdecea', color: '#c0392b' }}>
                    {s}
                  </Tag>
                ))}
              </Space>
            ) : (
              <span style={{ color: '#8e9aab' }}>{t('noRiskSignals')}</span>
            )}
          </DetailField>

          <DetailField label={t('candidateDepts')} full>
            <Space size={[6, 6]} wrap>
              {candidateDepts.map((d, i) => (
                <Tag key={i} style={{ borderRadius: 6, border: '1px solid #bee5eb', background: '#ebf5fb', color: '#2471a3' }}>
                  {tDept(d)}
                </Tag>
              ))}
            </Space>
          </DetailField>

          <DetailField label={t('aiSuggestedPri')}>
            <span
              className={`tf-priority-badge ${task.ai_suggested_priority}`}
              style={{ fontSize: 13 }}
            >
              {task.ai_suggested_priority?.toUpperCase()}
            </span>
          </DetailField>

          <DetailField label={t('triageStatus')}>
            <StatusBadge status={task.triage_status === 'completed' ? 'completed' : 'pending'} />
          </DetailField>
        </div>
      </Card>

      {/* Rule Engine Card (only if triggered) */}
      {hasRuleOverride && (
        <Card
          className="tf-card tf-detail-section"
          title={
            <>
              <span className="tf-detail-section-icon" style={{ background: '#fef3e6', color: '#e67e22' }}>
                <SafetyCertificateOutlined />
              </span>
              {t('ruleTriggered')}
            </>
          }
        >
          <div className="tf-detail-grid">
            <DetailField label={t('ruleTriggered')} value={task.rule_triggered} />
            <DetailField label={t('ruleReason')} value={task.rule_reason} />
          </div>
        </Card>
      )}
    </div>
  )
}
