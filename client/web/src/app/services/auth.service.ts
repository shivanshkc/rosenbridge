import { Injectable } from '@angular/core';

const AUTH_KEY = 'rb_auth';
const USERNAME_KEY = 'rb_username';

@Injectable({ providedIn: 'root' })
export class AuthService {
  saveCredentials(username: string, password: string): void {
    sessionStorage.setItem(AUTH_KEY, btoa(`${username}:${password}`));
    sessionStorage.setItem(USERNAME_KEY, username);
  }

  getCredentials(): { username: string; password: string } | null {
    const encoded = sessionStorage.getItem(AUTH_KEY);
    if (!encoded) return null;

    const decoded = atob(encoded);
    const colonIndex = decoded.indexOf(':');
    if (colonIndex === -1) return null;

    return {
      username: decoded.substring(0, colonIndex),
      password: decoded.substring(colonIndex + 1),
    };
  }

  getBasicAuthHeader(): string | null {
    const encoded = sessionStorage.getItem(AUTH_KEY);
    return encoded ? `Basic ${encoded}` : null;
  }

  getUsername(): string | null {
    return sessionStorage.getItem(USERNAME_KEY);
  }

  isAuthenticated(): boolean {
    return sessionStorage.getItem(AUTH_KEY) !== null;
  }

  clearCredentials(): void {
    sessionStorage.removeItem(AUTH_KEY);
    sessionStorage.removeItem(USERNAME_KEY);
  }
}
