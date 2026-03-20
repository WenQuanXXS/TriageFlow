import { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { Card, Descriptions, Tag, Alert, Button, Spin, message } from 'antd'
import { ArrowLeftOutlined } from '@ant-design/icons'
import { getTaskDetail } from '../api/tasks'
import { useLocale } from '../locales'

const priorityColors = {
  urgent: 'red',
  high: 'volcano',
  normal: 'blue',
  low: 'default',
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

  return (
    <div style={{ maxWidth: 800, margin: '0 auto' }}>
      <Button
        icon={<ArrowLeftOutlined />}
        onClick={() => navigate('/')}
        style={{ marginBottom: 16 }}
      >
        {t('back')}
      </Button>

      {hasRuleOverride && (
        <Alert
          type="warning"
          showIcon
          message={`${t('ruleOverrideAlert')}：${task.rule_triggered}`}
          description={task.rule_reason}
          style={{ marginBottom: 16 }}
        />
      )}

      <Card title={t('triageDetail')}>
        <Descriptions column={2} bordered>
          <Descriptions.Item label={t('patientName')}>{task.patient_name}</Descriptions.Item>
          <Descriptions.Item label={t('status')}>
            <Tag color={task.status === 'completed' ? 'green' : task.status === 'in_progress' ? 'blue' : 'orange'}>
              {task.status}
            </Tag>
          </Descriptions.Item>
          <Descriptions.Item label={t('chiefComplaint')} span={2}>{task.chief_complaint}</Descriptions.Item>

          {task.age > 0 && (
            <Descriptions.Item label={t('age')}>{task.age}</Descriptions.Item>
          )}
          {task.gender && (
            <Descriptions.Item label={t('gender')}>
              {t(task.gender === 'male' ? 'male' : task.gender === 'female' ? 'female' : 'otherGender')}
            </Descriptions.Item>
          )}
          {task.temperature > 0 && (
            <Descriptions.Item label={t('temperature')}>{task.temperature}°C</Descriptions.Item>
          )}
          {task.pain_level > 0 && (
            <Descriptions.Item label={t('painLevel')}>{task.pain_level}/10</Descriptions.Item>
          )}
          {task.special_condition && (
            <Descriptions.Item label={t('specialCondition')} span={2}>{task.special_condition}</Descriptions.Item>
          )}

          <Descriptions.Item label={t('symptoms')} span={2}>
            {symptoms.map((s, i) => <Tag key={i}>{s}</Tag>)}
          </Descriptions.Item>

          <Descriptions.Item label={t('riskSignals')} span={2}>
            {riskSignals.length > 0
              ? riskSignals.map((s, i) => <Tag key={i} color="red">{s}</Tag>)
              : <span style={{ color: '#999' }}>{t('noRiskSignals')}</span>
            }
          </Descriptions.Item>

          <Descriptions.Item label={t('candidateDepts')} span={2}>
            {candidateDepts.map((d, i) => <Tag key={i} color="blue">{tDept(d)}</Tag>)}
          </Descriptions.Item>

          <Descriptions.Item label={t('aiSuggestedPri')}>
            <Tag color={priorityColors[task.ai_suggested_priority]}>
              {task.ai_suggested_priority?.toUpperCase()}
            </Tag>
          </Descriptions.Item>

          <Descriptions.Item label={t('finalPriority')}>
            <Tag color={priorityColors[task.final_priority]}>
              {task.final_priority?.toUpperCase()}
            </Tag>
            {hasRuleOverride && task.ai_suggested_priority !== task.final_priority && (
              <span style={{ color: '#fa8c16', marginLeft: 8 }}>
                ({task.ai_suggested_priority} → {task.final_priority})
              </span>
            )}
          </Descriptions.Item>

          <Descriptions.Item label={t('finalDepartment')}>{tDept(task.final_department)}</Descriptions.Item>
          <Descriptions.Item label={t('triageStatus')}>
            <Tag color={task.triage_status === 'completed' ? 'green' : 'orange'}>
              {task.triage_status}
            </Tag>
          </Descriptions.Item>

          {hasRuleOverride && (
            <>
              <Descriptions.Item label={t('ruleTriggered')}>{task.rule_triggered}</Descriptions.Item>
              <Descriptions.Item label={t('ruleReason')}>{task.rule_reason}</Descriptions.Item>
            </>
          )}

          <Descriptions.Item label={t('created')} span={2}>
            {new Date(task.created_at).toLocaleString()}
          </Descriptions.Item>
        </Descriptions>
      </Card>
    </div>
  )
}
