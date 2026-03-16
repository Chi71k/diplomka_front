import React from 'react'
import { useNavigate } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'

const Dashboard = () => {
  const { profile } = useAuth()
  const navigate = useNavigate()

  const firstName = profile?.firstName || ''

  const quickActions = [
    {
      id: 'update-profile',
      title: 'Update Profile',
      description: 'Keep your study info up to date so we can find better matches.',
      onClick: () => navigate('/profile'),
    },
  ]

  return (
    <div className="page-content">
      <header className="page-header">
        <h1 className="page-title">Welcome back, {firstName || 'student'} 👋</h1>
        <p className="page-subtitle">
          Here's what's happening with your study matches.
        </p>
      </header>

      <section className="dashboard-section">
        <div className="dashboard-grid">
          <div className="dashboard-card">
            <div className="dashboard-info">
              <h2 className="dashboard-card-title">
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  width="24"
                  height="24"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  strokeWidth="2"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                >
                  <path d="M4 14a1 1 0 0 1-.78-1.63l9.9-10.2a.5.5 0 0 1 .86.46l-1.92 6.02A1 1 0 0 0 13 10h7a1 1 0 0 1 .78 1.63l-9.9 10.2a.5.5 0 0 1-.86-.46l1.92-6.02A1 1 0 0 0 11 14z" />
                </svg>
                Points
              </h2>
              <p className="dashboard-card-subtitle">350</p>
            </div>
          </div>
          <div className="dashboard-card">
            <div className="dashboard-info">
              <h2 className="dashboard-card-title">
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  width="24"
                  height="24"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  strokeWidth="2"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                >
                  <path d="M11.525 2.295a.53.53 0 0 1 .95 0l2.31 4.679a2.123 2.123 0 0 0 1.595 1.16l5.166.756a.53.53 0 0 1 .294.904l-3.736 3.638a2.123 2.123 0 0 0-.611 1.878l.882 5.14a.53.53 0 0 1-.771.56l-4.618-2.428a2.122 2.122 0 0 0-1.973 0L6.396 21.01a.53.53 0 0 1-.77-.56l.881-5.139a2.122 2.122 0 0 0-.611-1.879L2.16 9.795a.53.53 0 0 1 .294-.906l5.165-.755a2.122 2.122 0 0 0 1.597-1.16z" />
                </svg>
                Rating
              </h2>
              <p className="dashboard-card-subtitle">4.6</p>
            </div>
          </div>
          <div className="dashboard-card">
            <div className="dashboard-info">
              <h2 className="dashboard-card-title">
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  width="24"
                  height="24"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  strokeWidth="2"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                >
                  <path d="M12 7v14"></path>
                  <path d="M3 18a1 1 0 0 1-1-1V4a1 1 0 0 1 1-1h5a4 4 0 0 1 4 4 4 4 0 0 1 4-4h5a1 1 0 0 1 1 1v13a1 1 0 0 1-1 1h-6a3 3 0 0 0-3 3 3 3 0 0 0-3-3z"></path>
                </svg>
                Courses
              </h2>
              <p className="dashboard-card-subtitle">3</p>
            </div>
          </div>
          <div className="dashboard-card">
            <div className="dashboard-info">
              <h2 className="dashboard-card-title">
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  width="24"
                  height="24"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  strokeWidth="2"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                >
                  <path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"></path>
                  <circle cx="9" cy="7" r="4"></circle>
                  <path d="M22 21v-2a4 4 0 0 0-3-3.87"></path>
                  <path d="M16 3.13a4 4 0 0 1 0 7.75"></path>
                </svg>
                Matches
              </h2>
              <p className="dashboard-card-subtitle">3</p>
            </div>
          </div>
        </div>
      </section>

      <section className="dashboard-section">
        <div className="dashboard-section-header">
          <h2 className="dashboard-section-title">Quick Actions</h2>
          <span className="dashboard-section-subtitle">
            Shortcuts to keep your profile and matches fresh.
          </span>
        </div>

        <div className="dashboard-actions dashboard-actions-grid">
          {quickActions.map((action) => (
            <div key={action.id} className="dashboard-action-card">
              <div className="dashboard-action-text">
                <h3 className="dashboard-action-title">{action.title}</h3>
                <p className="dashboard-action-meta">{action.description}</p>
              </div>
              <div className="dashboard-action-footer">
                <button
                  type="button"
                  className="btn btn-primary btn-sm"
                  onClick={action.onClick}
                >
                  {action.title}
                </button>
              </div>
            </div>
          ))}
        </div>
      </section>
    </div>
  )
}

export default Dashboard