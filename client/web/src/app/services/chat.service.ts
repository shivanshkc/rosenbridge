import { computed, Injectable, signal } from '@angular/core';
import { ChatMessage } from '../models';

@Injectable({ providedIn: 'root' })
export class ChatService {
  private chatOrder = signal<string[]>([]);
  private chatMessages = signal<Record<string, ChatMessage[]>>({});
  private activeUser = signal<string | null>(null);

  readonly orderedChats = this.chatOrder.asReadonly();
  readonly activeUsername = this.activeUser.asReadonly();

  readonly activeMessages = computed(() => {
    const username = this.activeUser();
    if (!username) return [];
    return this.chatMessages()[username] ?? [];
  });

  startChat(username: string): void {
    if (!this.chatOrder().includes(username)) {
      this.chatOrder.update((order) => [username, ...order]);
      this.chatMessages.update((msgs) => ({ ...msgs, [username]: [] }));
    }
    this.activeUser.set(username);
  }

  selectChat(username: string): void {
    this.activeUser.set(username);
  }

  deselectChat(): void {
    this.activeUser.set(null);
  }

  addMessage(chatUsername: string, message: ChatMessage): void {
    if (!this.chatOrder().includes(chatUsername)) {
      this.chatOrder.update((order) => [chatUsername, ...order]);
    } else {
      this.chatOrder.update((order) => {
        const filtered = order.filter((u) => u !== chatUsername);
        return [chatUsername, ...filtered];
      });
    }

    this.chatMessages.update((msgs) => ({
      ...msgs,
      [chatUsername]: [...(msgs[chatUsername] ?? []), message],
    }));
  }

  clearAll(): void {
    this.chatOrder.set([]);
    this.chatMessages.set({});
    this.activeUser.set(null);
  }
}
