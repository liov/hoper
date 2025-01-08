import tinymce from "tinymce/tinymce";
import "tinymce/themes/silver/theme";
import "tinymce/plugins/image";
import "tinymce/plugins/link";
import "tinymce/plugins/code";
import "tinymce/plugins/table";
import "tinymce/plugins/lists";
import "tinymce/plugins/wordcount";
import {getGlobal} from "@hopeio/utils/compatible";

const getTinymce = () => {
    const global = getGlobal();

    return global && global.tinymce ? global.tinymce : null;
};

export { getTinymce };
