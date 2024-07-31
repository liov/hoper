//鉴权
import type { NavigationGuard } from "vue-router";

import axios from "axios";
import { userStore } from "@pc/store";

export const authenticated: NavigationGuard = (_to, _from, next) => {
  if (userStore.auth) next();
  else next({ name: "Login", query: { back: _to.path } });
};

export const completedAuthenticated: NavigationGuard = async (
  _to,
  _from,
  next,
) => {
  if (userStore.auth && userStore.auth.avatarUrl) next();
  else {
    const res = await axios.get(`/api/v1/user/0`);
    if (res.data.code == 0) userStore.auth = res.data.data.user;
    next({ name: "Login", query: { back: _to.path } });
  }
};
