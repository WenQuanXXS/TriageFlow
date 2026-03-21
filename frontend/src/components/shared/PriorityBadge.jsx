const priorityLabelKeys = {
  urgent: 'priorityUrgent',
  high: 'priorityHigh',
  normal: 'priorityNormal',
  low: 'priorityLow',
}

export { priorityLabelKeys }

export default function PriorityBadge({ priority, label }) {
  return (
    <span className={`tf-priority-badge ${priority}`}>
      {label}
    </span>
  )
}
