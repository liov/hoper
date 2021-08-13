import { ObjMap } from "@/plugin/utils/user";

const state = {
  moment: null,
  commentCache: new Map<number, []>(),
};

const mutations = {
  setMoment: function (state, moment) {
    state.moment = moment;
  },
  setCommentCache: function (state, comment) {
    state.commentCache.clear();
  },
};

const actions = {};

const getters = {
  getMoment(state, getters, rootState) {
    return state.moment;
  },
  getCommentCache: (state, getters, rootState) => (rooId) => {
    return state.commentCache.get(rooId);
  },
};

export const content = {
  state,
  mutations,
  actions,
  getters,
};
