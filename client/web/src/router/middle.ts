//鉴权
import type { NavigationGuard } from "vue-router";

import axios from "axios";
import { useUserStore } from "@/store/modules/user";

export const authenticated: NavigationGuard = (_to, _from, next) => {
  if (useUserStore().auth) next();
  else next({ name: "Login", query: { back: _to.path } });
};

export const completedAuthenticated: NavigationGuard = async (
  _to,
  _from,
  next
) => {
  if (useUserStore().auth && useUserStore().auth.avatarUrl) next();
  else {
    const res = await axios.get(`/api/user/0`);
    if (res.data.code == 0) useUserStore().auth = res.data.data.user;
    next({ name: "Login", query: { back: _to.path } });
  }
};
