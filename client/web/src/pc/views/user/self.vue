<template>
  <a-row>
    <a-col :span="5" style="text-align: right">
      <a-avatar shape="square" :size="100" :src="user.avatar_url" />
      <br />
      <a-upload
        v-show="edit"
        name="file"
        action="/api/upload/avatar"
        :before-upload="beforeUpload"
        @change="uploadAvatarChange"
      >
        <a-button>
          <a-icon type="upload" />
          上传头像
        </a-button>
      </a-upload>
      <a-upload
        v-show="edit"
        name="file"
        action="/api/upload/cover"
        :before-upload="beforeUpload"
        @change="uploadCoverChange"
      >
        <a-button>
          <a-icon type="upload" />
          上传背景
        </a-button>
      </a-upload>
      <a-button v-show="!edit" icon="edit" @click="editable">
        修改资料
      </a-button>
      <a-button v-show="edit" icon="to-top" @click="commit">
        提交修改
      </a-button>
    </a-col>
    <a-col
      :span="16"
      :style="{
        background: 'url(' + user.cover_url + ') no-repeat',
        backgroundSize: 'cover',
      }"
    >
      <a-form-item
        label="用户名"
        :label-col="formItemLayout.labelCol"
        :wrapper-col="formItemLayout.wrapperCol"
      >
        <a-input v-model="user.name" disabled>
          <a-icon slot="prefix" type="user" />
        </a-input>
      </a-form-item>
      <a-row>
        <a-col :span="10">
          <a-form-item
            label="性别"
            :label-col="{ span: 7 }"
            :wrapper-col="{ span: 10 }"
          >
            <a-radio-group v-model="user.sex" :disabled="!edit">
              <a-radio-button value="男"> 男 </a-radio-button>
              <a-radio-button value="女"> 女 </a-radio-button>
            </a-radio-group>
          </a-form-item>
        </a-col>
        <a-col :span="13">
          <a-form-item
            label="生日"
            :label-col="{ span: 3 }"
            :wrapper-col="{ span: 10 }"
          >
            <a-date-picker
              v-model="user.birthday"
              show-time
              :disabled="!edit"
            />
          </a-form-item>
        </a-col>
      </a-row>

      <a-form-item
        label="邮箱"
        :label-col="formItemLayout.labelCol"
        :wrapper-col="formItemLayout.wrapperCol"
      >
        <a-input v-model="user.email" disabled>
          <a-icon slot="prefix" type="mail" />
        </a-input>
      </a-form-item>
      <a-form-item
        label="手机"
        :label-col="formItemLayout.labelCol"
        :wrapper-col="formItemLayout.wrapperCol"
      >
        <a-input v-model="user.phone" disabled>
          <a-icon slot="prefix" type="phone" />
        </a-input>
      </a-form-item>
      <a-form-item
        label="个人简介"
        :label-col="formItemLayout.labelCol"
        :wrapper-col="formItemLayout.wrapperCol"
      >
        <a-input v-model="user.introduction" :disabled="!edit">
          <a-icon slot="prefix" type="profile" />
        </a-input>
      </a-form-item>
      <a-form-item
        label="个人签名"
        :label-col="formItemLayout.labelCol"
        :wrapper-col="formItemLayout.wrapperCol"
      >
        <a-input v-model="user.signature" :disabled="!edit">
          <a-icon slot="prefix" type="profile" />
        </a-input>
      </a-form-item>
      <a-form-item
        label="地址"
        :label-col="formItemLayout.labelCol"
        :wrapper-col="formItemLayout.wrapperCol"
      >
        <a-cascader
          v-model="address"
          style="width: 50%"
          :options="options"
          placeholder="请选择"
          :disabled="!edit"
          @change="addrChange"
        >
          <a-icon slot="prefix" type="smile" />
        </a-cascader>
      </a-form-item>
      <a-form-item
        v-show="edu_exps.length > 0"
        label="教育经历"
        :label-col="formItemLayout.labelCol"
        :wrapper-col="{ span: 18 }"
      >
        <a-input-group v-for="(item, index) in edu_exps" :key="index">
          <a-col :span="6">
            <a-input
              v-model="edu_exps[index].school"
              :disabled="!edit"
              placeholder="请输入学校!"
            />
          </a-col>
          <a-col :span="6">
            <a-input
              v-model="edu_exps[index].speciality"
              :disabled="!edit"
              placeholder="请输入专业!"
            />
          </a-col>
          <a-col :span="9">
            <a-range-picker
              v-model="edu_exps[index].time"
              :disabled="!edit"
              @change="eduTime(index)"
            >
              <a-icon slot="suffixIcon" type="smile" />
            </a-range-picker>
          </a-col>
          <a-button v-show="index === 0" icon="plus" @click="changeEdu(-1)" />
          <a-button v-show="index > 0" icon="minus" @click="changeEdu(index)" />
        </a-input-group>
      </a-form-item>
      <a-form-item
        v-show="work_exps.length > 0"
        label="职业经历"
        :label-col="formItemLayout.labelCol"
        :wrapper-col="{ span: 18 }"
      >
        <a-input-group
          v-for="(item, index) in work_exps"
          :key="index + edu_exps.length"
        >
          <a-col :span="6">
            <a-input
              v-model="work_exps[index].company"
              :disabled="!edit"
              placeholder="请输入公司!"
            />
          </a-col>
          <a-col :span="6">
            <a-input
              v-model="work_exps[index].title"
              :disabled="!edit"
              placeholder="请输入职务!"
            />
          </a-col>
          <a-col :span="9">
            <a-range-picker
              v-model="work_exps[index].time"
              :disabled="!edit"
              @change="workTime(index)"
            >
              <a-icon slot="suffixIcon" type="smile" />
            </a-range-picker>
          </a-col>
          <a-button v-show="index === 0" icon="plus" @click="changeWork(-1)" />
          <a-button
            v-show="index > 0"
            icon="minus"
            @click="changeWork(index)"
          />
        </a-input-group>
      </a-form-item>
      <a-row>
        <a-col :span="12" />
        <a-col :span="12" />
      </a-row>
    </a-col>
    <a-col :span="3">
      关&nbsp;&nbsp;&nbsp;注：{{ user.follow_count }}<br />
      粉&nbsp;&nbsp;&nbsp;丝：{{ user.followed_count }}<br />
      积&nbsp;&nbsp;&nbsp;分：{{ user.Score }}<br />
      文&nbsp;&nbsp;&nbsp;章：{{ user.article_count }}<br />
      瞬&nbsp;&nbsp;&nbsp;间：{{ user.moment_count }}<br />
      日记本：{{ user.diary_book_count }}<br />
      日&nbsp;&nbsp;&nbsp;记：{{ user.diary_count }}<br />
      <nuxt-link to="/like">
        <a-button icon="book"> 喜欢 </a-button>
      </nuxt-link>
      <br />
      <nuxt-link to="/collection">
        <a-button icon="book"> 收藏 </a-button>
      </nuxt-link>
    </a-col>
  </a-row>
</template>

<script>
export default {
  middleware: "auth",
  data() {
    return {
      formItemLayout: {
        labelCol: { span: 3 },
        wrapperCol: { span: 15 },
      },
      edu_exps: [],
      work_exps: [],
      loading: false,
      edit: false,
      options: [
        {
          value: "zhejiang",
          label: "Zhejiang",
          children: [
            {
              value: "hangzhou",
              label: "Hangzhou",
              children: [
                {
                  value: "xihu",
                  label: "West Lake",
                },
              ],
            },
          ],
        },
        {
          value: "jiangsu",
          label: "Jiangsu",
          children: [
            {
              value: "nanjing",
              label: "Nanjing",
              children: [
                {
                  value: "zhonghuamen",
                  label: "Zhong Hua Men",
                },
              ],
            },
          ],
        },
      ],
    };
  },
  async asyncData({ $axios }) {
    const res = await $axios.$get(`/api/user/edit`);
    return {
      user: res.data,
      address: res.data.address.split(" "),
    };
  },
  created() {
    this.user.birthday = this.$toDayjs(
      this.user.birthday
    );
    for (const v of this.user.edu_exps) {
      this.edu_exps.push({
        id: v.id,
        school: v.school,
        speciality: v.speciality,
        start_time: this.$toDayjs(v.start_time),
        end_time: this.$toDayjs(v.end_time),
        time: [
          this.$toDayjs(v.start_time),
          this.$toDayjs(v.end_time)
        ],
      });
    }

    for (const v of this.user.work_exps) {
      this.work_exps.push({
        id: v.id,
        company: v.company,
        title: v.title,
        start_time: this.$toDayjs(v.start_time),
        end_time: this.$toDayjs(v.end_time),
        time: [
          this.$toDayjs(v.start_time),
          this.$toDayjs(v.end_time),
        ],
      });
    }
  },
  methods: {
    getStatus: function () {
      fetch("https://hoper.xyz/user/1", {
        method: "get",
        headers: {
          "Content-Type": "application/json",
        },
      }).then((res) => {
        return res.json().status;
      });
    },
    uploadAvatarChange(info) {
      if (info.file.status === "uploading") {
        this.loading = true;
        return;
      }
      if (info.file.status === "done") {
        this.user.avatar_url = info.file.response.data.url;
        this.loading = false;
      }
    },
    uploadCoverChange(info) {
      if (info.file.status === "uploading") {
        this.loading = true;
        return;
      }
      if (info.file.status === "done") {
        this.user.cover_url = info.file.response.data.url;
        this.loading = false;
      }
    },
    beforeUpload(file) {
      const isImg = /image\//.test(file.type);
      if (!isImg) {
        this.$message.error("只能上传图片!");
      }
      const isLt2M = file.size / 1024 / 1024 < 4;
      if (!isLt2M) {
        this.$message.error("不能超过 4MB!");
      }
      return isImg && isLt2M;
    },
    changeEdu(index) {
      if (index < 0 && this.edu_exps.length < 5) {
        this.edu_exps.push({
          id: 0,
          name: "",
          speciality: "",
          start_time: "",
          end_time: "",
          time: [],
        });
        return;
      } else {
        this.$message.warning("添加过多!");
      }
      if (index >= 0 && this.edu_exps.length > 0) {
        this.edu_exps.splice(index, 1);
      }
    },
    changeWork(index) {
      if (index < 0 && this.work_exps.length < 5) {
        this.work_exps.push({
          id: 0,
          company: "",
          title: "",
          start_time: "",
          end_time: "",
          time: [],
        });
        return;
      } else {
        this.$message.warning("添加过多!");
      }

      if (index >= 0 && this.work_exps.length > 0) {
        this.work_exps.splice(index, 1);
      }
    },
    editable() {
      this.edit = true;
      if (this.edu_exps.length === 0) {
        this.changeEdu(-1);
      }
      if (this.work_exps.length === 0) {
        this.changeWork(-1);
      }
    },
    async commit() {
      const res = await this.$axios.$put(`/api/user`, {
        sex: this.user.sex,
        birthday: this.user.birthday,
        introduction: this.user.introduction,
        signature: this.user.signature,
        avatar_url: this.user.avatar_url,
        cover_url: this.user.cover_url,
        address: this.user.address,
        edu_exps: this.edu_exps,
        work_exps: this.work_exps,
      });
      if (res !== undefined && res.code === 200) {
        this.$store.commit("SET_USER", res.data);
        this.edit = false;
        this.$message.info("修改成功");
      } else {
        this.$message.info("修改失败");
      }
    },
    eduTime(index) {
      this.edu_exps[index].start_time = this.edu_exps[index].time[0];
      this.edu_exps[index].end_time = this.edu_exps[index].time[1];
    },
    workTime(index) {
      this.work_exps[index].start_time = this.work_exps[index].time[0];
      this.work_exps[index].end_time = this.work_exps[index].time[1];
    },
    addrChange(value) {
      this.user.address = value.join(" ");
    },
  },
};
</script>

<style scoped></style>
