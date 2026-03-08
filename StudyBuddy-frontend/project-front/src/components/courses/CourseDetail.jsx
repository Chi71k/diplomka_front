import React, { useState, useEffect } from 'react'
import { Link, useParams, useNavigate } from 'react-router-dom'
import { apiGetCourse, apiDeleteCourse } from '../../api'

const CourseDetail = () => {
  const { id } = useParams()
  const navigate = useNavigate()
  const [course, setCourse] = useState(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [deleting, setDeleting] = useState(false)

  useEffect(() => {
    let cancelled = false
    async function load() {
      setLoading(true)
      setError('')
      try {
        const data = await apiGetCourse(id)
        if (!cancelled) setCourse(data)
      } catch (e) {
        if (!cancelled) setError(e.status === 404 ? 'Course not found' : (e.error || 'Failed to load'))
      } finally {
        if (!cancelled) setLoading(false)
      }
    }
    load()
    return () => { cancelled = true }
  }, [id])

  const handleDelete = async () => {
    if (!window.confirm('Delete this course?')) return
    setDeleting(true)
    setError('')
    try {
      await apiDeleteCourse(id)
      navigate('/courses')
    } catch (e) {
      setError(e.error || 'Failed to delete')
    } finally {
      setDeleting(false)
    }
  }

  if (loading) {
    return (
      <div className="page-content">
        <div className="profile-loading">Loading...</div>
      </div>
    )
  }

  if (error && !course) {
    return (
      <div className="page-content">
        <div className="profile-card">
          <div className="auth-error">{error}</div>
          <Link to="/courses" className="btn btn-primary">Back to courses</Link>
        </div>
      </div>
    )
  }

  if (!course) return null

  return (
    <div className="page-content">
      <header className="page-header">
        <h1 className="page-title">{course.title}</h1>
        <p className="page-subtitle">{course.subject} · {course.level}</p>
      </header>

      <section className="profile-card">
        <div className="profile-info-row">
          <span className="profile-label">Subject</span>
          <span className="profile-value">{course.subject}</span>
        </div>
        <div className="profile-info-row">
          <span className="profile-label">Level</span>
          <span className="profile-value">{course.level}</span>
        </div>
        <div className="profile-info-row profile-info-bio">
          <span className="profile-label">Description</span>
          <span className="profile-value">{course.description}</span>
        </div>

        <div className="profile-form-actions profile-form-actions-top">
          <Link to={`/courses/${id}/edit`} className="btn btn-primary">Edit</Link>
          <button type="button" className="btn btn-danger" onClick={handleDelete} disabled={deleting}>
            {deleting ? 'Deleting...' : 'Delete course'}
          </button>
        </div>
        {error && <div className="auth-error auth-error-margin">{error}</div>}
      </section>
    </div>
  )
}

export default CourseDetail
