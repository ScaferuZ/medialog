interface ApiErrorPayload {
  error?: string
  details?: Record<string, string>
}

const isRecord = (value: unknown): value is Record<string, unknown> => {
  return typeof value === 'object' && value !== null
}

const getApiErrorMessage = (error: unknown): string => {
  if (isRecord(error)) {
    const data = error.data
    if (isRecord(data)) {
      const payload: ApiErrorPayload = {}

      if (typeof data.error === 'string') {
        payload.error = data.error
      }

      if (isRecord(data.details)) {
        const details = Object.entries(data.details).reduce<Record<string, string>>((acc, [key, value]) => {
          if (typeof value === 'string') {
            acc[key] = value
          }

          return acc
        }, {})

        if (Object.keys(details).length > 0) {
          payload.details = details
        }
      }

      if (payload.error === 'validation failed' && payload.details) {
        const firstDetail = Object.values(payload.details)[0]
        if (firstDetail) {
          return firstDetail
        }
      }

      if (payload.error) {
        return payload.error
      }
    }

    if (typeof error.statusMessage === 'string' && error.statusMessage.length > 0) {
      return error.statusMessage
    }

    if (typeof error.message === 'string' && error.message.length > 0) {
      return error.message
    }
  }

  if (error instanceof Error) {
    return error.message
  }

  return 'API request failed'
}

export const useApi = () => {
  const config = useRuntimeConfig()
  const baseURL = process.server ? config.apiBaseUrl : config.public.apiBaseUrl

  const request = async <T>(path: string, options: Parameters<typeof $fetch<T>>[1] = {}) => {
    try {
      return await $fetch<T>(path, {
        baseURL,
        ...options
      })
    } catch (error) {
      const message = getApiErrorMessage(error)
      throw new Error(message)
    }
  }

  return {
    baseURL,
    request
  }
}
