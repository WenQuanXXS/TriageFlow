import { createContext, useContext, useState } from 'react'

const messages = {
  en: {
    dashboard: 'Dashboard',
    newTask: 'New Task',
    newTriageTask: 'New Triage Task',
    patientName: 'Patient Name',
    chiefComplaint: 'Chief Complaint',
    priority: 'Priority',
    department: 'Department',
    submit: 'Submit',
    enterPatientName: 'Enter patient name',
    enterChiefComplaint: 'Describe the chief complaint',
    selectDepartment: 'Select department',
    pleaseEnterPatientName: 'Please enter patient name',
    pleaseEnterChiefComplaint: 'Please enter chief complaint',
    taskCreatedSuccess: 'Task created successfully',
    taskCreatedFail: 'Failed to create task',
    loadTasksFail: 'Failed to load tasks',
    updateStatusFail: 'Failed to update status',
    loadDashboardFail: 'Failed to load dashboard',
    toggleStatus: 'Toggle Status',
    patient: 'Patient',
    status: 'Status',
    created: 'Created',
    action: 'Action',
    total: 'Total',
    allStatus: 'All Status',
    allPriority: 'All Priority',
    pending: 'Pending',
    inProgress: 'In Progress',
    completed: 'Completed',
    urgent: 'Urgent',
    high: 'High',
    normal: 'Normal',
    low: 'Low',
    internalMedicine: 'Internal Medicine',
    surgery: 'Surgery',
    pediatrics: 'Pediatrics',
    neurology: 'Neurology',
    emergency: 'Emergency',
    statusPending: 'PENDING',
    statusInProgress: 'IN PROGRESS',
    statusCompleted: 'COMPLETED',
    priorityUrgent: 'URGENT',
    priorityHigh: 'HIGH',
    priorityNormal: 'NORMAL',
    priorityLow: 'LOW',
  },
  zh: {
    dashboard: '仪表盘',
    newTask: '新建任务',
    newTriageTask: '新建分诊任务',
    patientName: '患者姓名',
    chiefComplaint: '主诉',
    priority: '优先级',
    department: '科室',
    submit: '提交',
    enterPatientName: '请输入患者姓名',
    enterChiefComplaint: '请描述主诉症状',
    selectDepartment: '请选择科室',
    pleaseEnterPatientName: '请输入患者姓名',
    pleaseEnterChiefComplaint: '请输入主诉',
    taskCreatedSuccess: '任务创建成功',
    taskCreatedFail: '任务创建失败',
    loadTasksFail: '加载任务失败',
    updateStatusFail: '更新状态失败',
    loadDashboardFail: '加载仪表盘失败',
    toggleStatus: '切换状态',
    patient: '患者',
    status: '状态',
    created: '创建时间',
    action: '操作',
    total: '总计',
    allStatus: '全部状态',
    allPriority: '全部优先级',
    pending: '待处理',
    inProgress: '进行中',
    completed: '已完成',
    urgent: '紧急',
    high: '高',
    normal: '普通',
    low: '低',
    internalMedicine: '内科',
    surgery: '外科',
    pediatrics: '儿科',
    neurology: '神经科',
    emergency: '急诊',
    statusPending: '待处理',
    statusInProgress: '进行中',
    statusCompleted: '已完成',
    priorityUrgent: '紧急',
    priorityHigh: '高',
    priorityNormal: '普通',
    priorityLow: '低',
  },
}

const LocaleContext = createContext()

export function LocaleProvider({ children }) {
  const [lang, setLang] = useState('zh')

  const t = (key) => messages[lang]?.[key] ?? key
  const toggleLang = () => setLang((l) => (l === 'zh' ? 'en' : 'zh'))

  return (
    <LocaleContext.Provider value={{ lang, t, toggleLang }}>
      {children}
    </LocaleContext.Provider>
  )
}

export function useLocale() {
  return useContext(LocaleContext)
}
