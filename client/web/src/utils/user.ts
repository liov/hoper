export interface UserBase {
  id: number;
  name: string;
  score: number;
  gender: string;
  avatarUrl: string;
}

export type Users = UserBase[];

export function userMap(users: Users): Map<number, UserBase> {
  const map = new Map<number, UserBase>();
  for (const user of users) {
    map.set(user.id, user);
  }
  return map;
}

export function appendUserMap(
  map: Map<number, UserBase>,
  users: Users
): Map<number, UserBase> {
  for (const user of users) {
    map.set(user.id, user);
  }
  return map;
}
