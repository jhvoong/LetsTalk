export interface RoomDetails {
    roomID: string;
    roomName: string;
}

export interface JoinedRoom extends RoomDetails {
    icon: string;
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