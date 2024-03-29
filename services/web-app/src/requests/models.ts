export type User = {
    id: string;
    username: string;
    email: string;
};

export type Room = {
    id: string;
    name: string;
};

export type LoginRes = {
    session_id: string;
    user_id: string;
};

export type wsMessage<T> = {
    content: T;
    error: any;
};

export type RoomEvent = {
    type: string;

    user_id: string;
    username: string;
    message: string;
};
