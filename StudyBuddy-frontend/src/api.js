const API_BASE = import.meta.env.VITE_API_BASE || ''

// Callback, который AuthContext регистрирует для разлогина при 401
let onUnauthorized = null
export function setOnUnauthorized(fn) {
  onUnauthorized = fn
}

function getToken() {
  return localStorage.getItem('accessToken')
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

// Обёртка над fetch: при 401 пробует обновить токен и повторяет запрос.
// При неудаче — вызывает onUnauthorized (разлогин).
async function apiFetch(url, options) {
  let res = await fetch(url, options)

  if (res.status === 401) {
    try {
      const refreshRes = await fetch(`${API_BASE}/api/v1/auth/refresh`, {
        method: 'POST',
        credentials: 'include',
      })
      if (!refreshRes.ok) throw new Error('refresh failed')

      const refreshData = await refreshRes.json().catch(() => ({}))
      if (refreshData.accessToken) {
        localStorage.setItem('accessToken', refreshData.accessToken)
      }

      // Повторяем исходный запрос с новым токеном
      const newToken = getToken()
      res = await fetch(url, {
        ...options,
        headers: {
          ...options.headers,
          ...(newToken ? { Authorization: `Bearer ${newToken}` } : {}),
        },
      })
    } catch {
      if (onUnauthorized) onUnauthorized()
      throw { status: 401, error: 'Session expired. Please log in again.' }
    }
  }

  return res
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
  const res = await apiFetch(`${API_BASE}/api/v1/users/me`, {
    method: 'GET',
    headers: authHeaders(),
    credentials: 'include',
  })
  return handleResponse(res)
}

export async function apiUpdateProfile(body) {
  const res = await apiFetch(`${API_BASE}/api/v1/users/me`, {
    method: 'PUT',
    headers: authHeaders(),
    credentials: 'include',
    body: JSON.stringify(body),
  })
  return handleResponse(res)
}

export async function apiDeleteProfile() {
  const res = await apiFetch(`${API_BASE}/api/v1/users/me`, {
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
  const res = await apiFetch(url, {
    method: 'GET',
    headers: authHeaders(),
    credentials: 'include',
  })
  return handleResponse(res)
}

export async function apiGetCourse(id) {
  const res = await apiFetch(`${API_BASE}/api/v1/courses/${id}`, {
    method: 'GET',
    headers: authHeaders(),
    credentials: 'include',
  })
  return handleResponse(res)
}

export async function apiCreateCourse(body) {
  const res = await apiFetch(`${API_BASE}/api/v1/courses`, {
    method: 'POST',
    headers: authHeaders(),
    credentials: 'include',
    body: JSON.stringify(body),
  })
  return handleResponse(res)
}

export async function apiUpdateCourse(id, body) {
  const res = await apiFetch(`${API_BASE}/api/v1/courses/${id}`, {
    method: 'PATCH',
    headers: authHeaders(),
    credentials: 'include',
    body: JSON.stringify(body),
  })
  return handleResponse(res)
}

export async function apiDeleteCourse(id) {
  const res = await apiFetch(`${API_BASE}/api/v1/courses/${id}`, {
    method: 'DELETE',
    headers: authHeaders(),
    credentials: 'include',
  })
  if (res.status === 204) return
  return handleResponse(res)
}

export { getToken }

// --- Interests ---
export async function apiGetInterestsCatalog() {
  const res = await apiFetch(`${API_BASE}/api/v1/interests`, {
    method: 'GET',
    headers: authHeaders(),
    credentials: 'include',
  })
  return handleResponse(res)  
}

export async function apiGetMyInterests() {
  const res = await apiFetch(`${API_BASE}/api/v1/users/me/interests`, {
    method: 'GET',
    headers: authHeaders(),
    credentials: 'include',
  })
  return handleResponse(res)  
}

export async function apiReplaceMyInterests(interests) {
  const res = await apiFetch(`${API_BASE}/api/v1/users/me/interests`, {
    method: 'PUT',
    headers: authHeaders(),
    credentials: 'include',
    body: JSON.stringify({ interest_ids: interests }),
  })
  return handleResponse(res)
}