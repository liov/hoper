
export function get(key:string): string | null{
  return localStorage.getItem(key)
}
