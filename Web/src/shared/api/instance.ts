import axios from "axios";
import { getTokenInfo } from "./token";

export const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL,
  timeout: 1000,
  withCredentials: true,
  headers: { 'Content-Type': 'application/json' }
})

api.interceptors.request.use(async (config) => {
  const token = getTokenInfo().token
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})