import React from "react";

export default function useFetchState<T>(initialDataValue: T) {
    const [isLoading, setIsLoading] = React.useState(true);
    const [data, setData] = React.useState<T>(initialDataValue);
    const [isErrored, setIsErrored] = React.useState(false);

    return {
        isLoading,
        setIsLoading,
        data,
        setData,
        isErrored,
        setIsErrored,
    }
}