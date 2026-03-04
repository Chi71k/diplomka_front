import React, { useState, useEffect } from 'react'
import './App.css'
import Login from './components/Login'
import Register from './components/Register'
import Profile from './components/Profile'

const API_BASE = import.meta.env.VITE_API_BASE || ''

const App = () => {
	const [mode, setMode] = useState('login')
	const [profile, setProfile] = useState(null)
	const [loadingProfile, setLoadingProfile] = useState(true)
	const [accessToken, setAccessToken] = useState(() => sessionStorage.getItem('accessToken') || null)

	// try to load profile on mount; if cookie is present backend will return data
	const fetchProfile = async (token) => {
		setLoadingProfile(true)
		try {
			const header = token || accessToken ? { Authorization: `Bearer ${token || accessToken}` } : {}
			const res = await fetch(`${API_BASE}/api/v1/users/me`, {
				method: 'GET',
				credentials: 'include',
				headers: header,
			})
			if (res.ok) {
				const data = await res.json()
				setProfile(data)
			} else {
				setProfile(null)
			}
		} catch (e) {
			console.error(e)
			setProfile(null)
		} finally {
			setLoadingProfile(false)
		}
	}

	useEffect(() => {
		fetchProfile()
	}, [])

	const handleLoginSuccess = (token) => {
		if (token) {
			setAccessToken(token)
			sessionStorage.setItem('accessToken', token)
		}
		// immediately fetch profile using the provided token
		fetchProfile(token)
	}

	const handleRegisterSuccess = (token) => {
		if (token) {
			setAccessToken(token)
			sessionStorage.setItem('accessToken', token)
		}
		fetchProfile(token)
	}

	const handleLogout = () => {
		// clear local state and stored token
		setProfile(null)
		setAccessToken(null)
		sessionStorage.removeItem('accessToken')
	}

	if (loadingProfile) {
		return <div className="app-root">Загрузка...</div>
	}

	if (profile) {
		return <Profile profile={profile} onLogout={handleLogout} />
	}

	return (
		<div className="app-root">
			<div className="auth-container">
				{mode === 'login' ? (
					<Login onSwitch={() => setMode('register')} onSuccess={handleLoginSuccess} />
				) : (
					<Register onSwitch={() => setMode('login')} onSuccess={handleRegisterSuccess} />
				)}
			</div>
		</div>
	)
}

export default App