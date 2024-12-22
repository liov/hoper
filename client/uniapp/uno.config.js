"use strict";
var __spreadArray = (this && this.__spreadArray) || function (to, from, pack) {
    if (pack || arguments.length === 2) for (var i = 0, l = from.length, ar; i < l; i++) {
        if (ar || !(i in from)) {
            if (!ar) ar = Array.prototype.slice.call(from, 0, i);
            ar[i] = from[i];
        }
    }
    return to.concat(ar || Array.prototype.slice.call(from));
};
var _a, _b, _c;
Object.defineProperty(exports, "__esModule", { value: true });
// uno.config.ts
var unocss_1 = require("unocss");
var unocss_applet_1 = require("unocss-applet");
// @see https://unocss.dev/presets/legacy-compat
// import { presetLegacyCompat } from '@unocss/preset-legacy-compat'
var isMp = (_c = (_b = (_a = process.env) === null || _a === void 0 ? void 0 : _a.UNI_PLATFORM) === null || _b === void 0 ? void 0 : _b.startsWith('mp')) !== null && _c !== void 0 ? _c : false;
var presets = [];
if (isMp) {
    // 使用小程序预设
    presets.push((0, unocss_applet_1.presetApplet)(), (0, unocss_applet_1.presetRemRpx)());
}
else {
    presets.push(
    // 非小程序用官方预设
    (0, unocss_1.presetUno)(), 
    // 支持css class属性化
    (0, unocss_1.presetAttributify)());
}
exports.default = (0, unocss_1.defineConfig)({
    presets: __spreadArray(__spreadArray([], presets, true), [
        // 支持图标，需要搭配图标库，eg: @iconify-json/carbon, 使用 `<button class="i-carbon-sun dark:i-carbon-moon" />`
        (0, unocss_1.presetIcons)({
            scale: 1.2,
            warn: true,
            extraProperties: {
                display: 'inline-block',
                'vertical-align': 'middle',
            },
        }),
    ], false),
    /**
     * 自定义快捷语句
     * @see https://github.com/unocss/unocss#shortcuts
     */
    shortcuts: [['center', 'flex justify-center items-center']],
    transformers: [
        // 启用 @apply 功能
        (0, unocss_1.transformerDirectives)(),
        // 启用 () 分组功能
        // 支持css class组合，eg: `<div class="hover:(bg-gray-400 font-medium) font-(light mono)">测试 unocss</div>`
        (0, unocss_1.transformerVariantGroup)(),
        // Don't change the following order
        (0, unocss_applet_1.transformerAttributify)({
            // 解决与第三方框架样式冲突问题
            prefixedOnly: true,
            prefix: 'fg',
        }),
    ],
    rules: [
        [
            'p-safe',
            {
                padding: 'env(safe-area-inset-top) env(safe-area-inset-right) env(safe-area-inset-bottom) env(safe-area-inset-left)',
            },
        ],
        ['pt-safe', { 'padding-top': 'env(safe-area-inset-top)' }],
        ['pb-safe', { 'padding-bottom': 'env(safe-area-inset-bottom)' }],
    ],
});
/**
 * 最终这一套组合下来会得到：
 * mp 里面：mt-4 => margin-top: 32rpx  == 16px
 * h5 里面：mt-4 => margin-top: 1rem == 16px
 *
 * 另外，我们还可以推算出 UnoCSS 单位与设计稿差别4倍。
 * 375 * 4 = 1500，把设计稿设置为1500，那么设计稿里多少px，unocss就写多少述职。
 * 举个例子，设计稿显示某元素宽度100px，就写w-100即可。
 *
 * 如果是传统方式写样式，则推荐设计稿设置为 750，这样设计稿1px，代码写1rpx。
 * rpx是响应式的，可以让不同设备的屏幕显示效果保持一致。
 */
