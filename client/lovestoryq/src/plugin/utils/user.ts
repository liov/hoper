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

interface Obj<T> {
  id: T;
}

export function appendObjMap<T extends Obj<any>>(
  map: Map<number, T>,
  objs: T[]
): Map<number, T> {
  for (const obj of objs) {
    map.set(obj.id, obj);
  }
  return map;
}

export class ObjMap<K, V extends Obj<any>> {
  _map = new Map<K, V>();

  append(objs: V[]) {
    for (const obj of objs) {
      this._map.set(obj.id, obj);
    }
  }
  get(id: K): Obj<any> {
    return this._map.get(id)!;
  }
  set(id: K, v: V) {
    this._map.set(id, v);
  }
  has(key: K): boolean{
    return this._map.has(key)
  }
}
