import React, { useState, useEffect } from 'react'
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { AuthProvider, useAuth } from './context/AuthContext'
import { apiGetProfile } from './api'
import './App.css'
import Login from './components/Login'
import Register from './components/Register'
import AppLayout from './components/AppLayout'
import Profile from './components/Profile'
import CourseList from './components/courses/CourseList'
import CourseDetail from './components/courses/CourseDetail'
import CourseForm from './components/courses/CourseForm'

function AppRoutes() {
  const { token, setToken, profile, setProfile } = useAuth()
  const [loadingProfile, setLoadingProfile] = useState(true)

  const fetchProfile = async () => {
    setLoadingProfile(true)
    try {
      if (token) {
        const data = await apiGetProfile()
        setProfile(data)
      } else {
        setProfile(null)
      }
    } catch {
      setProfile(null)
    } finally {
      setLoadingProfile(false)
    }
  }

  useEffect(() => {
    if (token) {
      fetchProfile()
    } else {
      setLoadingProfile(false)
      setProfile(null)
    }
  }, [token])

  const handleLoginSuccess = (newToken) => {
    if (newToken) setToken(newToken)
  }

  const handleRegisterSuccess = (newToken) => {
    if (newToken) setToken(newToken)
  }

  const handleLogout = () => {
    setProfile(null)
    setToken(null)
  }

  if (loadingProfile && token) {
    return <div className="app-root app-loading">Loading...</div>
  }

  return (
    <Routes>
      <Route path="/login" element={
        token ? <Navigate to="/profile" replace /> : <Login onSuccess={handleLoginSuccess} />
      } />
      <Route path="/register" element={
        token ? <Navigate to="/profile" replace /> : <Register onSuccess={handleRegisterSuccess} />
      } />
      <Route path="/" element={token ? <AppLayout onLogout={handleLogout} /> : <Navigate to="/login" replace />}>
        <Route index element={<Navigate to="/profile" replace />} />
        <Route path="profile" element={<Profile />} />
        <Route path="courses" element={<CourseList />} />
        <Route path="courses/new" element={<CourseForm />} />
        <Route path="courses/:id" element={<CourseDetail />} />
        <Route path="courses/:id/edit" element={<CourseForm edit />} />
      </Route>
      <Route path="*" element={<Navigate to={token ? '/profile' : '/login'} replace />} />
    </Routes>
  )
}

const App = () => (
  <AuthProvider>
    <BrowserRouter>
      <div className="app-root">
        <AppRoutes />
      </div>
    </BrowserRouter>
  </AuthProvider>
)

export default App
