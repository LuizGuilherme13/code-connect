import type { LoginResponse, UserResponse, ErrorResponse } from './types';

const BASE_URL = '';

async function request<T>(endpoint: string, options?: RequestInit): Promise<T> {
  const res = await fetch(`${BASE_URL}${endpoint}`, {
    headers: { 'Content-Type': 'application/json' },
    ...options,
  });

  if (!res.ok) {
    const body = await res.json() as ErrorResponse;
    throw new Error(body.error || `Request failed with status ${res.status}`);
  }

  return res.json() as Promise<T>;
}

export function login(email: string, password: string) {
  return request<LoginResponse>('/api/login', {
    method: 'POST',
    body: JSON.stringify({ email, password }),
  });
}

export function register(name: string, email: string, password: string) {
  return request<UserResponse>('/api/register', {
    method: 'POST',
    body: JSON.stringify({ name, email, password }),
  });
}

export function getMe(token: string) {
  return request<UserResponse>('/api/users/me', {
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    },
  });
}
