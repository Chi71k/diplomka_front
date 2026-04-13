import { createContext, useState, useCallback } from 'react'
import { getToken } from '../api'

const AuthContext = createContext(null)

export function AuthProvider({ children }) {
  const [token, setTokenState] = useState(() => getToken())
  const [profile, setProfile] = useState(null)

  const setToken = useCallback((newToken) => {
    if (newToken) {
      localStorage.setItem('accessToken', newToken)
      setTokenState(newToken)
    } else {
      localStorage.removeItem('accessToken')
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

export { AuthContext }
