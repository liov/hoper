<template>
  <div v-if="msgs.length > 0">
    <div v-for="(item, idx) in msgs" :key="idx">
      <div :class="item.sendUserId == user.id ? 'right' : 'left'">
        <van-popover
          :show="true"
          :placement="item.sendUserId == user.id ? 'left' : 'right'"
        >
          <template #default>
            <span>{{ item.content }}</span>
          </template>
          <template #reference>
            <img class="avatar" :src="staticDir + user.avatarUrl" />
          </template>
        </van-popover>
      </div>
    </div>
    <div class="placeholder"></div>
  </div>
  <div class="input">
    <van-field
      v-model="message"
      rows="1"
      autosize
      type="textarea"
      placeholder="回复内容"
      :rules="[{ required: true, message: '输入内容为空' }]"
      @click-input="onFocus"
      @blur="onBlur"
      ref="commentRef"
    >
      <template #button>
        <div class="button">
          <van-button size="small" type="primary" @click="handleSubmit"
            >发送</van-button
          >
        </div>
      </template>
    </van-field>
  </div>
</template>

<script setup lang="ts">
import { API_HOST, STATIC_DIR as staticDir } from "@/plugin/config";
import { reactive, ref, onMounted, nextTick } from "vue";
import { useUserStore } from "@/store/user";
import { useRouter, useRoute } from "vue-router";

const router = useRouter();
const route = useRoute();

const userStore = useUserStore();
const message = ref("");
let ws: WebSocket; // Our websocket
const newMsg = ref(""); // Holds new messages to be sent to the server
const recv = 0; // Email address used for grabbing an avatar
const msgs = reactive([]);
const focus = ref(false);

const user = ref(userStore.auth);

onMounted(() => newWs());

function beforeDestroy() {
  ws.close();
}
function newWs() {
  ws = new WebSocket(
    document.location.protocol.replace("http", "ws") +
      "//" +
      API_HOST +
      "/api/ws/chat"
  );
  ws.onopen = () => {
    // console.log('建立websocket连接')
    if (message.value !== "") {
      handleSubmit();
    }
  };
  ws.onmessage = (evt) => {
    msgs.value = msgs.value.concat(JSON.parse(evt.data));
    console.log(msgs);
    message.value = "";
    nextTick(() => {
      document.querySelector(".placeholder")!.scrollIntoView();
    });
  };

  ws.onerror = () => {
    newWs();
  };
  ws.onclose = () => {
    console.log("websocket连接关闭");
  };
}
function handleSubmit() {
  if (message.value == "") {
    return;
  }

  if (ws.readyState !== 1) {
    newWs();
    return;
  }
  const msg: any = {
    recvUserId: recv,
    sendUserId: user.value.id,
    content: message, // Strip out html
  };
  msgs.value = msgs.value.concat(msg);
  ws.send(JSON.stringify(msg));
}
async function onFocus() {
  if (!userStore.auth) {
    await userStore.getAuth();
  }
  if (userStore.auth) focus.value = true;
  else
    await router.push({
      name: "Login",
      query: { back: route.path },
    });
}
function onBlur(e: FocusEvent) {
  if (!e.relatedTarget) focus.value = false;
}
</script>

<style scoped lang="less">
.left {
  text-align: left;
}
.right {
  text-align: right;
}
.input {
  position: fixed;
  bottom: 47px;
  width: 100%;
}
.button {
  display: grid;
}
.placeholder {
  height: 100px;
}
@avatar: 30px;
.avatar {
  flex-shrink: 0;
  width: @avatar;
  height: @avatar;
  border-radius: 40px;
  position: relative;
  margin: 0 16px;
}
</style>
