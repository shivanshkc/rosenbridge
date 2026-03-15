import { Component, inject, signal } from '@angular/core';
import { FormBuilder, ReactiveFormsModule } from '@angular/forms';
import { Router, RouterLink } from '@angular/router';
import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { MatSnackBar } from '@angular/material/snack-bar';
import { AuthService } from '../../services/auth.service';
import { WebSocketService } from '../../services/websocket.service';

@Component({
  selector: 'app-login',
  imports: [
    ReactiveFormsModule,
    RouterLink,
    MatFormFieldModule,
    MatInputModule,
    MatButtonModule,
    MatIconModule,
    MatCardModule,
  ],
  templateUrl: './login.html',
  styleUrl: './login.scss',
})
export class Login {
  private fb = inject(FormBuilder);
  private auth = inject(AuthService);
  private ws = inject(WebSocketService);
  private router = inject(Router);
  private snackBar = inject(MatSnackBar);

  form = this.fb.group({ username: [''], password: [''] });
  hidePassword = signal(true);
  loading = signal(false);

  async onSubmit(): Promise<void> {
    const { username, password } = this.form.value;
    if (!username || !password) {
      this.snackBar.open('Please enter username and password', 'Close', { duration: 3000 });
      return;
    }

    this.loading.set(true);
    try {
      await this.ws.connect(username, password);
      this.auth.saveCredentials(username, password);
      this.router.navigate(['/home']);
    } catch {
      this.snackBar.open('Login failed. Please check your credentials.', 'Close', {
        duration: 3000,
      });
    } finally {
      this.loading.set(false);
    }
  }
}
