export async function errCatch1(asyncFunc) {
  try {
    const res = await asyncFunc();
    return [null, res];
  } catch (error) {
    return [error, null];
  }
}

export function errCatch2(asyncFunc): any[] {
  const value = new Array(2);
  asyncFunc()
    .catch((e) => {
      value[0] = e;
    })
    .then((res) => {
      value[1] = res;
    });
  return value;
}

export function errCatch(asyncFunc): [any, any] {
  let error;
  let result;
  asyncFunc()
    .catch((e) => (error = e))
    .then((res) => (result = res));
  return [error, result];
}
