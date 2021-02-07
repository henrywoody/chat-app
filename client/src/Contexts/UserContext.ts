import React from "react";

export type UserContextValue = {
    username: string;
}

const UserContext = React.createContext<UserContextValue>(undefined!);

export default UserContext;
