import React from "react";
import { getRooms } from "./Modules/api";
import { Room } from "./Modules/types";
import useFetchState from "./Hooks/use-fetch-state";
import RoomsPanel from "./Components/RoomsPanel";
import ChatRoom from "./Components/ChatRoom";
import UserContext from "./Contexts/UserContext";
import styles from "./app.module.css";

export default function App() {
    const {
        isLoading,
        setIsLoading,
        data: rooms,
        setData: setRooms,
        isErrored,
        setIsErrored,
    } = useFetchState<Room[]>([]);
    const [activeRoomName, setActiveRoomName] = React.useState("");
    const [username, setUsername] = React.useState("");

    React.useEffect(() => {
        if (username === "") {
            const givenName = window.prompt("Set your username:");
            setUsername(givenName ?? "");
        }
    }, [username]);

    React.useEffect(() => {
        (async () => {
            setIsLoading(true);

            const { data: rooms, error } = await getRooms();
            if (error !== undefined) {
                setIsErrored(true);
                setRooms([]);
            } else {
                setIsErrored(false);
                setRooms(rooms!);
            }

            setIsLoading(false);
        })()
    }, []);

    return (
        <div className={styles["app"]}>
            <UserContext.Provider value={{ username }}>
                <RoomsPanel
                    isLoading={isLoading}
                    rooms={rooms}
                    isErrored={isErrored}
                    activeRoomName={activeRoomName}
                    onSelectActiveRoom={roomName => setActiveRoomName(roomName)}
                />

                { activeRoomName === "" ? (
                    <ChatRoomPlaceholder/>
                ) : (
                    <ChatRoom key={activeRoomName} name={activeRoomName}/>
                )}
            </UserContext.Provider>
        </div>
    );
}

function ChatRoomPlaceholder() {
    return (
        <div className={styles["chatroom-placeholder"]}>
            Select a room to get started.
        </div>
    )
}