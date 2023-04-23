<template>
  <div>IndexDB</div>
</template>

<script setup lang="ts">

import axios from "axios";
import { onMounted, reactive, ref } from "vue";

const DB:Record<string, any> = {
  name: "hoper",
  version: 1,
  db: null
};


const params = {
  pageNo: 0,
  pageSize: 5
};
const { data } = await axios.get(`/api/moment`, { params });
const momentList = reactive(data.data);
const total = ref(data.count);
const topCount = ref(data.top_count);
const user_action = data.user_action != null
  ? reactive(data.user_action)
  : reactive({
    collect: [],
    like: [],
    approve: [],
    comment: [],
    browse: []
  });
onMounted(async () => {
  DB.db = await openDB(DB.name, DB.version);
  addData(DB.db, "moments", momentList);
});

function openDB(name, version = 1) {
  return new Promise(function(resolve, reject) {
    const request = window.indexedDB.open(name);
    request.onerror = function(e) {
      reject(e);
    };
    request.onsuccess = function(e) {
      resolve(e.target!.result);
    };
    request.onupgradeneeded = function(e) {
      const db = e.target!.result;
      if (!db.objectStoreNames.contains("moments")) {
        db.createObjectStore("moments", { keyPath: "id" });
      }
      console.log("DB version changed to " + version);
    };
  });
}

function closeDB(db) {
  db.close();
}

function deleteDB(name) {
  indexedDB.deleteDatabase(name);
}

function addData(db, storeName, data) {
  const transaction = db.transaction(storeName, "readwrite");
  const store = transaction.objectStore(storeName);

  for (let i = 0; i < data.length; i++) {
    store.add(data[i]);
  }
}

,

function getDataByKey(db, storeName, key) {
  const transaction = db.transaction(storeName, "readwrite");
  const store = transaction.objectStore(storeName);
  const request = store.get(key);
  request.onsuccess = function(e) {
    const data = e.target.result;
    console.log(data.name);
  };
}

function updateDataByKey(db, storeName, key, newdata) {
  const transaction = db.transaction(storeName, "readwrite");
  const store = transaction.objectStore(storeName);
  const request = store.get(key);
  request.onsuccess = function(e) {
    let data = e.target.result;
    data = newdata;
    store.put(student);
  };
}

function deleteDataByKey(db, storeName, key) {
  const transaction = db.transaction(storeName, "readwrite");
  const store = transaction.objectStore(storeName);
  store.delete(key);
}

function clearObjectStore(db, storeName) {
  const transaction = db.transaction(storeName, "readwrite");
  const store = transaction.objectStore(storeName);
  store.clear();
}
}

</script>

<style scoped></style>
