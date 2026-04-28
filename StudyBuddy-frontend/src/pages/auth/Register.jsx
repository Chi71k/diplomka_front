import React, { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { apiRegister } from '../../api'
import AuthLayout from '../../components/AuthLayout'

const Register = ({ onSuccess }) => {
  const navigate = useNavigate()
  const [firstName, setFirstName] = useState('')
  const [lastName, setLastName] = useState('')
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  const handleSubmit = async (e) => {
    e.preventDefault()
    setError('')
    if (!firstName || !lastName || !email || !password) {
      setError('Fill in all fields')
      return
    }
    if (password.length < 8) {
      setError('Password must be at least 8 characters')
      return
    }
    setLoading(true)
    try {
      const data = await apiRegister({ email, password, firstName, lastName })
      if (onSuccess) onSuccess(data.accessToken)
      navigate('/profile')
    } catch (err) {
      setError(err.error || 'Registration failed')
    } finally {
      setLoading(false)
    }
  }

  return (
    <AuthLayout>
      <div className="auth-form-wrap">
        <h2 className="auth-form-title">Create account</h2>
        <p className="auth-form-subtitle">Sign up to find study partners</p>
        <form onSubmit={handleSubmit} className="auth-form">
          <label className="auth-label">First name</label>
          <input
            className="auth-input"
            type="text"
            placeholder="First name"
            value={firstName}
            onChange={(e) => setFirstName(e.target.value)}
            required
            disabled={loading}
          />
          <label className="auth-label">Last name</label>
          <input
            className="auth-input"
            type="text"
            placeholder="Last name"
            value={lastName}
            onChange={(e) => setLastName(e.target.value)}
            required
            disabled={loading}
          />
          <label className="auth-label">Email</label>
          <input
            className="auth-input"
            type="email"
            placeholder="you@university.edu"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
            disabled={loading}
          />
          <label className="auth-label">Password</label>
          <input
            className="auth-input"
            type="password"
            placeholder="At least 8 characters"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
            minLength={8}
            disabled={loading}
          />
          {error && <div className="auth-error">{error}</div>}
          <button className="auth-button auth-button-primary" type="submit" disabled={loading}>
            {loading ? 'Signing up...' : 'Sign up'}
          </button>
        </form>
        <p className="auth-switch">
          Already have an account? <Link to="/login" className="auth-link">Sign in</Link>
        </p>
      </div>
    </AuthLayout>
  )
}

export default Register
