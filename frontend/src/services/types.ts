export interface LoginResponse {
  token: string;
}

export interface UserResponse {
  id: string;
  name: string;
  email: string;
}

export interface ErrorResponse {
  error: string;
}
