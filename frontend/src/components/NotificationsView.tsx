import { useQuery } from "@connectrpc/connect-query"
import { listNotifications } from "../gen/v1/notification-NotificationService_connectquery.ts"

function formatDate(seconds: bigint | undefined): string {
  if (seconds === undefined) return ""
  return new Date(Number(seconds) * 1000).toLocaleString()
}

function channelLabel(channel: number): string {
  switch (channel) {
    case 1:
      return "Email"
    case 2:
      return "SMS"
    case 3:
      return "Push"
    default:
      return "—"
  }
}

export function NotificationsView() {
  const notifQ = useQuery(listNotifications, { pageSize: 50, pageToken: "" }, { refetchInterval: 5000 })

  if (notifQ.isPending) return <p>Loading notifications…</p>
  if (notifQ.isError) return <p className="error">Failed to load notifications: {notifQ.error.message}</p>

  const notifications = notifQ.data?.notifications ?? []
  if (notifications.length === 0) return <p>No notifications yet. Place an order to receive one.</p>

  return (
    <div className="notifications">
      <h2>Notifications</h2>
      <table className="table">
        <thead>
          <tr>
            <th>Status</th>
            <th>Channel</th>
            <th>Subject</th>
            <th>Body</th>
            <th>Created</th>
          </tr>
        </thead>
        <tbody>
          {notifications.map((n) => (
            <tr key={n.id}>
              <td>{n.status}</td>
              <td>{channelLabel(n.channel)}</td>
              <td>{n.subject}</td>
              <td>{n.body}</td>
              <td>{formatDate(n.createdAt?.seconds)}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}
