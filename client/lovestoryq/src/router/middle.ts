//鉴权
import { NavigationGuard } from "vue-router";
import store from "@/store/index";
import axios from "axios";

export const authenticated: NavigationGuard = (_to, _from, next) => {
  if (store.state.user.auth) next();
  else next({ name: "Login", query: { back: _to.path } });
};

export const completedAuthenticated: NavigationGuard = async (
  _to,
  _from,
  next
) => {
  if (store.state.user.auth && (store.state.user.auth as any).avatarUrl) next();
  else {
    const res = await axios.get(`/api/v1/user/0`);
    if (res.data.code == 0) store.commit("setAuth", res.data.details.user);
    next({ name: "Login", query: { back: _to.path } });
  }
};
