import { Room, Message } from "./types";

export type APIResponse<T> = Partial<{
    data: T;
    error: string;
}>

export async function getRooms(): Promise<APIResponse<Room[]>> {
    const response = await fetch("/rooms");
    if (!response.ok) {
        const text = await response.text();
        return {error: text}
    }

    const json = await response.json();
    return {data: json};
}

export async function getRoom(name: string): Promise<APIResponse<Room>> {
    const response = await fetch("/rooms/" + name);
    if (!response.ok) {
        const text = await response.text();
        return {error: text}
    }

    const json = await response.json();
    return {data: json};
}

export async function sendMessage(roomName: string, senderName: string, body: string): Promise<APIResponse<Message>> {
    const payload = { roomName, senderName, body };
    const response = await fetch("/messages", {method: "POST", body: JSON.stringify(payload)});
    if (!response.ok) {
        const text = await response.text();
        return {error: text};
    }

    const json = await response.json();
    return {data: json};
}
