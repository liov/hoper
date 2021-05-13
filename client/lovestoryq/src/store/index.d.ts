import {UserState} from "./user";

export interface RootState {
    loading: boolean;
}

export interface AllState extends RootState {
    user: UserState;
}
