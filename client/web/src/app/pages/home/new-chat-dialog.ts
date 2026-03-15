import { Component, inject } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatDialogModule, MatDialogRef } from '@angular/material/dialog';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';

@Component({
  selector: 'app-new-chat-dialog',
  imports: [FormsModule, MatDialogModule, MatFormFieldModule, MatInputModule, MatButtonModule],
  template: `
    <h2 mat-dialog-title>New Chat</h2>
    <mat-dialog-content>
      <mat-form-field appearance="outline" style="width: 100%; margin-top: 8px">
        <mat-label>Username</mat-label>
        <input matInput [(ngModel)]="username" (keyup.enter)="submit()" cdkFocusInitial />
      </mat-form-field>
    </mat-dialog-content>
    <mat-dialog-actions align="end">
      <button mat-button mat-dialog-close>Cancel</button>
      <button mat-flat-button (click)="submit()" [disabled]="!username.trim()">Start Chat</button>
    </mat-dialog-actions>
  `,
})
export class NewChatDialog {
  username = '';
  private dialogRef = inject(MatDialogRef<NewChatDialog>);

  submit(): void {
    const trimmed = this.username.trim();
    if (trimmed) {
      this.dialogRef.close(trimmed);
    }
  }
}
