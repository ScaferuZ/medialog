interface Log { id: string; media_id: string; status: string; rating?: number; note?: string; created_at: string }
export const useLogs = () => {
  const { request } = useApi(); const { getAuthHeaders } = useAuth()
  return {
    getMyLogs: (status?: string) => request(`/logs/me${status ? '?status='+status : ''}`, { headers: getAuthHeaders() }),
    createLog: (data: any) => request('/logs', { method: 'POST', body: data, headers: getAuthHeaders() }),
    updateLog: (id: string, data: any) => request(`/logs/${id}`, { method: 'PUT', body: data, headers: getAuthHeaders() }),
    deleteLog: (id: string) => request(`/logs/${id}`, { method: 'DELETE', headers: getAuthHeaders() })
  }
}
