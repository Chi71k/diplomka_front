import React, { useState, useEffect } from 'react'
import { Link, useParams, useNavigate } from 'react-router-dom'
import { apiGetCourse, apiCreateCourse, apiUpdateCourse } from '../../api'
import { useToast } from '../../context/ToastContext'

const CourseForm = ({ edit = false }) => {
  const { id } = useParams()
  const navigate = useNavigate()
  const toast = useToast()
  const [loading, setLoading] = useState(edit)
  const [saving, setSaving] = useState(false)
  const [validationError, setValidationError] = useState('')
  const [form, setForm] = useState({ title: '', description: '', subject: '', level: '' })

  useEffect(() => {
    if (!edit) return
    let cancelled = false
    async function load() {
      try {
        const data = await apiGetCourse(id)
        if (!cancelled) setForm({
          title: data.title || '',
          description: data.description || '',
          subject: data.subject || '',
          level: data.level || '',
        })
      } catch {
        if (!cancelled) toast.error('Course not found')
      } finally {
        if (!cancelled) setLoading(false)
      }
    }
    load()
    return () => { cancelled = true }
  }, [edit, id])

  const handleSubmit = async (e) => {
    e.preventDefault()
    if (!form.title || !form.description || !form.subject || !form.level) {
      setValidationError('Fill in all fields')
      return
    }
    setValidationError('')
    setSaving(true)
    try {
      if (edit) {
        await apiUpdateCourse(id, form)
        toast.success('Course updated')
        navigate('/courses')
      } else {
        const created = await apiCreateCourse(form)
        toast.success('Course created')
        navigate(`/courses/${created.id}`)
      }
    } catch (e) {
      toast.error(e.error || 'Failed to save')
    } finally {
      setSaving(false)
    }
  }

  if (edit && loading) {
    return (
      <div className="page-content">
        <div className="profile-loading">Loading...</div>
      </div>
    )
  }

  return (
    <div className="page-content">
      <header className="page-header">
        <h1 className="page-title">{edit ? 'Edit course' : 'New course'}</h1>
        <p className="page-subtitle">{edit ? 'Update course details' : 'Add a new study course'}</p>
      </header>

      <section className="profile-card">
        <form onSubmit={handleSubmit} className="profile-form">
          <label className="profile-label">Title</label>
          <input
            className="profile-input"
            placeholder="Course title"
            value={form.title}
            onChange={(e) => setForm((f) => ({ ...f, title: e.target.value }))}
            required
          />
          <label className="profile-label">Description</label>
          <textarea
            className="profile-input profile-textarea"
            placeholder="Description"
            value={form.description}
            onChange={(e) => setForm((f) => ({ ...f, description: e.target.value }))}
            rows={4}
            required
          />
          <label className="profile-label">Subject</label>
          <input
            className="profile-input"
            placeholder="Subject"
            value={form.subject}
            onChange={(e) => setForm((f) => ({ ...f, subject: e.target.value }))}
            required
          />
          <label className="profile-label">Level</label>
          <input
            className="profile-input"
            placeholder="Beginner, intermediate, advanced..."
            value={form.level}
            onChange={(e) => setForm((f) => ({ ...f, level: e.target.value }))}
            required
          />
          {validationError && <div className="auth-error">{validationError}</div>}
          <div className="profile-form-actions">
            <button type="submit" className="btn btn-primary" disabled={saving}>
              {saving ? 'Saving...' : (edit ? 'Save' : 'Create')}
            </button>
            <Link to={edit ? `/courses/${id}` : '/courses'} className="btn btn-secondary">
              Cancel
            </Link>
          </div>
        </form>
      </section>
    </div>
  )
}

export default CourseForm
