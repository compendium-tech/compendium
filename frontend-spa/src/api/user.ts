import { apiClient } from "./base"

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
    return (await apiClient.get("/account")).data
  },
  updateName: async (name) => {
    return (await apiClient.put("/account", { name })).data
  },
}
