import React, { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import { apiListCourses } from '../../api'
import { useAuth } from '../../context/AuthContext'

const CourseList = () => {
  const { profile } = useAuth()
  const [courses, setCourses] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [subject, setSubject] = useState('')
  const [level, setLevel] = useState('')

  const load = async () => {
    setLoading(true)
    setError('')
    try {
      const data = await apiListCourses({ subject: subject || undefined, level: level || undefined, limit: 50 })
      const all = Array.isArray(data) ? data : []
      setCourses(profile ? all.filter(c => c.ownerUserId === profile.id) : all)
    } catch (e) {
      setError(e.error || 'Failed to load courses')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    load()
  }, [subject, level, profile?.id])

  return (
    <div className="page-content">
      <header className="page-header">
        <div className="page-header-row">
          <div>
            <h1 className="page-title">Courses</h1>
            <p className="page-subtitle">Your study courses</p>
          </div>
          <Link to="/courses/new" className="btn btn-primary">New course</Link>
        </div>
      </header>

      <section className="profile-card">
        <div className="courses-filters">
          <input
            className="profile-input"
            placeholder="Subject (filter)"
            value={subject}
            onChange={(e) => setSubject(e.target.value)}
          />
          <input
            className="profile-input"
            placeholder="Level (filter)"
            value={level}
            onChange={(e) => setLevel(e.target.value)}
          />
        </div>

        {error && <div className="auth-error">{error}</div>}
        {loading && <div className="profile-loading">Loading...</div>}

        {!loading && courses.length === 0 && (
          <p className="page-muted">No courses yet. <Link to="/courses/new" className="auth-link">Create one</Link></p>
        )}

        {!loading && courses.length > 0 && (
          <ul className="course-list">
            {courses.map((c) => (
              <li key={c.id} className="course-list-item">
                <Link to={`/courses/${c.id}`} className="course-list-link">
                  <strong>{c.title}</strong>
                  <span className="course-meta">{c.subject} · {c.level}</span>
                </Link>
              </li>
            ))}
          </ul>
        )}
      </section>
    </div>
  )
}

export default CourseList
