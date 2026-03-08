import React from 'react'

const features = [
  {
    icon: (
      <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
        <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2" />
        <circle cx="9" cy="7" r="4" />
        <path d="M23 21v-2a4 4 0 0 0-3-3.87" />
        <path d="M16 3.13a4 4 0 0 1 0 7.75" />
      </svg>
    ),
    title: 'Find study partners',
    text: 'Connect with peers who share your courses and interests.',
  },
  {
    icon: (
      <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
        <path d="M4 19.5A2.5 2.5 0 0 1 6.5 17H20" />
        <path d="M6.5 2H20v20H6.5A2.5 2.5 0 0 1 4 19.5v-15A2.5 2.5 0 0 1 6.5 2z" />
        <path d="M8 7h8" />
        <path d="M8 11h8" />
      </svg>
    ),
    title: 'Track courses',
    text: 'Keep your current and past courses in one place.',
  },
  {
    icon: (
      <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
        <polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2" />
      </svg>
    ),
    title: 'Build reputation',
    text: 'Earn reviews and ratings from study sessions.',
  },
  {
    icon: (
      <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
        <polygon points="13 2 3 14 12 14 11 22 21 10 12 10 13 2" />
      </svg>
    ),
    title: 'Earn points',
    text: 'Get rewarded for being an active study partner.',
  },
]

export default function AuthLayout({ children }) {
  return (
    <div className="auth-split">
      <aside className="auth-promo">
        <div className="auth-promo-gradient" />
        <div className="auth-promo-content">
          <h1 className="auth-promo-brand">StudyBuddy</h1>
          <p className="auth-promo-tagline">Find your perfect study partner</p>
          <ul className="auth-promo-features">
            {features.map((f, i) => (
              <li key={i} className="auth-promo-feature">
                <span className="auth-promo-feature-icon">{f.icon}</span>
                <div>
                  <strong>{f.title}</strong>
                  <p>{f.text}</p>
                </div>
              </li>
            ))}
          </ul>
          <p className="auth-promo-copy">© {new Date().getFullYear()} StudyBuddy</p>
        </div>
      </aside>
      <main className="auth-form-panel">
        {children}
      </main>
    </div>
  )
}
