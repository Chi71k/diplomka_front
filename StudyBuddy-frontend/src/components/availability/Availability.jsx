import { useState, useEffect, useRef } from 'react'
import { useToast } from '../../context/ToastContext'
import {
  apiGetSlots, apiCreateSlot, apiDeleteSlot,
  apiGetGCalConnectUrl, apiImportGCal,
} from '../../api'

const DAYS = ['Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday', 'Sunday']

const Availability = () => {
  const toast = useToast()
  const [slots, setSlots] = useState([])
  const [loading, setLoading] = useState(true)
  const [adding, setAdding] = useState(false)
  const [gcalConnecting, setGcalConnecting] = useState(false)
  const [gcalImporting, setGcalImporting] = useState(false)
  const gcalMsgHandlerRef = useRef(null)
  const [form, setForm] = useState({
    dayOfWeek: 1,
    startTime: '09:00',
    endTime: '10:00',
    timezone: Intl.DateTimeFormat().resolvedOptions().timeZone,
  })

  const load = async () => {
    setLoading(true)
    try {
      const data = await apiGetSlots()
      setSlots(data.items ?? [])
    } catch (e) {
      toast.error(e.error || 'Failed to load slots')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    load()
    return () => {
      if (gcalMsgHandlerRef.current) {
        window.removeEventListener('message', gcalMsgHandlerRef.current)
      }
    }
  }, [])

  const handleAdd = async (e) => {
    e.preventDefault()
    setAdding(true)
    try {
      const slot = await apiCreateSlot({
        dayOfWeek: Number(form.dayOfWeek),
        startTime: form.startTime,
        endTime: form.endTime,
        timezone: form.timezone,
      })
      setSlots((prev) => [...prev, slot])
      setForm((f) => ({ ...f, startTime: '09:00', endTime: '10:00' }))
      toast.success('Slot added')
    } catch (e) {
      toast.error(e.error || 'Failed to add slot')
    } finally {
      setAdding(false)
    }
  }

  const handleDelete = async (id) => {
    try {
      await apiDeleteSlot(id)
      setSlots((prev) => prev.filter((s) => s.id !== id))
      toast.success('Slot removed')
    } catch (e) {
      toast.error(e.error || 'Failed to remove slot')
    }
  }

  const handleGCalConnect = async () => {
    if (gcalMsgHandlerRef.current) {
      window.removeEventListener('message', gcalMsgHandlerRef.current)
      gcalMsgHandlerRef.current = null
    }
    setGcalConnecting(true)
    try {
      const data = await apiGetGCalConnectUrl()
      window.open(data.authUrl, 'gcal-oauth', 'width=500,height=600')
      const handleMsg = (e) => {
        if (e.data?.type === 'GCAL_CONNECTED') {
          toast.success('Google Calendar connected')
          window.removeEventListener('message', handleMsg)
          gcalMsgHandlerRef.current = null
          load()
        }
      }
      gcalMsgHandlerRef.current = handleMsg
      window.addEventListener('message', handleMsg)
    } catch (e) {
      toast.error(e.error || 'Failed to connect Google Calendar')
    } finally {
      setGcalConnecting(false)
    }
  }

  const handleGCalImport = async () => {
    setGcalImporting(true)
    try {
      const data = await apiImportGCal()
      toast.success(`Imported ${data.imported} slot${data.imported !== 1 ? 's' : ''} from Google Calendar`)
      load()
    } catch (e) {
      toast.error(e.error || 'Failed to import from Google Calendar')
    } finally {
      setGcalImporting(false)
    }
  }

  const sortedSlots = [...slots].sort((a, b) => {
    if (a.dayOfWeek !== b.dayOfWeek) return a.dayOfWeek - b.dayOfWeek
    return a.startTime.localeCompare(b.startTime)
  })

  return (
    <div className="page-content">
      <header className="page-header">
        <h1 className="page-title">Availability</h1>
        <p className="page-subtitle">Set your weekly study schedule</p>
      </header>

      <section className="profile-card">
        <h3 className="profile-card-title">Add time slot</h3>
        <form onSubmit={handleAdd} className="profile-form">
          <div style={{ display: 'flex', gap: '12px', flexWrap: 'wrap' }}>
            <div style={{ flex: '1', minWidth: '140px' }}>
              <label className="profile-label">Day</label>
              <select
                className="profile-input"
                value={form.dayOfWeek}
                onChange={(e) => setForm((f) => ({ ...f, dayOfWeek: Number(e.target.value) }))}
              >
                {DAYS.map((name, i) => (
                  <option key={i} value={i}>{name}</option>
                ))}
              </select>
            </div>
            <div style={{ flex: '1', minWidth: '120px' }}>
              <label className="profile-label">Start</label>
              <input
                className="profile-input"
                type="time"
                value={form.startTime}
                onChange={(e) => setForm((f) => ({ ...f, startTime: e.target.value }))}
                required
              />
            </div>
            <div style={{ flex: '1', minWidth: '120px' }}>
              <label className="profile-label">End</label>
              <input
                className="profile-input"
                type="time"
                value={form.endTime}
                onChange={(e) => setForm((f) => ({ ...f, endTime: e.target.value }))}
                required
              />
            </div>
            <div style={{ flex: '2', minWidth: '180px' }}>
              <label className="profile-label">Timezone</label>
              <input
                className="profile-input"
                value={form.timezone}
                onChange={(e) => setForm((f) => ({ ...f, timezone: e.target.value }))}
                placeholder="Europe/Kyiv"
                required
              />
            </div>
          </div>
          <div className="profile-form-actions">
            <button type="submit" className="btn btn-primary" disabled={adding}>
              {adding ? 'Adding...' : 'Add slot'}
            </button>
          </div>
        </form>
      </section>

      <section className="profile-card">
        <h3 className="profile-card-title">Your slots</h3>
        {loading && <div className="profile-loading">Loading...</div>}
        {!loading && sortedSlots.length === 0 && (
          <p className="page-muted">No slots yet. Add your availability above.</p>
        )}
        {!loading && sortedSlots.length > 0 && (
          <ul className="course-list">
            {sortedSlots.map((s) => (
              <li key={s.id} className="course-list-item">
                <div
                  className="course-list-link"
                  style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}
                >
                  <span>
                    <strong>{DAYS[s.dayOfWeek]}</strong>
                    <span className="course-meta">{s.startTime} – {s.endTime} · {s.timezone}</span>
                  </span>
                  <button
                    type="button"
                    className="btn btn-danger btn-sm"
                    onClick={() => handleDelete(s.id)}
                  >
                    Remove
                  </button>
                </div>
              </li>
            ))}
          </ul>
        )}
      </section>

      <section className="profile-card">
        <h3 className="profile-card-title">Google Calendar</h3>
        <p style={{ margin: '0 0 16px', fontSize: '0.9rem', color: '#64748b' }}>
          Connect Google Calendar to automatically import your availability slots.
        </p>
        <div style={{ display: 'flex', gap: '12px', flexWrap: 'wrap' }}>
          <button
            type="button"
            className="btn btn-primary"
            onClick={handleGCalConnect}
            disabled={gcalConnecting}
          >
            {gcalConnecting ? 'Connecting...' : 'Connect Google Calendar'}
          </button>
          <button
            type="button"
            className="btn btn-secondary"
            onClick={handleGCalImport}
            disabled={gcalImporting}
          >
            {gcalImporting ? 'Importing...' : 'Import from Google Calendar'}
          </button>
        </div>
      </section>
    </div>
  )
}

export default Availability
