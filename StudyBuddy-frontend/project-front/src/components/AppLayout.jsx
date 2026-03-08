import React from 'react'
import { NavLink, Outlet } from 'react-router-dom'

const navItems = [
  { to: '/profile', label: 'Profile', icon: 'person' },
  { to: '/courses', label: 'Courses', icon: 'book' },
]

const iconSvg = (icon) => {
  if (icon === 'person') {
    return (
      <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
        <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2" />
        <circle cx="12" cy="7" r="4" />
      </svg>
    )
  }
  if (icon === 'book') {
    return (
      <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
        <path d="M4 19.5A2.5 2.5 0 0 1 6.5 17H20" />
        <path d="M6.5 2H20v20H6.5A2.5 2.5 0 0 1 4 19.5v-15A2.5 2.5 0 0 1 6.5 2z" />
        <path d="M8 7h8" />
        <path d="M8 11h8" />
      </svg>
    )
  }
  return null
}

export default function AppLayout({ onLogout }) {
  return (
    <div className="app-layout">
      <aside className="app-sidebar">
        <div className="app-sidebar-brand">StudyBuddy</div>
        <nav className="app-sidebar-nav">
          {navItems.map(({ to, label, icon }) => (
            <NavLink
              key={to}
              to={to}
              className={({ isActive }) => `app-sidebar-link ${isActive ? 'app-sidebar-link-active' : ''}`}
            >
              <span className="app-sidebar-link-icon">{iconSvg(icon)}</span>
              {label}
            </NavLink>
          ))}
        </nav>
        <button type="button" className="app-sidebar-logout" onClick={() => onLogout && onLogout()}>
          <span className="app-sidebar-link-icon">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
              <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4" />
              <polyline points="16 17 21 12 16 7" />
              <line x1="21" y1="12" x2="9" y2="12" />
            </svg>
          </span>
          Log out
        </button>
      </aside>
      <main className="app-main">
        <Outlet />
      </main>
    </div>
  )
}
