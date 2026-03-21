import { useState } from 'react'
import { Form, Input, Button, Card, message, Spin, InputNumber, Select, Slider, Row, Col } from 'antd'
import { useNavigate } from 'react-router-dom'
import { UserOutlined, HeartOutlined, SendOutlined } from '@ant-design/icons'
import { createTask } from '../api/tasks'
import { useLocale } from '../locales'

const { TextArea } = Input

export default function TaskForm({ title, onSuccess }) {
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
      if (onSuccess) {
        onSuccess(data)
      } else {
        navigate(`/console/tasks/${data.id}`)
      }
    } catch (err) {
      message.error(t('taskCreatedFail'))
    } finally {
      setLoading(false)
    }
  }

  return (
    <Card
      className="tf-form-card"
      title={title || t('newTriageTask')}
      style={{ maxWidth: 680, margin: '0 auto' }}
    >
      <Spin spinning={loading} className="tf-analyzing-spin" tip={t('analyzing')}>
        <Form form={form} layout="vertical" onFinish={onFinish} size="large">
          {/* Section: Basic Info */}
          <div className="tf-form-section">
            <span className="tf-form-section-icon"><UserOutlined /></span>
            {t('patientName')}
          </div>

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

          {/* Section: Demographics & Vitals */}
          <div className="tf-form-section" style={{ marginTop: 28 }}>
            <span className="tf-form-section-icon"><HeartOutlined /></span>
            {t('age')} / {t('gender')} / {t('temperature')}
          </div>

          <Row gutter={16}>
            <Col span={8}>
              <Form.Item name="age" label={t('age')}>
                <InputNumber
                  min={0}
                  max={150}
                  style={{ width: '100%' }}
                  placeholder={t('enterAge')}
                />
              </Form.Item>
            </Col>
            <Col span={8}>
              <Form.Item name="gender" label={t('gender')}>
                <Select placeholder={t('selectGender')} allowClear>
                  <Select.Option value="male">{t('male')}</Select.Option>
                  <Select.Option value="female">{t('female')}</Select.Option>
                  <Select.Option value="other">{t('otherGender')}</Select.Option>
                </Select>
              </Form.Item>
            </Col>
            <Col span={8}>
              <Form.Item name="temperature" label={t('temperature')}>
                <InputNumber
                  min={35.0}
                  max={42.0}
                  step={0.1}
                  style={{ width: '100%' }}
                  placeholder={t('enterTemperature')}
                />
              </Form.Item>
            </Col>
          </Row>

          <Form.Item
            name="pain_level"
            label={t('painLevel')}
            style={{ marginBottom: 28 }}
          >
            <Slider
              min={0}
              max={10}
              marks={{ 0: '0', 2: '2', 4: '4', 6: '6', 8: '8', 10: '10' }}
            />
          </Form.Item>

          <Form.Item name="special_condition" label={t('specialCondition')}>
            <TextArea rows={2} placeholder={t('enterSpecialCondition')} />
          </Form.Item>

          <Form.Item style={{ marginTop: 32, marginBottom: 0 }}>
            <Button
              className="tf-submit-btn"
              type="primary"
              htmlType="submit"
              loading={loading}
              block
              icon={<SendOutlined />}
            >
              {t('submit')}
            </Button>
          </Form.Item>
        </Form>
      </Spin>
    </Card>
  )
}
