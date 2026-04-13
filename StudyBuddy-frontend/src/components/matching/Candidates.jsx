import { useState, useEffect } from 'react'
import { useToast } from '../../context/ToastContext'
import { apiGetCandidates, apiSendMatchRequest } from '../../api'

const scoreColor = (score) => {
  if (score >= 0.7) return '#15803d'
  if (score >= 0.4) return '#d97706'
  return '#94a3b8'
}

const Candidates = () => {
  const toast = useToast()
  const [candidates, setCandidates] = useState([])
  const [loading, setLoading] = useState(true)
  const [sending, setSending] = useState(null)
  const [messages, setMessages] = useState({})

  useEffect(() => {
    const load = async () => {
      setLoading(true)
      try {
        const data = await apiGetCandidates(20)
        setCandidates(data.items ?? [])
      } catch (e) {
        toast.error(e.error || 'Failed to load candidates')
      } finally {
        setLoading(false)
      }
    }
    load()
  }, [])

  const handleSend = async (userId) => {
    setSending(userId)
    try {
      await apiSendMatchRequest(userId, messages[userId] || '')
      toast.success('Match request sent!')
      setCandidates((prev) => prev.filter((c) => c.userId !== userId))
    } catch (e) {
      toast.error(e.error || 'Failed to send request')
    } finally {
      setSending(null)
    }
  }

  return (
    <div className="page-content">
      <header className="page-header">
        <h1 className="page-title">Find study partners</h1>
        <p className="page-subtitle">
          Ranked by shared interests (40%), availability (40%), and courses (20%)
        </p>
      </header>

      {loading && <div className="profile-loading">Loading candidates...</div>}

      {!loading && candidates.length === 0 && (
        <section className="profile-card">
          <p className="page-muted">
            No candidates found. Make sure your interests, courses, and availability are filled in.
          </p>
        </section>
      )}

      {!loading && candidates.map((c) => (
        <section key={c.userId} className="profile-card" style={{ marginBottom: '16px' }}>
          <div style={{ display: 'flex', alignItems: 'flex-start', gap: '16px', flexWrap: 'wrap' }}>

            {/* Avatar */}
            <div
              style={{
                width: '48px', height: '48px', borderRadius: '50%', flexShrink: 0,
                background: 'linear-gradient(135deg, #60a5fa 0%, #3b82f6 100%)',
                display: 'flex', alignItems: 'center', justifyContent: 'center',
                color: '#fff', fontWeight: 600, fontSize: '1.1rem', overflow: 'hidden',
              }}
            >
              {c.avatarUrl
                ? <img src={c.avatarUrl} alt="" style={{ width: '100%', height: '100%', objectFit: 'cover' }} />
                : (c.firstName?.[0] || '?').toUpperCase()
              }
            </div>

            {/* Info */}
            <div style={{ flex: 1, minWidth: 0 }}>
              <div style={{ display: 'flex', alignItems: 'center', gap: '10px', flexWrap: 'wrap' }}>
                <strong style={{ fontSize: '1rem', color: '#1e293b' }}>
                  {c.firstName} {c.lastName}
                </strong>
                <span style={{
                  fontSize: '0.8rem', fontWeight: 700, padding: '2px 8px', borderRadius: '6px',
                  background: '#f8fafc', border: '1px solid #e2e8f0',
                  color: scoreColor(c.overallScore),
                }}>
                  {Math.round(c.overallScore * 100)}% match
                </span>
              </div>
              {c.bio && (
                <p style={{ margin: '6px 0 0', fontSize: '0.875rem', color: '#64748b' }}>{c.bio}</p>
              )}
              <div style={{ marginTop: '6px', display: 'flex', gap: '12px', flexWrap: 'wrap', fontSize: '0.8rem', color: '#94a3b8' }}>
                {c.commonCourses?.length > 0 && (
                  <span>{c.commonCourses.length} common course{c.commonCourses.length > 1 ? 's' : ''}</span>
                )}
                {c.commonSlots?.length > 0 && (
                  <span>{c.commonSlots.length} common slot{c.commonSlots.length > 1 ? 's' : ''}</span>
                )}
              </div>
            </div>

            {/* Action */}
            <div style={{ display: 'flex', flexDirection: 'column', gap: '8px', minWidth: '200px' }}>
              <input
                className="profile-input"
                placeholder="Message (optional)"
                value={messages[c.userId] || ''}
                onChange={(e) => setMessages((m) => ({ ...m, [c.userId]: e.target.value }))}
                maxLength={500}
              />
              <button
                type="button"
                className="btn btn-primary btn-sm"
                disabled={sending === c.userId}
                onClick={() => handleSend(c.userId)}
              >
                {sending === c.userId ? 'Sending...' : 'Send request'}
              </button>
            </div>
          </div>
        </section>
      ))}

    </div>
  )
}

export default Candidates
