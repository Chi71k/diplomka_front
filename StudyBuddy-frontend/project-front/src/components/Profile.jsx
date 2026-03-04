import React, { useEffect, useState } from 'react'

const API_BASE = import.meta.env.VITE_API_BASE || ''

// simple profile page; expects httpOnly auth cookie so fetch includes credentials
const Profile = ({ profile: initialProfile, onLogout }) => {
  const [profile, setProfile] = useState(initialProfile || null)
  const [loading, setLoading] = useState(!initialProfile)
  const [error, setError] = useState('')

  const load = async () => {
    setError('')
    setLoading(true)
    try {
      const res = await fetch(`${API_BASE}/api/v1/users/me`, {
        method: 'GET',
        credentials: 'include',
      })
      if (!res.ok) {
        if (res.status === 401 || res.status === 404) {
          // not logged in or no profile
          setProfile(null)
        } else {
          const data = await res.json().catch(() => ({}))
          setError(data.error || 'Ошибка при загрузке профиля')
        }
        return
      }
      const data = await res.json()
      setProfile(data)
    } catch (e) {
      console.error(e)
      setError('Сетевая ошибка')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    if (!initialProfile) {
      load()
    }
  }, [])

  const handleLogout = () => {
    // there is no logout endpoint, so just drop the profile state
    setProfile(null)
    if (onLogout) onLogout()
  }

  if (loading) {
    return <div className="auth-card">Загрузка профиля...</div>
  }

  if (error) {
    return (
      <div className="auth-card">
        <div style={{ color: 'crimson', textAlign: 'left' }}>{error}</div>
        <button onClick={load} className="auth-button">
          Попробовать снова
        </button>
      </div>
    )
  }

  if (!profile) {
    return (
      <div className="auth-card">
        <div>Вы не авторизованы</div>
      </div>
    )
  }

  return (
    <div className="auth-card">
      <h2 className="auth-title">Профиль</h2>
      <div className="profile-field">
        <strong>Имя:</strong> {profile.firstName}
      </div>
      <div className="profile-field">
        <strong>Фамилия:</strong> {profile.lastName}
      </div>
      <div className="profile-field">
        <strong>Email:</strong> {profile.email}
      </div>
      {profile.bio && (
        <div className="profile-field">
          <strong>О себе:</strong> {profile.bio}
        </div>
      )}
      {profile.avatarUrl && (
        <div className="profile-field">
          <img
            src={profile.avatarUrl}
            alt="avatar"
            style={{ maxWidth: 150, borderRadius: '50%' }}
          />
        </div>
      )}
      <button onClick={handleLogout} className="auth-button">
        Выйти
      </button>
    </div>
  )
}

export default Profile
