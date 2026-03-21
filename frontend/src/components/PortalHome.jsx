import { Link } from 'react-router-dom'
import { Button } from 'antd'
import {
  ControlOutlined,
  MedicineBoxOutlined,
  UserOutlined,
  TranslationOutlined,
} from '@ant-design/icons'
import { useLocale } from '../locales'

const portals = [
  {
    key: 'console',
    path: '/console',
    icon: <ControlOutlined />,
    titleKey: 'consolePortal',
    descKey: 'consolePortalDesc',
    color: '#2e9b6e',
    bg: '#e8f5ef',
  },
  {
    key: 'nurse',
    path: '/nurse',
    icon: <MedicineBoxOutlined />,
    titleKey: 'nursePortal',
    descKey: 'nursePortalDesc',
    color: '#2471a3',
    bg: '#ebf5fb',
  },
  {
    key: 'patient',
    path: '/patient',
    icon: <UserOutlined />,
    titleKey: 'patientPortal',
    descKey: 'patientPortalDesc',
    color: '#8e44ad',
    bg: '#f4ecf7',
  },
]

export default function PortalHome() {
  const { lang, t, toggleLang } = useLocale()

  return (
    <div className="tf-portal-page">
      <Button
        className="tf-portal-lang-btn"
        type="text"
        icon={<TranslationOutlined />}
        onClick={toggleLang}
      >
        {lang === 'zh' ? 'EN' : '中文'}
      </Button>

      <div className="tf-portal-hero">
        <div className="tf-portal-logo">
          <MedicineBoxOutlined />
        </div>
        <h1 className="tf-portal-title">
          Triage<span>Flow</span>
        </h1>
        <p className="tf-portal-subtitle">{t('portalSubtitle')}</p>
      </div>

      <div className="tf-portal-cards">
        {portals.map((p) => (
          <Link key={p.key} to={p.path} className="tf-portal-card">
            <div className="tf-portal-card-icon" style={{ background: p.bg, color: p.color }}>
              {p.icon}
            </div>
            <div className="tf-portal-card-title">{t(p.titleKey)}</div>
            <div className="tf-portal-card-desc">{t(p.descKey)}</div>
            <div className="tf-portal-card-enter" style={{ color: p.color }}>
              {t('enterPortal')} →
            </div>
          </Link>
        ))}
      </div>
    </div>
  )
}
