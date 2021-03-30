export interface UserBaseInfo {
    id: number;
    name: string;
    score: number;
    gender: string;
    avatarUrl: string;
}

export function userMap(users: Users): Map<number, UserBaseInfo> {
    const map = new Map<number, UserBaseInfo>();
    for (const user of users) {
        map.set(user.id, user)
    }
    return map
}

export function appendUserMap(map: Map<number, UserBaseInfo>, users: Users): Map<number, UserBaseInfo> {
    for (const user of users) {
        map.set(user.id, user)
    }
    return map
}

export type Users = UserBaseInfo[];