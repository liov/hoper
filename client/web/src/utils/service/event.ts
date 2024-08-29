export function getEvent(url: string): EventSource {
  const client = new EventSource(url);
  client.onmessage = (evt) => {
    console.log(evt);
  };
  return client;
}
