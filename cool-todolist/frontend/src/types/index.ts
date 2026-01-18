export interface Todo {
  id: number
  title: string
  completed: boolean
  order: number
  created_at: string
  updated_at: string
}

export type Theme = 'light' | 'dark' | 'auto'

export interface TodoStats {
  total: number
  completed: number
  pending: number
  completionRate: number
}
