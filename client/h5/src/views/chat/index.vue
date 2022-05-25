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

<script lang="ts">
import { Options, Vue } from "vue-class-component";
import Moment from "@/components/moment/Moment.vue";
import ActionMore from "@/components/action/More.vue";
import { API_HOST, STATIC_DIR } from "@/plugin/config";

@Options({
  components: { Moment, ActionMore },
})
export default class Chat extends Vue {
  message = "";
  user = null;
  ws: WebSocket = null; // Our websocket
  newMsg = ""; // Holds new messages to be sent to the server
  recv = 0; // Email address used for grabbing an avatar
  msgs = [];
  focus = false;
  staticDir = STATIC_DIR;

  created() {
    this.user = this.$store.state.user.auth;
    this.newWs();
  }
  beforeDestroy() {
    this.ws.close();
  }
  newWs() {
    this.ws = new WebSocket(
      document.location.protocol.replace("http", "ws") +
        "//" +
        API_HOST +
        "/api/ws/chat"
    );
    this.ws.onopen = () => {
      // console.log('建立websocket连接')
      if (this.message !== "") {
        this.handleSubmit();
      }
    };
    this.ws.onmessage = (evt) => {
      this.msgs = this.msgs.concat(JSON.parse(evt.data));
      console.log(this.msgs);
      this.message = "";
      this.$nextTick(function () {
        document.querySelector(".placeholder").scrollIntoView();
      });
    };

    this.ws.onerror = () => {
      this.newWs();
    };
    this.ws.onclose = () => {
      console.log("websocket连接关闭");
    };
  }
  handleSubmit() {
    if (this.message == "") {
      return;
    }

    if (this.ws.readyState !== 1) {
      this.newWs();
      return;
    }
    const msg:any = {
      recvUserId: this.recv,
      sendUserId: this.user.id,
      content: this.message, // Strip out html
    };
    this.msgs = this.msgs.concat(msg);
    this.ws.send(JSON.stringify(msg));
  }
  async onFocus() {
    if (!this.$store.state.user.auth) {
      await this.$store.dispatch("getAuth");
    }
    if (this.$store.state.user.auth) this.focus = true;
    else
      await this.$router.push({
        name: "Login",
        query: { back: this.$route.path },
      });
  }
  onBlur(e: FocusEvent) {
    if (!e.relatedTarget) this.focus = false;
  }
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
