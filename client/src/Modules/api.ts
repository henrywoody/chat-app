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

export function connectToRoom(roomName: string, callback: (data: Message[]) => void) {
    const socket = new WebSocket("ws://localhost:8080/messages");
    function selectRoom() {
        socket.send(JSON.stringify({selectRoom: roomName}));
    }
    socket.addEventListener("open", selectRoom);
    socket.addEventListener("message", event => {
        const data = JSON.parse(event.data);
        callback(data);
    })

    function closeSocket(socket: WebSocket) {
        if (socket.readyState === WebSocket.CONNECTING) {
            socket.removeEventListener("open", selectRoom);
            socket.addEventListener("open", () => socket.close());
        } else if (socket.readyState === WebSocket.OPEN) {
            socket.close();
        }
    }

    return {socket, closeSocket}
}
