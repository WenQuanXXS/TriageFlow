import { useState } from 'react'
import { Form, Input, Select, Button, Card, message } from 'antd'
import { useNavigate } from 'react-router-dom'
import { createTask } from '../api/tasks'
import { useLocale } from '../locales'

const { TextArea } = Input

export default function TaskForm() {
  const [loading, setLoading] = useState(false)
  const [form] = Form.useForm()
  const navigate = useNavigate()
  const { t } = useLocale()

  const priorityOptions = [
    { value: 'urgent', label: t('urgent') },
    { value: 'high', label: t('high') },
    { value: 'normal', label: t('normal') },
    { value: 'low', label: t('low') },
  ]

  const departmentOptions = [
    { value: 'Internal Medicine', label: t('internalMedicine') },
    { value: 'Surgery', label: t('surgery') },
    { value: 'Pediatrics', label: t('pediatrics') },
    { value: 'Neurology', label: t('neurology') },
    { value: 'Emergency', label: t('emergency') },
  ]

  const onFinish = async (values) => {
    setLoading(true)
    try {
      await createTask(values)
      message.success(t('taskCreatedSuccess'))
      form.resetFields()
      navigate('/')
    } catch (err) {
      message.error(t('taskCreatedFail'))
    } finally {
      setLoading(false)
    }
  }

  return (
    <Card title={t('newTriageTask')} style={{ maxWidth: 600, margin: '0 auto' }}>
      <Form form={form} layout="vertical" onFinish={onFinish}>
        <Form.Item
          name="patient_name"
          label={t('patientName')}
          rules={[{ required: true, message: t('pleaseEnterPatientName') }]}
        >
          <Input placeholder={t('enterPatientName')} />
        </Form.Item>

        <Form.Item
          name="chief_complaint"
          label={t('chiefComplaint')}
          rules={[{ required: true, message: t('pleaseEnterChiefComplaint') }]}
        >
          <TextArea rows={3} placeholder={t('enterChiefComplaint')} />
        </Form.Item>

        <Form.Item name="priority" label={t('priority')} initialValue="normal">
          <Select options={priorityOptions} />
        </Form.Item>

        <Form.Item name="department" label={t('department')}>
          <Select options={departmentOptions} allowClear placeholder={t('selectDepartment')} />
        </Form.Item>

        <Form.Item>
          <Button type="primary" htmlType="submit" loading={loading} block>
            {t('submit')}
          </Button>
        </Form.Item>
      </Form>
    </Card>
  )
}
