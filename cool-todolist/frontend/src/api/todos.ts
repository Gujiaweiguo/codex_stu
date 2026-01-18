import axios from 'axios'
import type { Todo } from '@/types'

const api = axios.create({
  baseURL: '/api',
  timeout: 10000,
})

export const todoApi = {
  async getAll(): Promise<Todo[]> {
    const { data } = await api.get<Todo[]>('/todos')
    return data
  },

  async create(title: string, order: number): Promise<Todo> {
    const { data } = await api.post<Todo>('/todos', { title, order })
    return data
  },

  async update(id: number, updates: Partial<Pick<Todo, 'title' | 'completed' | 'order'>>): Promise<Todo> {
    const { data } = await api.put<Todo>(`/todos/${id}`, updates)
    return data
  },

  async toggle(id: number): Promise<Todo> {
    const { data } = await api.patch<Todo>(`/todos/${id}/toggle`)
    return data
  },

  async delete(id: number): Promise<void> {
    await api.delete(`/todos/${id}`)
  },

  async reorder(ids: number[]): Promise<void> {
    await api.post('/todos/reorder', { ids })
  },
}
