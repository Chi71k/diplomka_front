import { useState, useEffect } from 'react'
import { useAuth } from '../../context/useAuth'
import { useToast } from '../../context/ToastContext'
import { apiGetMatchRequests, apiRespondMatchRequest, apiCancelMatchRequest, apiGetUserById } from '../../api'

const STATUS_LABELS = {
  pending: 'Pending',
  accepted: 'Accepted',
  declined: 'Declined',
  canceled: 'Canceled',
}
const STATUS_STYLES = {
  pending:  { background: '#fffbeb', color: '#d97706', border: '1px solid #fcd34d' },
  accepted: { background: '#f0fdf4', color: '#15803d', border: '1px solid #86efac' },
  declined: { background: '#fef2f2', color: '#dc2626', border: '1px solid #fca5a5' },
  canceled: { background: '#f8fafc', color: '#94a3b8', border: '1px solid #e2e8f0' },
}

const formatDate = (iso) =>
  new Date(iso).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })

const Requests = () => {
  const { profile } = useAuth()
  const toast = useToast()
  const [tab, setTab] = useState('incoming')
  const [requests, setRequests] = useState([])
  const [loading, setLoading] = useState(true)
  const [acting, setActing] = useState(null)
  const [userNames, setUserNames] = useState({})

  useEffect(() => {
    const load = async () => {
      setLoading(true)
      try {
        const data = await apiGetMatchRequests({ limit: 50 })
        const items = data.items ?? []
        setRequests(items)

        const myId = profile?.id
        const ids = [...new Set(
          items.map((r) => r.requesterId === myId ? r.receiverId : r.requesterId)
        )]
        const results = await Promise.allSettled(ids.map((id) => apiGetUserById(id)))
        const names = {}
        ids.forEach((id, i) => {
          if (results[i].status === 'fulfilled' && results[i].value) {
            const p = results[i].value
            names[id] = `${p.firstName} ${p.lastName}`.trim()
          }
        })
        setUserNames(names)
      } catch (e) {
        toast.error(e.error || 'Failed to load requests')
      } finally {
        setLoading(false)
      }
    }
    load()
  }, [])

  const incoming = requests.filter((r) => r.receiverId === profile?.id)
  const outgoing = requests.filter((r) => r.requesterId === profile?.id)
  const shown = tab === 'incoming' ? incoming : outgoing
  const pendingIncoming = incoming.filter((r) => r.status === 'pending').length

  const handleRespond = async (id, accept) => {
    setActing(id)
    try {
      const updated = await apiRespondMatchRequest(id, accept)
      setRequests((prev) => prev.map((r) => (r.id === id ? updated : r)))
      toast.success(accept ? 'Request accepted' : 'Request declined')
    } catch (e) {
      toast.error(e.error || 'Failed to respond')
    } finally {
      setActing(null)
    }
  }

  const handleCancel = async (id) => {
    setActing(id)
    try {
      await apiCancelMatchRequest(id)
      setRequests((prev) => prev.map((r) => (r.id === id ? { ...r, status: 'canceled' } : r)))
      toast.success('Request canceled')
    } catch (e) {
      toast.error(e.error || 'Failed to cancel')
    } finally {
      setActing(null)
    }
  }

  const otherUserId = (r) =>
    r.requesterId === profile?.id ? r.receiverId : r.requesterId

  const displayName = (id) => userNames[id] || id.slice(0, 8) + '...'

  return (
    <div className="page-content">
      <header className="page-header">
        <h1 className="page-title">Match requests</h1>
        <p className="page-subtitle">Manage incoming and outgoing study partner requests</p>
      </header>

      <div style={{ display: 'flex', gap: '8px', marginBottom: '20px' }}>
        <button
          type="button"
          className={`btn btn-sm ${tab === 'incoming' ? 'btn-primary' : 'btn-secondary'}`}
          onClick={() => setTab('incoming')}
        >
          Incoming {pendingIncoming > 0 ? `(${pendingIncoming})` : ''}
        </button>
        <button
          type="button"
          className={`btn btn-sm ${tab === 'outgoing' ? 'btn-primary' : 'btn-secondary'}`}
          onClick={() => setTab('outgoing')}
        >
          Outgoing
        </button>
      </div>

      {loading && <div className="profile-loading">Loading...</div>}

      {!loading && shown.length === 0 && (
        <section className="profile-card">
          <p className="page-muted">
            {tab === 'incoming' ? 'No incoming requests.' : 'No outgoing requests.'}
          </p>
        </section>
      )}

      {!loading && shown.map((r) => {
        const otherId = otherUserId(r)
        return (
          <section key={r.id} className="profile-card" style={{ marginBottom: '12px' }}>
            <div style={{
              display: 'flex', alignItems: 'flex-start',
              justifyContent: 'space-between', gap: '16px', flexWrap: 'wrap',
            }}>
              <div style={{ flex: 1 }}>
                <div style={{ display: 'flex', alignItems: 'center', gap: '10px', flexWrap: 'wrap' }}>
                  <span style={{ fontSize: '0.85rem', color: '#475569' }}>
                    {tab === 'incoming' ? 'From: ' : 'To: '}
                    <strong style={{ color: '#1e293b' }}>{displayName(otherId)}</strong>
                  </span>
                  <span style={{
                    fontSize: '0.78rem', fontWeight: 600, padding: '2px 8px', borderRadius: '6px',
                    ...(STATUS_STYLES[r.status] || STATUS_STYLES.pending),
                  }}>
                    {STATUS_LABELS[r.status] || r.status}
                  </span>
                  <span style={{ fontSize: '0.8rem', color: '#94a3b8' }}>{formatDate(r.createdAt)}</span>
                </div>
                {r.message && (
                  <p style={{ margin: '8px 0 0', fontSize: '0.9rem', color: '#334155' }}>
                    "{r.message}"
                  </p>
                )}
              </div>

              <div style={{ display: 'flex', gap: '8px', flexShrink: 0 }}>
                {tab === 'incoming' && r.status === 'pending' && (
                  <>
                    <button
                      type="button"
                      className="btn btn-primary btn-sm"
                      disabled={acting === r.id}
                      onClick={() => handleRespond(r.id, true)}
                    >
                      Accept
                    </button>
                    <button
                      type="button"
                      className="btn btn-danger btn-sm"
                      disabled={acting === r.id}
                      onClick={() => handleRespond(r.id, false)}
                    >
                      Decline
                    </button>
                  </>
                )}
                {tab === 'outgoing' && r.status === 'pending' && (
                  <button
                    type="button"
                    className="btn btn-secondary btn-sm"
                    disabled={acting === r.id}
                    onClick={() => handleCancel(r.id)}
                  >
                    {acting === r.id ? 'Canceling...' : 'Cancel'}
                  </button>
                )}
              </div>
            </div>
          </section>
        )
      })}
    </div>
  )
}

export default Requests
