export type UserInfo = {
  id: number;
  /** 头像 */
  avatar?: string;
  /** 用户名 */
  name: string;
  account?: string;
  /** 联系电话 */
  phone: string;
  roles: Array<string>;
  role: number;
  permissions: Array<string>;
};
