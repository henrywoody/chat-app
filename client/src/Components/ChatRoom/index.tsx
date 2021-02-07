import React from "react";
import { Message } from "../../Modules/types";
import useFetchState from "../../Hooks/use-fetch-state";
import { getRoom, sendMessage } from "../../Modules/api";
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
        setIsErrored,
    } = useFetchState<Message[]>([]);

    React.useEffect(() => {
        (async () => {
            setIsLoading(true);

            const { data: room, error } = await getRoom(name);
            if (error !== undefined) {
                setIsErrored(true);
                setMessages([]);
            } else {
                setIsErrored(false);
                setMessages(room!.messages);
            }

            setIsLoading(false);
        })();
    }, [name]);

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
            <div className={styles["chatroom__message-feed"]}>
                {content}
            </div>

            <MessageComposer roomName={name}/>
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
}

function MessageComposer({roomName}: MessageComposerProps) {
    const { username } = React.useContext(UserContext);

    const [body, setBody] = React.useState("");
    const [isSubmitting, setIsSubmitting] = React.useState(false);

    React.useEffect(() => {
        setBody("");
    }, [roomName]);

    const onSubmit = React.useCallback(async (event: React.FormEvent) => {
        event.preventDefault();

        if (body === "") {
            return;
        }
        setIsSubmitting(true);
        const {error} = await sendMessage(roomName, username, body);
        if (error === undefined) {
            setBody("");
        }
        setIsSubmitting(false);
    }, [roomName, username, body]);

    return (
        <div className={styles["chatroom__message-composer"]}>
            <form onSubmit={onSubmit}>
                <textarea name="body" value={body} onChange={e => setBody(e.target.value)}/>
                <button type="submit" disabled={isSubmitting}>Send</button>
            </form>
        </div>
    )
}