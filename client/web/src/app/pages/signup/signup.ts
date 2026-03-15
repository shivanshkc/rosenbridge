import { Component, inject, signal } from '@angular/core';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { Router, RouterLink } from '@angular/router';
import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import { firstValueFrom } from 'rxjs';
import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { MatSnackBar } from '@angular/material/snack-bar';
import { AuthService } from '../../services/auth.service';
import { WebSocketService } from '../../services/websocket.service';
import { FormValidators } from '../../validators/form-validators';
import { ApiError } from '../../models';
import { environment } from '../../../environments/environment';

@Component({
  selector: 'app-signup',
  imports: [
    ReactiveFormsModule,
    RouterLink,
    MatFormFieldModule,
    MatInputModule,
    MatButtonModule,
    MatIconModule,
    MatCardModule,
  ],
  templateUrl: './signup.html',
  styleUrl: './signup.scss',
})
export class Signup {
  private fb = inject(FormBuilder);
  private http = inject(HttpClient);
  private auth = inject(AuthService);
  private ws = inject(WebSocketService);
  private router = inject(Router);
  private snackBar = inject(MatSnackBar);

  form = this.fb.group({
    username: ['', [Validators.required, FormValidators.username()]],
    password: ['', [Validators.required, FormValidators.password()]],
    repeatPassword: ['', [Validators.required, FormValidators.matchField('password')]],
  });

  hidePassword = signal(true);
  hideRepeatPassword = signal(true);
  loading = signal(false);

  constructor() {
    this.form.get('password')?.valueChanges.subscribe(() => {
      this.form.get('repeatPassword')?.updateValueAndValidity();
    });
  }

  getFieldError(field: string): string | null {
    const control = this.form.get(field);
    if (!control?.touched || !control.errors) return null;

    if (control.errors['required']) {
      const label = field === 'repeatPassword' ? 'Repeat password' : field.charAt(0).toUpperCase() + field.slice(1);
      return `${label} is required`;
    }
    if (control.errors['username']) return control.errors['username'];
    if (control.errors['password']) return control.errors['password'];
    if (control.errors['mismatch']) return control.errors['mismatch'];
    return null;
  }

  async onSubmit(): Promise<void> {
    this.form.markAllAsTouched();
    if (this.form.invalid) return;

    const { username, password } = this.form.value;
    this.loading.set(true);

    try {
      await firstValueFrom(this.http.post(`${environment.backendUrl}/api/user`, { username, password }));
    } catch (e) {
      const apiError = (e as HttpErrorResponse).error as ApiError;
      this.snackBar.open(apiError?.reason || 'Signup failed', 'Close', { duration: 3000 });
      this.loading.set(false);
      return;
    }

    try {
      await this.ws.connect(username!, password!);
      this.auth.saveCredentials(username!, password!);
      this.router.navigate(['/home']);
    } catch {
      this.snackBar.open('Account created but connection failed. Please log in.', 'Close', {
        duration: 4000,
      });
      this.router.navigate(['/login']);
    } finally {
      this.loading.set(false);
    }
  }
}
