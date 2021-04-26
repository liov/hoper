<template>
<div v-if="msgs.length > 0">
<van-popover v-for="(item,idx) in msgs" :key="idx" :show="true" :placement="item.sendUserId==user.Id?'left':'right'" >
  <span>{{item.msg}}</span>
  <template #reference>
    <img class="avatar" :src="user.avatarUrl" />
  </template>
</van-popover>
</div>
</template>

<script lang="ts">
import { Options, Vue } from "vue-class-component";
import Moment from "@/components/moment/Moment.vue";
import ActionMore from "@/components/action/More.vue";

@Options({
  components: { Moment, ActionMore },
})
export default class Chat extends Vue {
  submitting = false;
  value = "";
  user = null;
  ws:WebSocket = null; // Our websocket
  newMsg = ""; // Holds new messages to be sent to the server
  recipient = 0; // Email address used for grabbing an avatar
  msgs = [];

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
        window.location.host +
        ":8090/api/ws/chat"
    );
    this.ws.onopen = () => {
      // console.log('建立websocket连接')
      if (this.value !== "") {
        this.handleSubmit();
      }
    };
    this.ws.onmessage = (evt) => {
      this.submitting = false;
      this.msgs = [...this.msgs, JSON.parse(evt.data)];
      this.value = "";
      this.$nextTick(function () {
        document.querySelector("#bottom").scrollIntoView();
      });
    };

    this.ws.onerror = () => {
      this.newWs();
    };
    this.ws.onclose = () => {
      // console.log('websocket连接关闭')
    };
  }
  handleSubmit() {
    if (!this.value) {
      return;
    }

    if (this.ws.readyState !== 1) {
      this.newWs();
      return;
    }
    this.submitting = true;

    this.ws.send(
      JSON.stringify({
        recipient_user_id: this.recipient,
        sender_user_id:
          this.user !== null
            ? this.user.id
            : parseInt(localStorage.getItem("user")),
        content: this.value, // Strip out html
      })
    );
  }
  handleChange(e) {
    this.value = e.target.value;
  }
}
</script>

<style scoped lang="less">
</style>