import { ObjMap } from "@/utils/user";
import { defineStore } from "pinia";

export interface ContentStore {
  moment: any;
  note?: any;
  diary?: any;
  diaryBook?: any;
  fav?: any;
  collect?: any;
  comment?: any;
  commentCache: Map<number, any[]>;
}

const state: ContentStore = {
  moment: null,
  commentCache: new Map<number, any[]>(),
};

const actions = {};

const getters = {
  getMoment(state) {
    return state.moment;
  },
  getCommentCache: (state) => (rooId) => {
    return state.commentCache.get(rooId);
  },
};

export const useContentStore = defineStore({
  id: "content",
  state: () => state,
  getters,
  actions,
});

export const contentMutations = [
  (place) => {
    throw new Error(`${place} not implemented`);
  },
  (moment) => (state.moment = moment),
  (note) => (state.note = note),
  (diary) => (state.diary = diary),
  (diaryBook) => (state.diary = diaryBook),
  (fav) => (state.fav = fav),
  (collect) => (state.collect = collect),
  (comment) => (state.comment = comment),
];
