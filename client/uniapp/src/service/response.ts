
export interface Response<T> {
    code: number;
    message: string;
    details: T;
}