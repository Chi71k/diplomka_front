import { useState, useEffect } from 'react'
import { useToast } from '../../context/ToastContext'
import { apiGetMatchRequests, apiRespondMatchRequest, apiCancelMatchRequest, apiGetUserById } from '../../api'
import { useAuth } from '../../context/useAuth'

const statusLabel = {
  pending: 'Pending',
  accepted: 'Accepted',
  declined: 'Declined',
  canceled: 'Canceled',
}

const statusColor = {
  pending: '#d97706',
  accepted: '#15803d',
  declined: '#dc2626',
  canceled: '#94a3b8',
}

const Requests = () => {
  const toast = useToast()
  const { profile } = useAuth()
  const [tab, setTab] = useState('incoming')
  const [requests, setRequests] = useState([])
  const [userCache, setUserCache] = useState({})
  const [loading, setLoading] = useState(true)
  const [acting, setActing] = useState(null)

  const load = async () => {
    setLoading(true)
    try {
      const data = await apiGetMatchRequests({ limit: 50 })
      const items = data.items ?? []
      setRequests(items)

      const ids = [...new Set(items.flatMap(r => [r.requesterId, r.receiverId]))]
      const entries = await Promise.allSettled(ids.map(id => apiGetUserById(id).then(u => [id, u])))
      const cache = {}
      for (const r of entries) {
        if (r.status === 'fulfilled') {
          const [id, user] = r.value
          cache[id] = user
        }
      }
      setUserCache(cache)
    } catch (e) {
      toast.error(e.error || 'Failed to load requests')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => { load() }, [])

  const myId = profile?.id

  const incoming = requests.filter(r => r.receiverId === myId)
  const outgoing = requests.filter(r => r.requesterId === myId)
  const shown = tab === 'incoming' ? incoming : outgoing

  const handleRespond = async (id, accept) => {
    setActing(id)
    try {
      await apiRespondMatchRequest(id, accept)
      toast.success(accept ? 'Request accepted!' : 'Request declined')
      setRequests(prev => prev.map(r => r.id === id ? { ...r, status: accept ? 'accepted' : 'declined' } : r))
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
      toast.success('Request canceled')
      setRequests(prev => prev.map(r => r.id === id ? { ...r, status: 'canceled' } : r))
    } catch (e) {
      toast.error(e.error || 'Failed to cancel')
    } finally {
      setActing(null)
    }
  }

  const formatDate = (iso) =>
    new Date(iso).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })

  return (
    <div className="page-content">
      <header className="page-header">
        <h1 className="page-title">Match Requests</h1>
        <p className="page-subtitle">Manage your incoming and outgoing study partner requests</p>
      </header>

      <div style={{ display: 'flex', gap: '8px', marginBottom: '16px' }}>
        {['incoming', 'outgoing'].map(t => (
          <button
            key={t}
            type="button"
            className={`btn ${tab === t ? 'btn-primary' : 'btn-secondary'} btn-sm`}
            onClick={() => setTab(t)}
          >
            {t === 'incoming' ? `Incoming (${incoming.length})` : `Outgoing (${outgoing.length})`}
          </button>
        ))}
      </div>

      {loading && <div className="profile-loading">Loading requests...</div>}

      {!loading && shown.length === 0 && (
        <section className="profile-card">
          <p className="page-muted">No {tab} requests yet.</p>
        </section>
      )}

      {!loading && shown.map(r => {
        const otherId = tab === 'incoming' ? r.requesterId : r.receiverId
        const other = userCache[otherId]
        return (
          <section key={r.id} className="profile-card" style={{ marginBottom: '12px' }}>
            <div style={{ display: 'flex', alignItems: 'flex-start', gap: '12px', flexWrap: 'wrap' }}>
              <div style={{
                width: '40px', height: '40px', borderRadius: '50%', flexShrink: 0,
                background: 'linear-gradient(135deg, #60a5fa 0%, #3b82f6 100%)',
                display: 'flex', alignItems: 'center', justifyContent: 'center',
                color: '#fff', fontWeight: 600, fontSize: '1rem', overflow: 'hidden',
              }}>
                {other?.avatarUrl
                  ? <img src={other.avatarUrl} alt="" style={{ width: '100%', height: '100%', objectFit: 'cover' }} />
                  : (other?.firstName?.[0] || '?').toUpperCase()
                }
              </div>

              <div style={{ flex: 1, minWidth: 0 }}>
                <div style={{ display: 'flex', alignItems: 'center', gap: '8px', flexWrap: 'wrap' }}>
                  <strong style={{ color: '#1e293b' }}>
                    {other ? `${other.firstName} ${other.lastName}` : otherId}
                  </strong>
                  <span style={{
                    fontSize: '0.75rem', fontWeight: 600, padding: '2px 7px', borderRadius: '5px',
                    background: '#f8fafc', border: '1px solid #e2e8f0',
                    color: statusColor[r.status] || '#94a3b8',
                  }}>
                    {statusLabel[r.status] || r.status}
                  </span>
                </div>
                <p style={{ margin: '2px 0 0', fontSize: '0.8rem', color: '#94a3b8' }}>
                  {formatDate(r.createdAt)}
                </p>
                {r.message && (
                  <p style={{ margin: '6px 0 0', fontSize: '0.875rem', color: '#475569', fontStyle: 'italic' }}>
                    "{r.message}"
                  </p>
                )}
              </div>

              {r.status === 'pending' && tab === 'incoming' && (
                <div style={{ display: 'flex', gap: '8px' }}>
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
                </div>
              )}

              {r.status === 'pending' && tab === 'outgoing' && (
                <button
                  type="button"
                  className="btn btn-secondary btn-sm"
                  disabled={acting === r.id}
                  onClick={() => handleCancel(r.id)}
                >
                  Cancel
                </button>
              )}
            </div>
          </section>
        )
      })}
    </div>
  )
}

export default Requests
