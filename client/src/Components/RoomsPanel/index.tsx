import React from "react";
import { Room } from "../../Modules/types";
import styles from "./style.module.css";

export type RoomsPanelProps = {
    isLoading: boolean;
    rooms: Room[];
    isErrored: boolean;
    activeRoomName: string;
    onSelectActiveRoom: (s: string) => void;
}

export default function RoomsPanel({isLoading, rooms, isErrored, activeRoomName, onSelectActiveRoom}: RoomsPanelProps) {
    const content = React.useMemo(() => {
        if (isLoading) {
            return (
                <React.Fragment>
                    Loading...
                </React.Fragment>
            )
        }

        if (isErrored) {
            return (
                <React.Fragment>
                    Error loading rooms.
                </React.Fragment>
            )
        }

        if (rooms.length === 0) {
            return (
                <React.Fragment>
                    No rooms found.
                </React.Fragment>
            )
        }

        return (
            <React.Fragment>
                {rooms.map(room => (
                    <RoomTab
                        key={room.name}
                        room={room}
                        isActive={room.name === activeRoomName}
                        onClick={() => onSelectActiveRoom(room.name)}
                    />
                ))}
            </React.Fragment>
        )
    }, [isLoading, isErrored, rooms, activeRoomName, onSelectActiveRoom]);

    return (
        <div className={styles["rooms-panel"]}>
            <div className={styles["rooms-panel__header"]}>
                Rooms
            </div>

            {content}
        </div>
    )
}

type RoomTabProps = {
    room: Room;
    isActive: boolean;
    onClick: () => void;
}

function RoomTab({room, isActive, onClick}: RoomTabProps) {
    let className = styles["rooms-panel__room-tab"];
    if (isActive) {
        className += " " + styles["rooms-panel__room-tab--active"];
    }

    return (
        <button className={className} onClick={onClick}>
            {room.name}
        </button>
    )
}