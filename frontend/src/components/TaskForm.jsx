import { useState } from 'react'
import { Form, Input, Button, Card, message, Spin } from 'antd'
import { useNavigate } from 'react-router-dom'
import { createTask } from '../api/tasks'
import { useLocale } from '../locales'

const { TextArea } = Input

export default function TaskForm() {
  const [loading, setLoading] = useState(false)
  const [form] = Form.useForm()
  const navigate = useNavigate()
  const { t } = useLocale()

  const onFinish = async (values) => {
    setLoading(true)
    try {
      const { data } = await createTask(values)
      message.success(t('taskCreatedSuccess'))
      form.resetFields()
      navigate(`/tasks/${data.id}`)
    } catch (err) {
      message.error(t('taskCreatedFail'))
    } finally {
      setLoading(false)
    }
  }

  return (
    <Card title={t('newTriageTask')} style={{ maxWidth: 600, margin: '0 auto' }}>
      <Spin spinning={loading} tip={t('analyzing')}>
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
            <TextArea rows={4} placeholder={t('enterChiefComplaint')} />
          </Form.Item>

          <Form.Item>
            <Button type="primary" htmlType="submit" loading={loading} block>
              {t('submit')}
            </Button>
          </Form.Item>
        </Form>
      </Spin>
    </Card>
  )
}
