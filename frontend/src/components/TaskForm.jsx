import { useState } from 'react'
import { Form, Input, Button, Card, message, Spin, InputNumber, Select, Slider } from 'antd'
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

          <Form.Item name="age" label={t('age')}>
            <InputNumber min={0} max={150} style={{ width: '100%' }} placeholder={t('enterAge')} />
          </Form.Item>

          <Form.Item name="gender" label={t('gender')}>
            <Select placeholder={t('selectGender')} allowClear>
              <Select.Option value="male">{t('male')}</Select.Option>
              <Select.Option value="female">{t('female')}</Select.Option>
              <Select.Option value="other">{t('otherGender')}</Select.Option>
            </Select>
          </Form.Item>

          <Form.Item name="temperature" label={t('temperature')}>
            <InputNumber min={35.0} max={42.0} step={0.1} style={{ width: '100%' }} placeholder={t('enterTemperature')} />
          </Form.Item>

          <Form.Item name="pain_level" label={t('painLevel')}>
            <Slider min={0} max={10} marks={{ 0: '0', 5: '5', 10: '10' }} />
          </Form.Item>

          <Form.Item name="special_condition" label={t('specialCondition')}>
            <TextArea rows={2} placeholder={t('enterSpecialCondition')} />
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
