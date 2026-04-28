import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { useAuth } from '../../context/useAuth'
import { useToast } from '../../context/ToastContext'
import { apiGetProfile, apiUpdateProfile, apiDeleteProfile, apiGetMyInterests } from '../../api'

const Profile = () => {
  const navigate = useNavigate()
  const { profile, setProfile, setToken } = useAuth()
  const toast = useToast()
  const [loading, setLoading] = useState(!profile)
  const [loadError, setLoadError] = useState('')
  const [editing, setEditing] = useState(false)
  const [form, setForm] = useState({ firstName: '', lastName: '', bio: '', avatarUrl: '' })
  const [saving, setSaving] = useState(false)
  const [deleteConfirm, setDeleteConfirm] = useState(false)
  const [avatarError, setAvatarError] = useState(false)
  const [interests, setInterests] = useState([])

  const load = async () => {
    setLoadError('')
    setLoading(true)
    try {
      const [data, interestsData] = await Promise.all([apiGetProfile(), apiGetMyInterests()])
      setProfile(data)
      setInterests(interestsData.items ?? [])
      setForm({
        firstName: data.firstName || '',
        lastName: data.lastName || '',
        bio: data.bio || '',
        avatarUrl: data.avatarUrl || '',
      })
    } catch (e) {
      setLoadError(e.error || 'Failed to load profile')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    if (profile) {
      setForm({
        firstName: profile.firstName || '',
        lastName: profile.lastName || '',
        bio: profile.bio || '',
        avatarUrl: profile.avatarUrl || '',
      })
      apiGetMyInterests()
        .then((interestsData) => setInterests(interestsData.items ?? []))
        .catch(() => toast.error('Failed to load interests'))
    } else {
      load()
    }
  }, [])

  useEffect(() => {
    setAvatarError(false)
  }, [profile?.avatarUrl])

  const handleSave = async (e) => {
    e.preventDefault()
    setSaving(true)
    try {
      const body = {}
      if (form.firstName !== (profile?.firstName ?? '')) body.firstName = form.firstName
      if (form.lastName !== (profile?.lastName ?? '')) body.lastName = form.lastName
      if (form.bio !== (profile?.bio ?? '')) body.bio = form.bio
      if (form.avatarUrl !== (profile?.avatarUrl ?? '')) body.avatarUrl = form.avatarUrl
      const data = await apiUpdateProfile(body)
      setProfile(data)
      setEditing(false)
      toast.success('Profile saved')
    } catch (err) {
      toast.error(err.error || 'Failed to save')
    } finally {
      setSaving(false)
    }
  }

  const handleDeleteAccount = async () => {
    if (!deleteConfirm) {
      setDeleteConfirm(true)
      return
    }
    setSaving(true)
    try {
      await apiDeleteProfile()
      setProfile(null)
      setToken(null)
      navigate('/login')
    } catch (err) {
      toast.error(err.error || 'Failed to delete account')
    } finally {
      setSaving(false)
    }
  }

  const initial = (profile?.firstName?.[0] || profile?.email?.[0] || '?').toUpperCase()
  const showAvatarImage = profile?.avatarUrl && !avatarError

  if (loading) {
    return (
      <div className="page-content">
        <div className="profile-loading">Loading profile...</div>
      </div>
    )
  }

  if (loadError && !profile) {
    return (
      <div className="page-content">
        <div className="profile-card">
          <div className="auth-error">{loadError}</div>
          <button onClick={load} className="btn btn-primary">Try again</button>
        </div>
      </div>
    )
  }

  if (!profile) {
    return (
      <div className="page-content">
        <div className="profile-card">You are not signed in.</div>
      </div>
    )
  }

  return (
    <div className="page-content">
      <header className="page-header">
        <h1 className="page-title">Your profile</h1>
        <p className="page-subtitle">Manage your study profile and preferences</p>
      </header>

      <section className="profile-card profile-summary">
        <div className="profile-summary-left">
          {showAvatarImage ? (
            <img
              src={profile.avatarUrl}
              alt=""
              className="profile-summary-avatar"
              onError={() => setAvatarError(true)}
              onLoad={() => setAvatarError(false)}
            />
          ) : (
            <div className="profile-summary-avatar-initial">{initial}</div>
          )}
          <div>
            <h2 className="profile-summary-name">{profile.firstName} {profile.lastName}</h2>
            <p className="profile-summary-email">{profile.email}</p>
            {!profile.avatarUrl && (
              <p className="profile-avatar-hint">Add a photo via URL. Click Edit and paste an image link.</p>
            )}
          </div>
        </div>
      </section>

      <section className="profile-card">
        <h3 className="profile-card-title">Basic information</h3>
        {editing ? (
          <form onSubmit={handleSave} className="profile-form">
            <label className="profile-label">First name</label>
            <input
              className="profile-input"
              value={form.firstName}
              onChange={(e) => setForm((f) => ({ ...f, firstName: e.target.value }))}
              placeholder="First name"
              required
            />
            <label className="profile-label">Last name</label>
            <input
              className="profile-input"
              value={form.lastName}
              onChange={(e) => setForm((f) => ({ ...f, lastName: e.target.value }))}
              placeholder="Last name"
              required
            />
            <label className="profile-label">Bio</label>
            <textarea
              className="profile-input profile-textarea"
              value={form.bio}
              onChange={(e) => setForm((f) => ({ ...f, bio: e.target.value }))}
              placeholder="Tell us about yourself and your study goals"
              rows={4}
            />
            <label className="profile-label">Profile photo URL</label>
            <input
              className="profile-input"
              type="url"
              value={form.avatarUrl}
              onChange={(e) => setForm((f) => ({ ...f, avatarUrl: e.target.value }))}
              placeholder="https://example.com/photo.jpg"
            />
            <p className="profile-field-hint">Use a direct link to an image (URL that opens the image only).</p>
            <div className="profile-form-actions">
              <button type="submit" className="btn btn-primary" disabled={saving}>
                {saving ? 'Saving...' : 'Save'}
              </button>
              <button type="button" className="btn btn-secondary" onClick={() => setEditing(false)}>
                Cancel
              </button>
            </div>
          </form>
        ) : (
          <>
            <div className="profile-info-row">
              <span className="profile-label">First name</span>
              <span className="profile-value">{profile.firstName}</span>
            </div>
            <div className="profile-info-row">
              <span className="profile-label">Last name</span>
              <span className="profile-value">{profile.lastName}</span>
            </div>
            <div className="profile-info-row">
              <span className="profile-label">Email</span>
              <span className="profile-value">{profile.email}</span>
            </div>
            {profile.bio && (
              <div className="profile-info-row profile-info-bio">
                <span className="profile-label">Bio</span>
                <span className="profile-value">{profile.bio}</span>
              </div>
            )}
            <button type="button" className="btn btn-primary btn-sm" onClick={() => setEditing(true)}>
              Edit
            </button>
          </>
        )}
      </section>

      {interests.length > 0 && (
        <section className="profile-card">
          <h3 className="profile-card-title">Interests</h3>
          <div className="interests-grid" style={{marginTop: '8px'}}>
            {interests.map((interest) => (
              <span key={interest.ID} className="interest-chip selected" style={{cursor: 'default', pointerEvents: 'none'}}>
                {interest.Name}
              </span>
            ))}
          </div>
        </section>
      )}

      <section className="profile-card profile-danger-card">
        <h3 className="profile-card-title">Delete account</h3>
        <p className="profile-danger-text">Once deleted, your data cannot be recovered.</p>
        <button
          type="button"
          className="btn btn-danger"
          onClick={handleDeleteAccount}
          disabled={saving}
        >
          {deleteConfirm ? 'Click again to confirm' : 'Delete account'}
        </button>
        {deleteConfirm && (
          <button type="button" className="btn btn-ghost" onClick={() => setDeleteConfirm(false)}>
            Cancel
          </button>
        )}
      </section>
    </div>
  )
}

export default Profile
