# Changelog

thema-ui 是起步套件：项目拷贝后自己拥有，修复不会自动同步。每条记录写清改了哪个文件、为什么、已有项目是否需要跟进，方便手动搬运。

## v0.1.0 — 2026-11-03

第一个版本，按 `docs/design.md` 定稿“拷贝即拥有”的用法。

- 删除 `cmd/install` 和 `themaui.Install` / `themaui.Diff`：项目用 `cp -R` 拷贝 `theme/`、`ui/`、`locale/`。
- 根包 `themaui` 只保留 `FS`、`ReadBundle`、`Bundle.Assets`（`bundle.go` 可直接拷进项目）。
- `cmd/build` 去掉 `-release` 和 `VERSION` 文件，只保留合并、压缩、`-check`；banner 版本取 `-version` 或 `theme.json`。
- 不再提交 `theme/assets/dist/`：项目部署前自己运行 `cmd/build`；预览的 dist 模式临时生成。
- 新增 `examples/lodge.css`、`forest.css`、`ink.css`，取代预览里的 `demo.css` 配色开关。
- CSS 分层改为逐层可用：`common.css` 开头用 `:where(:root)` 定义全部变量的中性默认值，default + common 单独加载就完整可用；`tokens.css` 只保留 thema-ui 自己的风格值；动效变量移到 `style.css`。加载顺序改为 default → common → pages → tokens → style。新增测试 `TestCSSVariables`。
- 新增组件 `global/scroll-progress`（参照 kurmon-web）：`ui.Document.ScrollProgress = true` 时在页面顶部显示滚动进度条；`thema-ui.js` 设置 `--progress`，`common.css` 新增 `.ui-scroll-progress`，颜色和高度由 `--scroll-progress-color`（默认 `--accent`）、`--scroll-progress-h`（4px）控制，`style.css` 加淡光晕；主色偏暗的项目应设一个更亮的进度条颜色（lodge 用金色、forest 用亮绿）。
- 配色与样式改版（参考当下 WordPress 区块主题）：`tokens.css` 改为暖白底 + 白色面板 + 墨色文字 + 靛蓝主色，圆角 8/14/20px，胶囊按钮，页头 72px，头图更高；`style.css` 加大标题字号、收紧字距，正文 17px，eyebrow 前加短线，毛玻璃吸顶页头，卡片更柔和的分层阴影和悬停上浮、图片放大，输入框聚焦光圈，CTA 渐变，头图和 CTA 里的主按钮改为白色胶囊。`common.css` 新增 `--button-radius`（默认 `--radius`）。示例 lodge / forest / ink 补上 `--bg`，forest 和 ink 设 `--button-radius` 保持各自的方角风格。
- 预览新增首页（`/`，下拉选择示例 / 页面 / 语言），工具栏加示例切换，去掉 palette；layers 新增 `tokens` 一层。

项目跟进：新项目按 `README.md` 的步骤起步；`themes/<项目>/assets/dist/` 由项目自己生成。
