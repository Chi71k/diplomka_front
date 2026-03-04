import React, { useState } from 'react'

const API_BASE = import.meta.env.VITE_API_BASE || ''

const Register = ({ onSwitch, onSuccess }) => {
  const [name, setName] = useState('')
  const [surname, setSurname] = useState('')
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  const handleSubmit = async (e) => {
    e.preventDefault()
    setError('')
    if (!name || !surname || !email || !password) {
      setError('Все поля обязательны')
      return
    }
    setLoading(true)
    try {
      const body = {
        email,
        password,
        firstName: name,
        lastName: surname,
      }
      const res = await fetch(`${API_BASE}/api/v1/auth/register`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(body),
        credentials: 'include',
      })
      const data = await res.json()
      if (!res.ok) {
        setError(data.error || 'Ошибка регистрации')
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
      <h2 className="auth-title">Регистрация</h2>
      <form onSubmit={handleSubmit} className="auth-form">
        <input
          className="auth-input"
          type="text"
          placeholder="Имя"
          value={name}
          onChange={(e) => setName(e.target.value)}
          required
          disabled={loading}
        />
        <input
          className="auth-input"
          type="text"
          placeholder="Фамилия"
          value={surname}
          onChange={(e) => setSurname(e.target.value)}
          required
          disabled={loading}
        />
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
          {loading ? 'Загрузка...' : 'Создать аккаунт'}
        </button>
      </form>
      {error && <div style={{ color: 'crimson', textAlign: 'left' }}>{error}</div>}
      <div className="auth-switch">
        Уже есть аккаунт? <button className="linkish" onClick={onSwitch}>Войти</button>
      </div>
    </div>
  )
}

export default Register
