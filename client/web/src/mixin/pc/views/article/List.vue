<template>
  <a-row>
    <a-col :span="2">
      <router-link to="/article/edit?op=add"> 添加 </router-link>
    </a-col>
    <a-col :span="20">
      <a-breadcrumb>
        <a-breadcrumb-item>主页</a-breadcrumb-item>
        <a-breadcrumb-item>
          <router-link to=""> 博客 </router-link>
        </a-breadcrumb-item>
        <a-breadcrumb-item>列表</a-breadcrumb-item>
      </a-breadcrumb>
      <a-list
        item-layout="vertical"
        size="default"
        :data-source="listState.articleList"
        split
      >
        <template #footer>
          <b />
        </template>
        <template #renderItem="{ item }">
          <a-list-item :key="item.title">
            <template #actions>
              <span @click="star(item.id)">
                <star-outlined style="margin-right: 8px" />
                {{ item.collect_count }}
              </span>
              <span @click="like(item.id)">
                <like-outlined style="margin-right: 8px" />
                {{ item.like_count }}
              </span>
              <span @click="comment(item.id)">
                <message-outlined style="margin-right: 8px" />
                {{ item.comment_count }}
              </span>
              <span v-if="userStore.auth && userStore.auth.id === item.user.id">
                <edit-outlined style="margin-right: 8px" />
                <router-link
                  :to="'/article/edit?id=' + item.id"
                  style="color: rgba(0, 0, 0, 0.45)"
                >
                  编辑
                </router-link>
              </span>
              <a-divider type="vertical" />
              <a-tag
                v-for="(subitem, subindex) in item.categories"
                :key="subindex"
              >
                {{ subitem.name }}
              </a-tag>
              <a-divider type="vertical" />
              <a-checkable-tag
                v-for="(subitem, subindex) in item.tags"
                :key="subindex"
                :color="color[subindex]"
              >
                {{ subitem.name }}
              </a-checkable-tag>
            </template>
            <template #extra>
              <img
                v-if="item.image_url !== ''"
                height="120"
                alt="logo"
                :src="STATIC_DIR + item.image_url"
              />
            </template>
            <a-list-item-meta :description="item.intro">
              <template #title>
                <a-row type="flex">
                  <a-col :flex="1" style="font-size: 10px">
                    <router-link :to="'/user/' + item.user.id">
                      {{ item.user.name }}
                    </router-link>
                  </a-col>

                  <a-col :flex="3">
                    <router-link :to="'/article/' + item.id">
                      <p style="color: rgba(0, 0, 0, 0.85)">{{ item.title }}</p>
                    </router-link>
                  </a-col>

                  <a-col :flex="1" style="font-size: 10px">
                    <span> {{ date2s(item.created_at) }}</span>
                    <a-divider type="vertical" />
                    <span>{{ s2date(item.created_at).fromNow() }}</span>
                  </a-col>
                </a-row>
              </template>
              <template #avatar>
                <router-link :to="'/user/' + item.user.id">
                  <a-avatar
                    :src="STATIC_DIR + item.user.avatar_url"
                    alt="头像"
                  />
                </router-link>
              </template>
            </a-list-item-meta>
          </a-list-item>
        </template>
      </a-list>
      <a-modal v-model:visible="favState.visible" title="Title" @ok="handleOk">
        <template #footer>
          <a-button key="back" @click="handleCancel"> 取消 </a-button>
          <a-button
            key="submit"
            type="primary"
            :loading="favState.loading"
            @click="handleOk"
          >
            确定
          </a-button>
        </template>
        <a-form-item
          label="收藏夹"
          required
          :label-col="{ span: 4 }"
          :wrapper-col="{ span: 6 }"
        >
          <a-select
            v-model="favState.favorites"
            mode="multiple"
            placeholder="请选择收藏夹"
            style="width: 200px"
          >
            <a-select-option
              v-for="item in favState.existFavorites"
              :key="item.id"
            >
              {{ item.name }}
            </a-select-option>
          </a-select>
        </a-form-item>

        <a-row>
          <a-col :span="16">
            <a-form-item
              label="新收藏夹"
              :label-col="{ span: 6 }"
              :wrapper-col="{ span: 16 }"
            >
              <a-input v-model="favState.favorite" />
            </a-form-item>
          </a-col>
          <a-col :span="8">
            <a-button style="margin-top: 5px" @click="addFavorite">
              添加
            </a-button>
          </a-col>
        </a-row>
      </a-modal>
      <a-pagination
        v-model="listState.current"
        :page-size-options="listState.pageSizeOptions"
        :total="listState.total"
        show-quick-jumper
        show-size-changer
        :page-size="listState.pageSize"
        @showSizeChange="list"
        @change="list"
      >
        <template v-slot:buildOptionText="props">
          <span v-if="props.value !== '50'">{{ props.value }}条/页</span>
          <span v-if="props.value === '50'">全部</span>
        </template>
      </a-pagination>
    </a-col>
    <a-col :span="2"> col-8 </a-col>
  </a-row>
</template>

<script setup lang="ts">
import { onMounted, onUpdated, reactive, ref } from "vue";
import type { Ref } from "vue";
import axios from "axios";
import { message } from "ant-design-vue";
import {
  StarOutlined,
  LikeOutlined,
  MessageOutlined,
  EditOutlined,
} from "@ant-design/icons-vue";
import { useUserStore } from "@/mixin/store/user";
import ArticleClient from "@/mixin/service/article";
import { date2s, s2date } from "../../../../../../../thirdparty/diamond/src/time/time";
import { STATIC_DIR } from "@/mixin/plugin/config";
import { tagColor } from "@/mixin/pc/views/article/const";

const userStore = useUserStore();

const actions: Record<string, string>[] = [
  { type: "StarOutlined", attr: "collect_count" },
  { type: "LikeOutlined", attr: "like_count" },
  { type: "MessageOutlined", attr: "comment_count" },
];

const listState: any = reactive({
  pageSizeOptions: ["5", "10", "15", "20"],
  pageSize: 5,
  current: 1,
  loading: false,
  articleList: [],
  total: 0,
  refId: 0,
});
const color: Ref<string[]> = ref(tagColor);
const favState: any = reactive({
  visible: false,
  loading: false,
  favorites: [],
  existFavorites: [],
  favorite: "",
});

await list(listState.current, listState.pageSize);
async function list(current: number, pageSize: number) {
  listState.pageSize = pageSize;
  const { articleList, total } = await ArticleClient.list(
    current - 1,
    listState.pageSize,
  );
  listState.articleList = articleList;
  listState.total = total;
}

async function star(id) {
  listState.refId = id;
  favState.visible = true;
  if (favState.existFavorites.length > 0) {
    return;
  }
  const { data } = await axios.get(`/api/favorites`);
  if (data !== undefined) {
    favState.existFavorites = data.data;
    favState.favorites.push(favState.existFavorites[0].id);
  } else {
    message.error("无法获取收藏夹");
  }
}

async function handleOk(e: MouseEvent) {
  favState.loading = true;
  const params = {
    ref_id: listState.refId,
    kind: "Moment",
    favorites_ids: favState.favorites,
  };
  const { data } = await axios.post("/api/collection", params);
  if (data.code === 200) {
    message.info("收藏成功");
  }
  favState.loading = false;
  favState.visible = false;
}

function handleCancel(e) {
  favState.visible = false;
}
function like(id) {
  listState.refId = id;
}
function comment(id) {
  listState.refId = id;
}
function addFavorite() {
  if (favState.favorite === "") {
    message.error("标签为空");
    return;
  }
  for (const v of favState.existFavorites) {
    if (v.name === favState.tag) {
      message.error("标签重复");
      return;
    }
  }
  favState.existFavorites.push({ name: favState.tag });
  favState.favorites.push(favState.favorite);
  favState.favorite = "";
}
</script>

<style scoped></style>
