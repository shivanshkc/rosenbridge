export interface SocketEvent {
  event_type: string;
  event_body: MessageEventBody;
}

export interface MessageEventBody {
  message: string;
  sender: string;
}

export interface ChatMessage {
  text: string;
  sender: string;
  own: boolean;
}

export interface ApiError {
  status: string;
  reason: string;
}
