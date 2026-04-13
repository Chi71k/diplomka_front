import { useState, useEffect } from "react"
import { useAuth } from "../../context/useAuth"
import {useToast} from "../../context/ToastContext"
import { apiGetMatchRequests, apiGetUserById } from "../../api"

const Partners = () => {
  const {profile} = useAuth()
  const toast = useToast()
  const [partners, setPartners] = useState([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    const load = async () => {
      setLoading(true)
      try {
        // 1. Получаем все принятые запросы
        const data = await apiGetMatchRequests({ status: "accepted", limit: 100 })
        const items = data.items || []

        // 2. Для каждого запроса определяем ID партнера
        const myId = profile?.id
        const partnerIds = items.map((r) => 
          r.requesterId === myId ? r.receiverId : r.requesterId
        )

        // 3. Получаем данные каждого партнера
        const results = await Promise.allSettled(
          partnerIds.map((id) => apiGetUserById(id))
        )

        // 4. Фильтруем успешные запросы и сохраняем партнеров
        const enriched = items
          .map((request, i) => {
            if (results[i].status !== "fulfilled") return null
            return {
              request,
              partner: results[i].value
            }
          })

        .filter(Boolean)

        setPartners(enriched)
      } catch (e) {
        toast.error(e.error || "Failed to load partners")
      } finally {
        setLoading(false)
      }
    }
    load()
  }, [profile?.id])

  const formatDate = (iso) =>
    new Date(iso).toLocaleDateString("en-US", {
      month: "short", day: "numeric", year: "numeric",
    })

    return (
      <div className="page-content">
        <header className="page-header">
          <h1 className="page-title">My Partners</h1>
          <p className="page-subtitle">
            {partners.length > 0
            ? `You have ${partners.length} study partner${partners.length > 1 ? "s" : ""}.`
            : 'Your accepted match requests will appear here once you have any.'}
          </p>
        </header>

        {loading && <div className="profile-loading">Loading partners...</div>} 

        {!loading && partners.length === 0 && (
          <section className="profile-card">
            <p className="page-muted">
              No partners yet. Go to {' '}
              <a href="/matching/candidates" className="link">Find partners</a> to send requests.
            </p>
          </section>
        )}

        {!loading && partners.map(({ request, partner }) => (
          <section
          key={request.id}
          className="profile-card"
          style={{ marginBottom: '16px' }}
        >
          <div style={{display: 'flex', alignItems: 'flex-start', gap: '16px', flexWrap: 'wrap'}}>

            {/* Аватар */}
            <div style={{
              width: '52px', height: '52px', borderRadius: '50%', flexShrink: 0,
              background: 'linear-gradient(135deg, #60a5fa 0%, #3b82f6 100%)',
              display: 'flex', alignItems: 'center', justifyContent: 'center',
              color: '#fff', fontWeight: 600, fontSize: '1.2rem', overflow: 'hidden',
              }}>
              {partner.avatarUrl
                ? <img src={partner.avatarUrl} alt="" style={{ width: '100%', height: '100%', objectFit: 'cover' }} />
                : (partner.firstName?.[0] || '?').toUpperCase()
              }
            </div>

            {/* Информация о партнере */}
            <div style={{flex: 1, minWidth: '0'}}>
              <strong style={{ fontSize: '1rem', color: '#1e293b'}}>
                {partner.firstName} {partner.lastName}
              </strong>
              <p style={{ margin: '2px 0 0', fontSize: '0.8rem', color: '#94a3b8' }}>
                Partners since {formatDate(request.updatedAt ?? request.createdAt)}
              </p>
              {partner.bio && (
                <p style={{margin: '6px 0 0', fontSize: '0.9rem', color: '#64748b'}}>
                  {partner.bio}
                </p>
              )}
              {request.message && (
                <p style={{ margin: '6px 0 0', fontSize: '0.85rem', color: '#475569', fontStyle: 'italic' }}>
                  "{request.message}"
                </p>
              )}
            </div>

          </div>
        </section>
        ))}
      </div>
    ) 
  }
export default Partners