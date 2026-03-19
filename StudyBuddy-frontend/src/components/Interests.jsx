import React, { useEffect } from "react";
import {apiGetInterestsCatalog, apiGetMyInterests, apiReplaceMyInterests} from "../api";
import {useToast} from "../context/ToastContext";

export default function Interests() {
    const toast = useToast()
    const [catalog, setCatalog] = React.useState([])
    const [selected, setSelected] = React.useState(new Set())
    const [loading, setLoading] = React.useState(true)
    const [saving, setSaving] = React.useState(false)

    useEffect(() => {
        let canceled = false
        Promise.all([apiGetInterestsCatalog(), apiGetMyInterests()])
            .then(([all, my]) => {
                if (canceled) return
                setCatalog(all.items ?? [])
                setSelected(new Set((my.items ?? []).map(i => i.ID)))
            })
            .catch(() => {toast.error('Failed to load interests')})
            .finally(() => {if (!canceled) setLoading(false)})
        return () => {canceled = true}
    }, [])

    function toggle(id) {
        setSelected((prev) => {
            const next = new Set(prev)
            next.has(id) ? next.delete(id) : next.add(id)
            return next
        })
    }

    async function handleSave() {
        setSaving(true)
        try {
            await apiReplaceMyInterests([...selected])
            toast.success('Interests save')
        } catch (error) {
            toast.error('Failed to save interests')
        } finally {
            setSaving(false)
        }
    }

    if (loading) {
        return (
            <div className="page-content">
                <p className="profile-loading">Loading...</p>
            </div>
        )
    }

    return (
        <div className="page-content">
            <header className="page-header">
                <div className="page-header-row">
                    <div>
                        <h1 className="page-title">Interests</h1>
                        <p className="page-subtitle">Pick topics you want to study or teach</p>
                    </div>
                    <button className="btn btn-primary" onClick={handleSave} disabled={saving}>
                        {saving ? 'Saving...' : 'Save interests'}
                    </button>
                </div>
            </header>

            <section className="profile-card">
                <p className="page-muted" style={{marginBottom: '16px'}}>
                    {selected.size} topic{selected.size !== 1 ? 's' : ''} selected. Click to toggle.
                </p>
                <div className="interests-grid">
                    {catalog.map((item) => (
                      <button
                        key={item.ID}
                        type="button"
                        className={`interest-chip ${selected.has(item.ID) ? 'selected' : ''}`}
                        onClick={() => toggle(item.ID)}
                        >
                        {item.Name}
                      </button>  
                    ))}
                </div>
                {catalog.length === 0 && (
                    <p className="page-muted">No interests available in the catalog yet.</p>
                )}
            </section>
        </div>
    )
}