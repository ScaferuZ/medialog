export const useApi = () => {
  const config = useRuntimeConfig()
  const baseURL = config.public.apiBaseUrl

  const request = async <T>(path: string, options: Parameters<typeof $fetch<T>>[1] = {}) => {
    try {
      return await $fetch<T>(path, {
        baseURL,
        ...options
      })
    } catch (error) {
      const message = error instanceof Error ? error.message : 'API request failed'
      throw new Error(message)
    }
  }

  return {
    baseURL,
    request
  }
}
