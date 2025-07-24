import { apiClient, handleAxiosError } from "./base"

export interface AccountDetails {
  createdAt: Date
  email: string
  id: string
  name: string
}

interface UserService {
  getAccountDetails: () => Promise<AccountDetails>
  updateName: (name: string) => Promise<AccountDetails>
}

export const userService: UserService = {
  getAccountDetails: async () => {
    try {
      const response = await apiClient.get("/account")

      return response.data
    } catch (error) {
      return handleAxiosError(error)
    }
  },
  updateName: async (name) => {
    try {
      const response = await apiClient.put("/account", { name })

      return response.data
    } catch (error) {
      return handleAxiosError(error)
    }
  },
}
