const state = {
  moment: null,
};

const mutations = {
  SET_MOMENT: function (state, moment) {
    state.moment = moment;
  },
};

const actions = {};

const getters = {
  getMoment(state, getters, rootState) {
    return state.moment;
  },
};

export const moment = {
  state,
  mutations,
  actions,
  getters,
};
