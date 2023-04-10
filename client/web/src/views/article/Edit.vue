<template>
  <div class="article">
    <a-row>
      <a-col :span="5" style="z-index: 0">
        <a-form-item
          label=""
          :label-col="{ span: 4, offset: 2 }"
          :wrapper-col="{ span: 24, offset: 2 }"
        >
          <a-radio-group v-model="editorType" @change="handleChange">
            <a-radio-button value="markdown"> markdown </a-radio-button>
            <a-radio-button value="html"> 富文本 </a-radio-button>
          </a-radio-group>
          <a-button v-if="editorType === 'html'" @click="save">
            <save-outlined />
          </a-button>
        </a-form-item>
      </a-col>
      <a-col :span="12">
        <a-form-item
          label="标题"
          required
          :label-col="{ span: 4 }"
          :wrapper-col="{ span: 18 }"
        >
          <a-input v-model="article.title" placeholder="请输入标题！" />
        </a-form-item>
      </a-col>
      <a-col :span="6">
        <a-form-item
          label="封面"
          :label-col="{ span: 6 }"
          :wrapper-col="{ span: 15 }"
        >
          <a-row>
            <a-col :span="16">
              <a-upload
                name="file"
                action="/api/upload/article"
                :before-upload="beforeUpload"
                :custom-request="customUpload"
                @change="uploadChange"
              >
                <a-button>
                  <upload-outlined />
                  上传封面
                </a-button>
              </a-upload>
            </a-col>
            <a-col :span="8">
              <a-button class="formbuttion" @click="showImage = !showImage">
                <span v-if="!showImage">显示封面</span>
                <span v-if="showImage">不显示封面</span>
              </a-button>
            </a-col>
          </a-row>
        </a-form-item>
      </a-col>
    </a-row>

    <div>
      <img v-if="showImage" :src="imageUrl" />
    </div>
    <div id="tag">
      <a-row>
        <a-col :span="6">
          <a-form-item
            label="分类"
            required
            :label-col="{ span: 4 }"
            :wrapper-col="{ span: 6 }"
          >
            <a-select
              v-model="categories"
              mode="multiple"
              placeholder="请选择分类"
              style="width: 200px"
            >
              <a-select-option v-for="item in existCategories" :key="item.id">
                {{ item.name }}
              </a-select-option>
            </a-select>
          </a-form-item>
        </a-col>
        <a-col :span="6">
          <a-form-item
            label="标签"
            required
            :label-col="{ span: 4 }"
            :wrapper-col="{ span: 6 }"
          >
            <a-select
              v-model="tags"
              mode="multiple"
              placeholder="请选择标签"
              style="width: 200px"
            >
              <a-select-option v-for="item in existTags" :key="item.name">
                {{ item.name }}
              </a-select-option>
            </a-select>
          </a-form-item>
        </a-col>
        <a-col :span="6">
          <a-form-item
            label="新标签"
            :label-col="{ span: 6 }"
            :wrapper-col="{ span: 12 }"
          >
            <a-row>
              <a-col :span="16">
                <a-input v-model="tag" />
              </a-col>
              <a-col :span="6">
                <a-button class="formbuttion" @click="addTag"> 添加 </a-button>
              </a-col>
            </a-row>
          </a-form-item>
        </a-col>
        <a-col :span="6">
          <a-form-item
            label="权限"
            :label-col="{ span: 4 }"
            :wrapper-col="{ span: 6 }"
          >
            <a-select
              v-model="article.permission"
              placeholder="请选择权限"
              :default-value="[0]"
              style="width: 200px"
            >
              <a-select-option :key="0"> 全部可见 </a-select-option>
              <a-select-option :key="1"> 自己可见 </a-select-option>
              <a-select-option :key="2" disabled> 部分可见 </a-select-option>
              <a-select-option :key="3" disabled> 陌生人可见 </a-select-option>
            </a-select>
          </a-form-item>
        </a-col>
      </a-row>
    </div>
    <div id="editor">
      <div id="vditor" v-show="editorType === 'markdown'" />
      <div v-show="editorType === 'html'">
        <textarea id="default" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
//import 'tinymce/skins/ui/oxide/skin.min.css'
//import 'tinymce/skins/ui/oxide/content.min.css'
//import 'tinymce/skins/content/default/content.css'
//import '../../assets/css/content.css'
//不要tinymce的content.css的body属性
import { SaveOutlined, UploadOutlined } from "@ant-design/icons-vue";
import { upload } from "@/plugin/utils/upload";
import {
  nextTick,
  onBeforeUnmount,
  onMounted,
  onUpdated,
  reactive,
  ref,
} from "vue";
import type { Ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import axios from "axios";
import { message } from "ant-design-vue";

const route = useRoute();
const router = useRouter();
const article: Record<string, any> = ref({
  categories: [],
  tags: [],
  permission: 0,
});
if (route.query.id) {
  const { data } = await axios.get(`/api/article/` + route.query.id);
  if (data) article.value = data;
}
const cres = await axios.get(`/api/category`);
const existCategories = ref(cres.data.data);

const params = { pageNo: 0, pageSize: 10 };
const tres = await axios.get(`/api/tag`, { params });
const existTags = ref(tres.data.data);
const editorType = ref("html");
const showImage = ref(false);
const imageUrl = ref("");
const tag = ref("");
const categories = ref([cres.data.data[0].id]);
const tags: Ref<any[]> = ref([]);
const loading = ref(false);
const useDarkMode = window.matchMedia("(prefers-color-scheme: dark)").matches;
const isSmallScreen = window.matchMedia("(max-width: 1023.5px)").matches;
let vditor = undefined;
const tinymceConfig = {
  selector: "textarea#default",
  base_url: "/node_modules/tinymce/",
  promotion: false,
  branding: false,
  language_url: "src/plugin/tinymce/zh_CN.js",
  language: "zh-Hans",
  plugins:
    "preview importcss searchreplace autolink autosave save directionality code visualblocks visualchars fullscreen image link media template codesample table charmap pagebreak nonbreaking anchor insertdatetime advlist lists wordcount help charmap quickbars emoticons",
  editimage_cors_hosts: ["picsum.photos"],
  menubar: "file edit view insert format tools table help",
  toolbar:
    "undo redo | bold italic underline strikethrough | fontfamily fontsize blocks | alignleft aligncenter alignright alignjustify | outdent indent |  numlist bullist | forecolor backcolor removeformat | pagebreak | charmap emoticons | fullscreen  preview save print | insertfile image media template link anchor codesample | ltr rtl",
  toolbar_sticky: true,
  toolbar_sticky_offset: isSmallScreen ? 102 : 108,
  autosave_ask_before_unload: true,
  autosave_interval: "30s",
  autosave_prefix: "{path}{query}-{id}-",
  autosave_restore_when_empty: false,
  autosave_retention: "2m",
  image_advtab: true,
  link_list: [
    { title: "My page 1", value: "https://www.tiny.cloud" },
    { title: "My page 2", value: "http://www.moxiecode.com" },
  ],
  image_list: [
    { title: "My page 1", value: "https://www.tiny.cloud" },
    { title: "My page 2", value: "http://www.moxiecode.com" },
  ],
  image_class_list: [
    { title: "None", value: "" },
    { title: "Some class", value: "class-name" },
  ],
  importcss_append: true,
  file_picker_callback: (callback, value, meta) => {
    /* Provide file and text for the link dialog */
    if (meta.filetype === "file") {
      callback("https://www.google.com/logos/google.jpg", { text: "My text" });
    }

    /* Provide image and alt text for the image dialog */
    if (meta.filetype === "image") {
      callback("https://www.google.com/logos/google.jpg", {
        alt: "My alt text",
      });
    }

    /* Provide alternative source and posted for the media dialog */
    if (meta.filetype === "media") {
      callback("movie.mp4", {
        source2: "alt.ogg",
        poster: "https://www.google.com/logos/google.jpg",
      });
    }
  },
  images_upload_handler: async (blobInfo, success, failure) => {
    const data = await upload("article", blobInfo.blob());
    success(data.url);
  },
  templates: [
    {
      title: "New Table",
      description: "creates a new table",
      content:
        '<div class="mceTmpl"><table width="98%%"  border="0" cellspacing="0" cellpadding="0"><tr><th scope="col"> </th><th scope="col"> </th></tr><tr><td> </td><td> </td></tr></table></div>',
    },
    {
      title: "Starting my story",
      description: "A cure for writers block",
      content: "Once upon a time...",
    },
    {
      title: "New list with dates",
      description: "New List with dates",
      content:
        '<div class="mceTmpl"><span class="cdate">cdate</span><br><span class="mdate">mdate</span><h2>My List</h2><ul><li></li><li></li></ul></div>',
    },
  ],
  template_cdate_format: "[Date Created (CDATE): %m/%d/%Y : %H:%M:%S]",
  template_mdate_format: "[Date Modified (MDATE): %m/%d/%Y : %H:%M:%S]",
  height: 600,
  image_caption: true,
  quickbars_selection_toolbar:
    "bold italic | quicklink h2 h3 blockquote quickimage quicktable",
  noneditable_class: "mceNonEditable",
  toolbar_mode: "sliding",
  contextmenu: "link image table",
  skin: useDarkMode ? "oxide-dark" : "oxide",
  content_css: useDarkMode ? "dark" : "default",
  content_style:
    "body { font-family:Helvetica,Arial,sans-serif; font-size:16px }",
};

onMounted(() => {
  if (route.query.id) {
    editorType.value = article.value.content_type === 0 ? "markdown" : "html";
    categories.value = [];
    for (const v of article.value.categories) {
      categories.value.push(v.id);
    }
    article.value.categories = [];
    tags.value = [];
    for (const v of article.value.tags) {
      tags.value.push(v.name);
    }
    article.value.tags = [];
  }
  nextTick(async () => {
    if (editorType.value == "markdown") {
      const Vditor = await import("vditor/dist/index.min.js");
      import("vditor/dist/index.css");
      vditor = new Vditor.default("vditor", {});
    } else {
      await import("../../plugin/tinymce/tinymce.js");
      window.tinymce.init(tinymceConfig).then((resolve) => {
        if (route.query.id && article.value.content_type === 1) {
          window.tinymce.activeEditor.setContent(article.value.html_content);
        }
      });
    }
  });
});
onBeforeUnmount(() => {
  if (editorType.value == "html") {
    window.tinymce.activeEditor.destroy();
  }
});

async function handleChange(e) {
  editorType.value = e.target.value;
  if (!vditor && editorType.value == "markdown") {
    const Vditor = await import("vditor/dist/index.min.js");
    import("vditor/dist/index.css");
    console.log(Vditor);
    vditor = new Vditor.default("vditor", {});
  }
  // 手动创建有bug，切换路由回来不渲染了,得destroy()了,另一个组件测试可以用show，无语
  // 两种方式，1.插件客户端渲染，2.有window后引入
  /*      if (typeof window !== 'undefined') {
      /!*        const Vditor = require('vditor')
        const vditor = new Vditor('content', {})
        vditor.focus() *!/
      require('../../plugins/filter/tinymce')
      if (e.target.value === 'html' && tinymce.activeEditor === null) {
        tinymce.init(this.init)
      }
    } */
}

function uploadChange(info) {
  if (info.file.status === "uploading") {
    loading.value = true;
    return;
  }
  if (info.file.status === "done") {
    article.value.image_url = info.file.response.data.url;
    // Get this url from response in real world.
    imageUrl.value = info.file.response.data.url;
    loading.value = false;
  }
}

function beforeUpload(file) {
  const isImg = /image\//.test(file.type);
  if (!isImg) {
    message.error("只能上传图片!");
  }
  const isLt2M = file.size / 1024 / 1024 < 4;
  if (!isLt2M) {
    message.error("不能超过 4MB!");
  }
  return isImg && isLt2M;
}

async function customUpload({
  action,
  data,
  file,
  filename,
  headers,
  onError,
  onProgress,
  onSuccess,
  withCredentials,
  classify,
}) {
  const res = await upload(classify, file);
  onSuccess({ data: res, status: 200 }, file);
  file.status = "done";
}
function addTag() {
  if (tag.value === "") {
    message.error("标签为空");
    return;
  }
  for (const v of existTags) {
    if (v.name === tag.value) {
      message.error("标签重复");
      return;
    }
  }
  existTags.value.push({ name: tag.value });
  tags.value.push(tag.value);
  tag.value = "";
}

async function imgAdd(pos, $file) {
  // 第一步.将图片上传到服务器.
  const data = await upload("article", $file);
  // 第二步.将返回的url替换到文本原位置![...](0) -> ![...](url)
  /**
   * $vm 指为mavonEditor实例，可以通过如下两种方式获取
   * 1. 通过引入对象获取: `import {mavonEditor} from ...` 等方式引入后，`$vm`为`mavonEditor`
   * 2. 通过$refs获取: html声明ref : `<mavon-editor ref=md ></mavon-editor>，`$vm`为 `this.$refs.md`
   */
}

async function save(value, render) {
  article.value.contentType = editorType.value;
  if (editorType.value === "markdown") {
    // this.article.html_content = render
    article.value.content = value;
    article.value.content_type = 0;
  } else {
    article.value.html_content = window.tinymce.activeEditor.getContent();
    article.value.content = window.tinymce.activeEditor.getContent({
      format: "text",
    });
    article.value.content_type = 1;
  }

  for (const i of tags.value) {
    article.value.tags.push({ name: i });
  }
  for (const i of categories.value) {
    article.value.categories.push({ id: i });
  }

  const { data } = route.query.id
    ? await axios.put(`/api/article/${route.query.id}`, article)
    : await axios.post(`/api/article`, article);
  // success
  if (data.code === 200) {
    router.push({ path: "/article" });
  } else message.error(data.msg);
}
</script>

<style scoped lang="less"></style>
