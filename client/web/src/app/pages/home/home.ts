import { Component, effect, ElementRef, inject, OnDestroy, OnInit, signal, viewChild } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { Router } from '@angular/router';
import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import { firstValueFrom, Subscription } from 'rxjs';
import { BreakpointObserver } from '@angular/cdk/layout';
import { MatButtonModule } from '@angular/material/button';
import { MatDialog } from '@angular/material/dialog';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { MatListModule } from '@angular/material/list';
import { MatSnackBar } from '@angular/material/snack-bar';
import { MatToolbarModule } from '@angular/material/toolbar';
import { MatTooltipModule } from '@angular/material/tooltip';
import { AuthService } from '../../services/auth.service';
import { WebSocketService } from '../../services/websocket.service';
import { ChatService } from '../../services/chat.service';
import { SoundService } from '../../services/sound.service';
import { ApiError, MessageEventBody } from '../../models';
import { NewChatDialog } from './new-chat-dialog';
import { environment } from '../../../environments/environment';

@Component({
  selector: 'app-home',
  imports: [
    FormsModule,
    MatToolbarModule,
    MatButtonModule,
    MatIconModule,
    MatListModule,
    MatFormFieldModule,
    MatInputModule,
    MatTooltipModule,
  ],
  templateUrl: './home.html',
  styleUrl: './home.scss',
})
export class Home implements OnInit, OnDestroy {
  private router = inject(Router);
  private http = inject(HttpClient);
  private dialog = inject(MatDialog);
  private snackBar = inject(MatSnackBar);
  private breakpointObserver = inject(BreakpointObserver);

  auth = inject(AuthService);
  ws = inject(WebSocketService);
  chat = inject(ChatService);
  sound = inject(SoundService);

  messageText = '';
  disconnected = signal(false);
  reconnecting = signal(false);
  isMobile = signal(false);
  currentTheme = signal<'azure-blue' | 'cyan-orange'>('azure-blue');

  private messagesEnd = viewChild<ElementRef>('messagesEnd');
  private subscriptions: Subscription[] = [];

  constructor() {
    effect(() => {
      this.chat.activeMessages();
      const el = this.messagesEnd()?.nativeElement;
      if (el) {
        setTimeout(() => el.scrollIntoView({ behavior: 'smooth' }));
      }
    });
  }

  ngOnInit(): void {
    this.subscriptions.push(
      this.breakpointObserver.observe('(max-width: 768px)').subscribe((result) => {
        this.isMobile.set(result.matches);
      }),
    );

    this.subscriptions.push(
      this.ws.connectionStatus$.subscribe((connected) => {
        this.disconnected.set(!connected);
      }),
    );

    this.subscriptions.push(
      this.ws.messages$.subscribe((event) => {
        if (event.event_type === 'MessageReceived') {
          const body = event.event_body as MessageEventBody;
          this.chat.addMessage(body.sender, {
            text: body.message,
            sender: body.sender,
            own: false,
          });
          this.sound.playReceiveSound();
        }
      }),
    );

    if (!this.ws.isConnected) {
      this.reconnect();
    }

    const savedTheme = localStorage.getItem('rb_theme') as 'azure-blue' | 'cyan-orange' | null;
    if (savedTheme) {
      this.currentTheme.set(savedTheme);
      this.applyTheme(savedTheme);
    }
  }

  ngOnDestroy(): void {
    this.subscriptions.forEach((s) => s.unsubscribe());
  }

  toggleTheme(): void {
    const next = this.currentTheme() === 'azure-blue' ? 'cyan-orange' : 'azure-blue';
    this.currentTheme.set(next);
    localStorage.setItem('rb_theme', next);
    this.applyTheme(next);
  }

  openNewChatDialog(): void {
    this.dialog
      .open(NewChatDialog, { width: '360px' })
      .afterClosed()
      .subscribe((username: string | undefined) => {
        if (username) {
          this.chat.startChat(username);
        }
      });
  }

  async sendMessage(): Promise<void> {
    const text = this.messageText.trim();
    const receiver = this.chat.activeUsername();
    if (!text || !receiver) return;

    try {
      await firstValueFrom(
        this.http.post(`${environment.backendUrl}/api/message`, { message: text, receivers: [receiver] }),
      );
      this.chat.addMessage(receiver, {
        text,
        sender: this.auth.getUsername()!,
        own: true,
      });
      this.messageText = '';
      this.sound.playSendSound();
    } catch (e) {
      const apiError = (e as HttpErrorResponse).error as ApiError | undefined;
      this.snackBar.open(apiError?.reason || 'Failed to send message', 'Close', {
        duration: 3000,
      });
    }
  }

  async reconnect(): Promise<void> {
    const creds = this.auth.getCredentials();
    if (!creds) {
      this.router.navigate(['/login']);
      return;
    }
    this.reconnecting.set(true);
    try {
      await this.ws.connect(creds.username, creds.password);
    } catch {
      this.snackBar.open('Reconnection failed. Try again.', 'Close', { duration: 3000 });
    } finally {
      this.reconnecting.set(false);
    }
  }

  logout(): void {
    this.ws.disconnect();
    this.auth.clearCredentials();
    this.chat.clearAll();
    this.router.navigate(['/login']);
  }

  private applyTheme(theme: string): void {
    document.documentElement.classList.toggle('cyan-orange-theme', theme === 'cyan-orange');
  }
}
