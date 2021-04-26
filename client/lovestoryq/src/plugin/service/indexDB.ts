function openDB(name: string, version = 1) {
  const request = window.indexedDB.open(name);
  request.onerror = (e) => {
    console.log(e);
  };
  request.onsuccess = (e) => {
    console.log(e);
  };
  request.onupgradeneeded = (e) => {
    const db = request.result;
    if (!db.objectStoreNames.contains("moments")) {
      db.createObjectStore("moments", { keyPath: "id" });
    }
    console.log("DB version changed to " + version);
  };
}

function deleteDB(name) {
  indexedDB.deleteDatabase(name);
}
function addData(db: IDBDatabase, storeName, data) {
  const transaction = db.transaction(storeName, "readwrite");
  const store = transaction.objectStore(storeName);

  for (let i = 0; i < data.length; i++) {
    store.add(data[i]);
  }
}
function getDataByKey(db: IDBDatabase, storeName, key: string) {
  const transaction = db.transaction(storeName, "readwrite");
  const store = transaction.objectStore(storeName);
  const request = store.get(key);
  request.onsuccess = (e) => {
    const data = request.result;
    console.log(data.name);
  };
}
function updateDataByKey(db: IDBDatabase, storeName, key, newdata) {
  const transaction = db.transaction(storeName, "readwrite");
  const store = transaction.objectStore(storeName);
  const request = store.get(key);
  request.onsuccess = (e) => {
    const data = request.result;
    store.put(data.name);
  };
}
function deleteDataByKey(db: IDBDatabase, storeName, key) {
  const transaction = db.transaction(storeName, "readwrite");
  const store = transaction.objectStore(storeName);
  store.delete(key);
}
function clearObjectStore(db: IDBDatabase, storeName) {
  const transaction = db.transaction(storeName, "readwrite");
  const store = transaction.objectStore(storeName);
  store.clear();
}
