<template>
  <div>
    <div>db</div>
    <div id="status" />
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from "vue";
import axios from "axios";

let db;
const msg = ref("");

const params = {
  pageNo: 0,
  pageSize: 5,
};
const { data } = await axios.get(`/api/v1/moment`, { params });
const momentList = reactive(data.data);
const total = ref(data.count);
const topCount = ref(data.top_count);
const user_action =
  data.user_action != null
    ? reactive(data.user_action)
    : reactive({
        collect: [],
        like: [],
        approve: [],
        comment: [],
        browse: [],
      });

onMounted(() => {
  db = window.openDatabase("hoper", "1.0", "hoper DB", 2 * 1024 * 1024);
  db.transaction(function (tx) {
    tx.executeSql("CREATE TABLE IF NOT EXISTS Moments (id unique, content)");
  });
  const vm = this;
  db.transaction(function (tx) {
    for (const item of momentList) {
      tx.executeSql(
        `INSERT INTO Moments (id, content)
         VALUES (${item.id}, '${item.content}')`,
      );
    }
  });

  db.transaction(function (tx) {
    tx.executeSql(
      "SELECT * FROM Moments",
      [],
      function (tx, results) {
        const len = results.rows.length;
        let i;
        let msg = "<p>查询记录条数: " + len + "</p>";
        document.querySelector("#status")!.innerHTML += msg;

        for (i = 0; i < len; i++) {
          msg = "<p><b>" + results.rows.item(i).content + "</b></p>";
          document.querySelector("#status")!.innerHTML += msg;
        }
      },
      null,
    );
  });
});
</script>

<style scoped></style>
