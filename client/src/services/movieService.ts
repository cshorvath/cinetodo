export interface Movie {
  id: number;
  title: string;
  originalTitle: string;
  director: string;
  year: number;
}

export interface UserMovie {
  id: number;
  title: string;
  originalTitle: string;
  director: string;
  year: number;
  seen: boolean;
}

import { API_BASE_URL } from '@/config/api';

const API_BASE = API_BASE_URL;

export const movieService = {
  async searchMovies(query: string, token: string): Promise<Movie[]> {
    const response = await fetch(`${API_BASE}/movie?query=${encodeURIComponent(query)}`, {
      headers: { 'Authorization': `Bearer ${token}` },
    });
    if (!response.ok) throw new Error('Search failed');
    return response.json();
  },

  async getUserMovies(token: string): Promise<UserMovie[]> {
    const response = await fetch(`${API_BASE}/user/movie`, {
      headers: { 'Authorization': `Bearer ${token}` },
    });
    if (!response.ok) throw new Error('Failed to fetch movies');
    return response.json();
  },

  async addMovie(movieId: number, token: string): Promise<void> {
    const response = await fetch(`${API_BASE}/user/movie/${movieId}`, {
      method: 'POST',
      headers: { 'Authorization': `Bearer ${token}` },
    });
    if (!response.ok) throw new Error('Failed to add movie');
  },

  async deleteMovie(movieId: number, token: string): Promise<void> {
    const response = await fetch(`${API_BASE}/user/movie/${movieId}`, {
      method: 'DELETE',
      headers: { 'Authorization': `Bearer ${token}` },
    });
    if (!response.ok) throw new Error('Failed to delete movie');
  },

  async toggleSeen(movieId: number, seen: boolean, token: string): Promise<void> {
    const response = await fetch(`${API_BASE}/user/movie/${movieId}`, {
      method: 'PATCH',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ seen }),
    });
    if (!response.ok) throw new Error('Failed to update movie');
  },
};
