export interface UserBaseInfo {
    id: number;
    name: string;
    score: number;
    gender: string;
    avatarUrl: string;
}

export type Users = UserBaseInfo[];

type GConstructor<T = Users> = new (...args: any[]) => T;

function UserMap<TBase extends GConstructor>(Base: TBase) {
    return class UserMapping extends Base {
        // Mixins may not declare private/protected properties
        // however, you can use ES2020 private fields
        _map = new Map<number, UserBaseInfo>();

        appendUserMap(users: Users) {
            for (const user of users) {
                this._map.set(user.id, user);
            }
        }
        getUser(id: number): UserBaseInfo {
            return this._map.get(id)!;
        }
    };
}
