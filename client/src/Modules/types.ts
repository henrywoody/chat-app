export type Room = {
    name: string;
    messages: Message[];
}

export type Message = {
    ID:         number;
	senderName: string;
	sentAt:     string;
	body:       string;
}
