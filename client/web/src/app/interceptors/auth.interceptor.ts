import { HttpErrorResponse, HttpInterceptorFn } from '@angular/common/http';
import { inject } from '@angular/core';
import { Router } from '@angular/router';
import { catchError, throwError } from 'rxjs';
import { AuthService } from '../services/auth.service';
import { WebSocketService } from '../services/websocket.service';
import { ChatService } from '../services/chat.service';

export const authInterceptor: HttpInterceptorFn = (req, next) => {
  const auth = inject(AuthService);
  const router = inject(Router);
  const ws = inject(WebSocketService);
  const chat = inject(ChatService);

  const header = auth.getBasicAuthHeader();
  const authReq = header ? req.clone({ setHeaders: { Authorization: header } }) : req;

  return next(authReq).pipe(
    catchError((error: HttpErrorResponse) => {
      if (error.status === 401) {
        ws.disconnect();
        auth.clearCredentials();
        chat.clearAll();
        router.navigate(['/login']);
      }
      return throwError(() => error);
    }),
  );
};
