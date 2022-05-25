<template>
  <div id="editor_t" />
</template>

<script lang="ts">
import { Options, Vue, prop } from "vue-class-component";
import { upload } from "@/plugin/utils/upload";
import tinymce from "tinymce/tinymce";
import "tinymce/skins/ui/oxide/skin.min.css";
import "tinymce/skins/ui/oxide/content.min.css";
import "tinymce/skins/content/default/content.css";
import "tinymce/themes/silver/theme";
import "tinymce/icons/default";
import "tinymce/plugins/image";
import "tinymce/plugins/link";
import "tinymce/plugins/code";
import "tinymce/plugins/table";
import "tinymce/plugins/lists";
import "tinymce/plugins/contextmenu";
import "tinymce/plugins/wordcount";
import "tinymce/plugins/colorpicker";
import "tinymce/plugins/textcolor";

@Options({
  components: {},
})
export default class Tinymce extends Vue {
  mounted() {
    if (!tinymce.activeEditor) {
      tinymce.init({
        selector: "#editor_t",
        mobile: {
          menubar: true,
          language_url: "//hoper.xyz/static/zh_CN.js",
          language: "zh_CN",
          skin: "oxide",
          height: 650,
          plugins:
            "link lists image code table colorpicker textcolor wordcount contextmenu",
          toolbar:
            "bold italic underline strikethrough | fontsizeselect | forecolor backcolor | alignleft aligncenter alignright alignjustify | bullist numlist | outdent indent blockquote | undo redo | link unlink image code | removeformat",
          branding: false,
          // 此处为图片上传处理函数，这个直接用了base64的图片形式上传图片，
          // 如需ajax上传可参考https://www.tiny.cloud/docs/configure/file-image-upload/#images_upload_handler
          images_upload_handler: async (blobInfo, success, failure) => {
            const url = await upload(
              new File([blobInfo.blob()], blobInfo.filename())
            );
            success(url);
          },
          // images_upload_url: '/api/upload/article'
        },
      });
    } else {
      tinymce.activeEditor.show();
    }
  }
}
</script>

<style scoped lang="less"></style>
