const API_BASE = import.meta.env.VITE_API_BASE || ''

function getToken() {
  return sessionStorage.getItem('accessToken')
}

function authHeaders() {
  const token = getToken()
  return {
    'Content-Type': 'application/json',
    ...(token ? { Authorization: `Bearer ${token}` } : {}),
  }
}

async function handleResponse(res) {
  const data = await res.json().catch(() => ({}))
  if (!res.ok) {
    throw { status: res.status, ...data }
  }
  return data
}

// --- Auth ---
export async function apiLogin(email, password) {
  const res = await fetch(`${API_BASE}/api/v1/auth/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, password }),
    credentials: 'include',
  })
  return handleResponse(res)
}

export async function apiRegister({ email, password, firstName, lastName }) {
  const res = await fetch(`${API_BASE}/api/v1/auth/register`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, password, firstName, lastName }),
    credentials: 'include',
  })
  return handleResponse(res)
}

// --- Users (profile) ---
export async function apiGetProfile() {
  const res = await fetch(`${API_BASE}/api/v1/users/me`, {
    method: 'GET',
    headers: authHeaders(),
    credentials: 'include',
  })
  return handleResponse(res)
}

export async function apiUpdateProfile(body) {
  const res = await fetch(`${API_BASE}/api/v1/users/me`, {
    method: 'PUT',
    headers: authHeaders(),
    credentials: 'include',
    body: JSON.stringify(body),
  })
  return handleResponse(res)
}

export async function apiDeleteProfile() {
  const res = await fetch(`${API_BASE}/api/v1/users/me`, {
    method: 'DELETE',
    headers: authHeaders(),
    credentials: 'include',
  })
  if (res.status === 204) return
  return handleResponse(res)
}

// --- Courses ---
export async function apiListCourses(params = {}) {
  const q = new URLSearchParams()
  if (params.subject) q.set('subject', params.subject)
  if (params.level) q.set('level', params.level)
  if (params.limit != null) q.set('limit', params.limit)
  if (params.offset != null) q.set('offset', params.offset)
  const query = q.toString()
  const url = `${API_BASE}/api/v1/courses${query ? `?${query}` : ''}`
  const res = await fetch(url, {
    method: 'GET',
    headers: authHeaders(),
    credentials: 'include',
  })
  return handleResponse(res)
}

export async function apiGetCourse(id) {
  const res = await fetch(`${API_BASE}/api/v1/courses/${id}`, {
    method: 'GET',
    headers: authHeaders(),
    credentials: 'include',
  })
  return handleResponse(res)
}

export async function apiCreateCourse(body) {
  const res = await fetch(`${API_BASE}/api/v1/courses`, {
    method: 'POST',
    headers: authHeaders(),
    credentials: 'include',
    body: JSON.stringify(body),
  })
  return handleResponse(res)
}

export async function apiUpdateCourse(id, body) {
  const res = await fetch(`${API_BASE}/api/v1/courses/${id}`, {
    method: 'PATCH',
    headers: authHeaders(),
    credentials: 'include',
    body: JSON.stringify(body),
  })
  return handleResponse(res)
}

export async function apiDeleteCourse(id) {
  const res = await fetch(`${API_BASE}/api/v1/courses/${id}`, {
    method: 'DELETE',
    headers: authHeaders(),
    credentials: 'include',
  })
  if (res.status === 204) return
  return handleResponse(res)
}

export { getToken }
