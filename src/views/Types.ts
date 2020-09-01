export interface RoomDetails {
    roomID: string;
    roomName: string;
}

export interface JoinedRoom extends RoomDetails {
    roomIcon: string;
}

export interface JoinRequest extends RoomDetails {
    requestingUserName: string;
    requestingUserID: string;
}

export interface RecentChatPreviewType extends RoomDetails {
    truncatedMessage: string; // Truncated message should be < 10 characters
}

export interface RecentChatPreview {
    [index: number]: RecentChatPreviewType;
}

export interface NewRoomRequest extends RoomDetails {
    requestingUserName: string;
    requestingUserID: string;
}

export interface RoomPageDetails {
    roomID: string;
    roomName: string;
    roomIcon: string;

    firstLoad: boolean;

    registeredUsers: string[];
    messages: Messages[];
}

export interface UsersIcon {
    [index: string]: string;
}

export interface UsersOnline {
    [index: string]: boolean;
}

export interface Messages {
    time: string;
    type: string; // Type denotes the message type, if INFO, FILE or MESSAGE type.
    name: string;
    userID: string;
    message: string;

    index: Int16Array;
}
