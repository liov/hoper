//鉴权
import type { NavigationGuard } from "vue-router";

import axios from "axios";
import { state } from "@/store/user";

export const authenticated: NavigationGuard = (_to, _from, next) => {
  if (state.auth) next();
  else next({ name: "Login", query: { back: _to.path } });
};

export const completedAuthenticated: NavigationGuard = async (
  _to,
  _from,
  next
) => {
  if (state.auth && (state.auth as any).avatarUrl) next();
  else {
    const res = await axios.get(`/api/v1/user/0`);
    if (res.data.code == 0) state.auth = res.data.details.user;
    next({ name: "Login", query: { back: _to.path } });
  }
};
