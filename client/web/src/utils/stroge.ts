
function storageLocal() {
  return {
    getItem<T>(key: string): T | null {
      const val = localStorage.getItem(key);
      if (val === null) return null;
      try {
        return JSON.parse(val) as T;
      } catch {
        return val as unknown as T;
      }
    },
    setItem(key: string, value: unknown): void {
      localStorage.setItem(key, JSON.stringify(value));
    },
    delItem(key: string): void {
      localStorage.removeItem(key);
    },
  };
}

export const storage = storageLocal();
