import { Routes, Route, Link, useLocation } from 'react-router-dom'
import { Layout, Menu, Button, ConfigProvider } from 'antd'
import { PlusOutlined, UnorderedListOutlined, TranslationOutlined } from '@ant-design/icons'
import zhCN from 'antd/locale/zh_CN'
import enUS from 'antd/locale/en_US'
import TaskForm from './components/TaskForm'
import TaskList from './components/TaskList'
import Dashboard from './components/Dashboard'
import TriageDetail from './components/TriageDetail'
import { useLocale } from './locales'

const { Header, Content } = Layout

function HomePage() {
  return (
    <>
      <Dashboard />
      <TaskList />
    </>
  )
}

export default function App() {
  const { lang, t, toggleLang } = useLocale()
  const location = useLocation()

  const menuItems = [
    { key: '/', label: <Link to="/">{t('dashboard')}</Link>, icon: <UnorderedListOutlined /> },
    { key: '/new', label: <Link to="/new">{t('newTask')}</Link>, icon: <PlusOutlined /> },
  ]

  return (
    <ConfigProvider locale={lang === 'zh' ? zhCN : enUS}>
      <Layout style={{ minHeight: '100vh' }}>
        <Header style={{ display: 'flex', alignItems: 'center' }}>
          <Menu
            theme="dark"
            mode="horizontal"
            selectedKeys={[location.pathname]}
            items={menuItems}
            style={{ flex: 1 }}
          />
          <Button
            type="text"
            icon={<TranslationOutlined />}
            onClick={toggleLang}
            style={{ color: '#fff' }}
          >
            {lang === 'zh' ? 'English' : '中文'}
          </Button>
        </Header>
        <Content style={{ padding: '24px', maxWidth: 1200, margin: '0 auto', width: '100%' }}>
          <Routes>
            <Route path="/" element={<HomePage />} />
            <Route path="/new" element={<TaskForm />} />
            <Route path="/tasks/:id" element={<TriageDetail />} />
          </Routes>
        </Content>
      </Layout>
    </ConfigProvider>
  )
}
