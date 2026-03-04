import React, { useState } from 'react'

const API_BASE = import.meta.env.VITE_API_BASE || ''

const Login = ({ onSwitch, onSuccess }) => {
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  const handleSubmit = async (e) => {
    e.preventDefault()
    setError('')
    if (!email || !password) {
      setError('Email и пароль обязательны')
      return
    }
    setLoading(true)
    try {
      const res = await fetch(`${API_BASE}/api/v1/auth/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, password }),
        credentials: 'include',
      })
      const data = await res.json()
      if (!res.ok) {
        setError(data.error || 'Ошибка входа')
        return
      }
      // backend returns tokens in body; pass access token to parent
      if (onSuccess) onSuccess(data.accessToken)
    } catch (err) {
      setError('Сетевая ошибка')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="auth-card">
      <h2 className="auth-title">Вход</h2>
      <form onSubmit={handleSubmit} className="auth-form">
        <input
          className="auth-input"
          type="email"
          placeholder="Email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          required
          disabled={loading}
        />
        <input
          className="auth-input"
          type="password"
          placeholder="Пароль"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          required
          disabled={loading}
        />
        <button className="auth-button" type="submit" disabled={loading}>
          {loading ? 'Загрузка...' : 'Войти'}
        </button>
      </form>
      {error && <div style={{ color: 'crimson', textAlign: 'left' }}>{error}</div>}
      <div className="auth-switch">
        Нет аккаунта? <button className="linkish" onClick={onSwitch}>Зарегистрироваться</button>
      </div>
    </div>
  )
}

export default Login
