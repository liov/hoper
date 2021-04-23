import {
  createStore,
  createLogger,
  Store,
  useStore as baseUseStore,
} from "vuex";
import { AllState, RootState } from "./index.d";
import { user } from "./user";
import { content } from "./content";
import { InjectionKey } from "vue";

const state: any = {
  loading: false,
};

const mutations = {};

const actions = {};

const getters = {};

const modules = {
  user,
  content,
};

const plugins = [createLogger()];
const store = createStore<any>({
  state,
  mutations,
  actions,
  getters,
  modules,
  plugins,
});

/*
export const key: InjectionKey<Store<RootState>> = Symbol();

export function useStore<T = AllState>() {
  return baseUseStore<T>(key);
}
*/

export default store;
