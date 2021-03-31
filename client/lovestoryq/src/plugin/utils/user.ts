export interface UserBaseInfo {
  id: number;
  name: string;
  score: number;
  gender: string;
  avatarUrl: string;
}

export type Users = UserBaseInfo[];

export function userMap(users: Users): Map<number, UserBaseInfo> {
  const map = new Map<number, UserBaseInfo>();
  for (const user of users) {
    map.set(user.id, user);
  }
  return map;
}

export function appendUserMap(
  map: Map<number, UserBaseInfo>,
  users: Users
): Map<number, UserBaseInfo> {
  for (const user of users) {
    map.set(user.id, user);
  }
  return map;
}

interface Obj {
  id: number;
}

export function appendObjMap<T extends Obj>(
  map: Map<number, T>,
  objs: T[]
): Map<number, T> {
  for (const obj of objs) {
    map.set(obj.id, obj);
  }
  return map;
}

export class ObjMap<T extends Obj> {
  _map = new Map<number, T>();

  appendMap(objs: T[]) {
    for (const obj of objs) {
      this._map.set(obj.id, obj);
    }
  }
  get(id: number): Obj {
    return this._map.get(id)!;
  }
}
