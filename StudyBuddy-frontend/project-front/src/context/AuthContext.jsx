import React, { createContext, useContext, useState, useCallback } from 'react'
import { getToken } from '../api'

const AuthContext = createContext(null)

export function AuthProvider({ children }) {
  const [token, setTokenState] = useState(() => getToken())
  const [profile, setProfile] = useState(null)

  const setToken = useCallback((newToken) => {
    if (newToken) {
      sessionStorage.setItem('accessToken', newToken)
      setTokenState(newToken)
    } else {
      sessionStorage.removeItem('accessToken')
      setTokenState(null)
      setProfile(null)
    }
  }, [])

  const value = {
    token,
    setToken,
    profile,
    setProfile,
    isAuthenticated: !!token,
  }
  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}

export function useAuth() {
  const ctx = useContext(AuthContext)
  if (!ctx) throw new Error('useAuth must be used within AuthProvider')
  return ctx
}
