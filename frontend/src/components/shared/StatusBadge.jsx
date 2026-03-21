const statusColors = {
  pending: '#e67e22',
  in_progress: '#3498db',
  completed: '#27ae60',
}

const statusLabelKeys = {
  pending: 'statusPending',
  in_progress: 'statusInProgress',
  completed: 'statusCompleted',
}

export { statusColors, statusLabelKeys }

export default function StatusBadge({ status, label }) {
  return (
    <span className="tf-status-badge">
      <span className="tf-status-dot" style={{ background: statusColors[status] }} />
      {label || status}
    </span>
  )
}
