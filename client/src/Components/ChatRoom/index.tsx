import React from "react";
import { Message } from "../../Modules/types";
import useFetchState from "../../Hooks/use-fetch-state";
import { connectToRoom } from "../../Modules/api";
import UserContext from "../../Contexts/UserContext";
import styles from "./style.module.css";

export type ChatRoomProps = {
    name: string;
}

export default function ChatRoom({name}: ChatRoomProps) {
    const {
        isLoading,
        setIsLoading,
        data: messages,
        setData: setMessages,
        isErrored,
    } = useFetchState<Message[]>([]);

    const socketRef = React.useRef<WebSocket | null>(null);
    const feedRef = React.useRef<HTMLDivElement>(null);

    React.useEffect(() => {
        const {socket, closeSocket} = connectToRoom(name, (data: Message[]) => {
            setMessages(data);
            setIsLoading(false);
        });
        socketRef.current = socket;

        return () => {
            closeSocket(socketRef.current!);
        }
    }, [name]);

    React.useEffect(() => {
        window.setTimeout(() => {
            if (feedRef.current === null) {
                return;
            }
            feedRef.current.scrollTo(0, feedRef.current.scrollHeight)
        }, 0);
    }, [messages]);

    const content = React.useMemo(() => {
        if (isLoading) {
            return <React.Fragment>Loading...</React.Fragment>
        }
        if (isErrored) {
            return <React.Fragment>Error loading messages.</React.Fragment>
        }
        return (
            <React.Fragment>
                {messages.map(message => (
                    <MessageDisplay key={message.ID} message={message}/>
                ))}
            </React.Fragment>
        )
    }, [isLoading, messages, isErrored]);

    return (
        <div className={styles["chatroom"]}>
            <div ref={feedRef} className={styles["chatroom__message-feed"]}>
                {content}
            </div>

            <MessageComposer roomName={name} socket={socketRef.current}/>
        </div>
    )
}

type MessageDisplayProps = {
    message: Message;
}

function MessageDisplay({message}: MessageDisplayProps) {
    const { username } = React.useContext(UserContext);
    const isFromUser = message.senderName === username;

    let className = styles["message-feed__message"];
    if (isFromUser) {
        className += " " + styles["message-feed__message--from-user"];
    }

    return (
        <div className={className}>
            <div className={styles["message-feed__message__sender-name"]}>{message.senderName}</div>
            <p className={styles["message-feed__message__body"]} title={message.sentAt}>
                {message.body}
            </p>
        </div>
    )
}

type MessageComposerProps = {
    roomName: string;
    socket: WebSocket | null
}

function MessageComposer({roomName, socket}: MessageComposerProps) {
    const { username } = React.useContext(UserContext);

    const [body, setBody] = React.useState("");
    const [isSubmitting, setIsSubmitting] = React.useState(false);
    const textAreaRef = React.useRef<HTMLTextAreaElement>(null);

    React.useEffect(() => {
        setBody("");
    }, [roomName]);

    const onSubmit = React.useCallback(async (event: React.FormEvent) => {
        event.preventDefault();

        if (body === "") {
            return;
        }
        setIsSubmitting(true);
        if (socket !== null && socket.readyState === WebSocket.OPEN) {
            socket?.send(JSON.stringify({roomName, senderName: username, body}));
            setBody("");
            if (textAreaRef.current !== null) {
                textAreaRef.current.focus();
            }
        } else {
            console.log("unable to send message; socket: ", socket);
        }
        setIsSubmitting(false);
    }, [roomName, username, body]);

    return (
        <div className={styles["chatroom__message-composer"]}>
            <form onSubmit={onSubmit}>
                <textarea ref={textAreaRef} name="body" value={body} onChange={e => setBody(e.target.value)} autoFocus/>
                <button type="submit" disabled={isSubmitting}>Send</button>
            </form>
        </div>
    )
}