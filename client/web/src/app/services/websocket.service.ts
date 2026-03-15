import { Injectable, NgZone } from '@angular/core';
import { Observable, Subject } from 'rxjs';
import { SocketEvent } from '../models';
import { environment } from '../../environments/environment';

@Injectable({ providedIn: 'root' })
export class WebSocketService {
  private ws: WebSocket | null = null;
  private messageSubject = new Subject<SocketEvent>();
  private statusSubject = new Subject<boolean>();

  readonly messages$: Observable<SocketEvent> = this.messageSubject.asObservable();
  readonly connectionStatus$: Observable<boolean> = this.statusSubject.asObservable();

  constructor(private ngZone: NgZone) {}

  connect(username: string, password: string): Promise<void> {
    return new Promise((resolve, reject) => {
      this.disconnect();

      let base: string;
      if (environment.backendUrl) {
        base = environment.backendUrl.replace(/^http/, 'ws');
      } else {
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        base = `${protocol}//${window.location.host}`;
      }
      const params = `username=${encodeURIComponent(username)}&password=${encodeURIComponent(password)}`;
      const url = `${base}/api/connect?${params}`;

      this.ws = new WebSocket(url);
      let opened = false;

      this.ws.onopen = () => {
        opened = true;
        this.ngZone.run(() => {
          this.statusSubject.next(true);
          resolve();
        });
      };

      this.ws.onmessage = (event) => {
        this.ngZone.run(() => {
          try {
            const data: SocketEvent = JSON.parse(event.data);
            this.messageSubject.next(data);
          } catch {
            console.error('Failed to parse WebSocket message');
          }
        });
      };

      this.ws.onclose = () => {
        this.ngZone.run(() => {
          this.statusSubject.next(false);
          if (!opened) {
            reject(new Error('Connection failed. Please check your credentials.'));
          }
        });
      };
    });
  }

  disconnect(): void {
    if (this.ws) {
      this.ws.onclose = null;
      this.ws.close();
      this.ws = null;
    }
  }

  get isConnected(): boolean {
    return this.ws?.readyState === WebSocket.OPEN;
  }
}
