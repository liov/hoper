<template>
  <div id="chat">
    <div v-if="msgs.length > 0">
      <div v-for="(item, idx) in msgs" :key="idx">
        <div :class="item.sendUserId === user.id ? 'right' : 'left'">
          <van-popover
            teleport="#chat"
            :show="true"
            :placement="item.sendUserId === user.id ? 'left' : 'right'"
          >
            <span>{{ item.content }}</span>
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
  </div>
</template>

<script setup lang="ts">
import { API_HOST, Env, STATIC_DIR as staticDir } from "@h5/plugin/config";
import {
  reactive,
  ref,
  onMounted,
  nextTick,
  onUnmounted,
  onBeforeUnmount,
  onDeactivated,
  onActivated,
} from "vue";
import { useUserStore } from "@h5/store/user";
import { useRouter, useRoute } from "vue-router";
import type { Ref, UnwrapRef, UnwrapNestedRefs } from "vue";

const router = useRouter();
const route = useRoute();
const userStore = useUserStore();
const message = ref("");
let ws: WebSocket; // Our websocket
const newMsg = ref(""); // Holds new messages to be sent to the server
const recv = 0; // Email address used for grabbing an avatar
const msgs: Ref<UnwrapRef<any[]>> = ref([]);
const focus = ref(false);

const user = ref(userStore.auth);

onMounted(() => newWs());
onUnmounted(() => ws.close());

function newWs() {
  const apiHost = Env.PROD ? "https:" + API_HOST : API_HOST;
  ws = new WebSocket(apiHost.replace("http", "ws") + "/api/ws/chat");
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
    content: message.value, // Strip out html
  };
  msgs.value = msgs.value.concat(msg);
  ws.send(JSON.stringify(msg));
  message.value = "";
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
