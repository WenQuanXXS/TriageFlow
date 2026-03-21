import { Routes, Route, Link, useLocation, useNavigate } from 'react-router-dom'
import { Layout, Menu, Button, ConfigProvider } from 'antd'
import {
  PlusOutlined,
  UnorderedListOutlined,
  TranslationOutlined,
  MedicineBoxOutlined,
  HomeOutlined,
  UserOutlined,
  PhoneOutlined,
} from '@ant-design/icons'
import zhCN from 'antd/locale/zh_CN'
import enUS from 'antd/locale/en_US'
import TaskForm from './components/TaskForm'
import TaskList from './components/TaskList'
import Dashboard from './components/Dashboard'
import TriageDetail from './components/TriageDetail'
import PortalHome from './components/PortalHome'
import NursePortal from './components/NursePortal'
import PatientQueue from './components/PatientQueue'
import { useLocale } from './locales'

const { Header, Content } = Layout

const theme = {
  token: {
    colorPrimary: '#2e9b6e',
    colorSuccess: '#2e9b6e',
    colorLink: '#2e9b6e',
    borderRadius: 8,
    fontFamily: "'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif",
    colorBgContainer: '#ffffff',
    controlHeight: 38,
    fontSize: 14,
  },
  components: {
    Table: {
      headerBg: '#fafcfb',
      rowHoverBg: '#e8f5ef',
      borderColor: '#e2e8f0',
    },
    Card: {
      headerFontSize: 16,
    },
    Input: {
      controlHeight: 40,
    },
    Select: {
      controlHeight: 40,
    },
    Button: {
      controlHeight: 38,
    },
  },
}

function ConsolePage() {
  return (
    <>
      <Dashboard />
      <TaskList />
    </>
  )
}

function PatientRegister() {
  const navigate = useNavigate()
  const { t } = useLocale()

  return (
    <TaskForm
      title={t('patientRegistration')}
      onSuccess={(data) => navigate(`/patient/queue/${data.id}`)}
    />
  )
}

function getMenuItems(pathname, t) {
  if (pathname.startsWith('/console')) {
    return [
      { key: '/console', label: <Link to="/console">{t('dashboard')}</Link>, icon: <UnorderedListOutlined /> },
      { key: '/console/new', label: <Link to="/console/new">{t('newTask')}</Link>, icon: <PlusOutlined /> },
    ]
  }
  if (pathname.startsWith('/nurse')) {
    return [
      { key: '/nurse', label: <Link to="/nurse">{t('nurseQueue')}</Link>, icon: <PhoneOutlined /> },
    ]
  }
  if (pathname.startsWith('/patient')) {
    return [
      { key: '/patient', label: <Link to="/patient">{t('patientRegistration')}</Link>, icon: <UserOutlined /> },
    ]
  }
  return []
}

function getSelectedKey(pathname) {
  if (pathname.startsWith('/console/tasks')) return '/console'
  if (pathname.startsWith('/patient/queue')) return '/patient'
  return pathname
}

export default function App() {
  const { lang, t, toggleLang } = useLocale()
  const location = useLocation()

  const isPortalHome = location.pathname === '/'
  const menuItems = getMenuItems(location.pathname, t)

  return (
    <ConfigProvider locale={lang === 'zh' ? zhCN : enUS} theme={theme}>
      <Layout style={{ minHeight: '100vh', background: '#f4f7f5' }}>
        {!isPortalHome && (
          <Header className="tf-header" style={{ display: 'flex', alignItems: 'center' }}>
            <Link to="/" className="tf-brand">
              <div className="tf-brand-icon">
                <MedicineBoxOutlined style={{ color: '#fff' }} />
              </div>
              <div className="tf-brand-text">
                Triage<span>Flow</span>
              </div>
            </Link>
            <Menu
              theme="dark"
              mode="horizontal"
              selectedKeys={[getSelectedKey(location.pathname)]}
              items={menuItems}
              style={{ flex: 1, background: 'transparent' }}
            />
            <Link to="/">
              <Button className="tf-lang-btn" type="text" icon={<HomeOutlined />}>
                {t('home')}
              </Button>
            </Link>
            <Button
              className="tf-lang-btn"
              type="text"
              icon={<TranslationOutlined />}
              onClick={toggleLang}
            >
              {lang === 'zh' ? 'EN' : '中文'}
            </Button>
          </Header>
        )}
        <Content className={isPortalHome ? '' : 'tf-content'}>
          <Routes>
            <Route path="/" element={<PortalHome />} />
            <Route path="/console" element={<ConsolePage />} />
            <Route path="/console/new" element={<TaskForm />} />
            <Route path="/console/tasks/:id" element={<TriageDetail />} />
            <Route path="/nurse" element={<NursePortal />} />
            <Route path="/patient" element={<PatientRegister />} />
            <Route path="/patient/queue/:id" element={<PatientQueue />} />
          </Routes>
        </Content>
      </Layout>
    </ConfigProvider>
  )
}
