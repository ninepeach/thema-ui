# thema-ui 整体设计

状态：设计定稿前的草案 · 2026-11-03
本文描述 thema-ui 的**目标形态**。第 16 节列出从当前代码到目标需要做的改动；确认本文后再改代码。

---

## 1. 定位

thema-ui 是一个**起步套件**（starter kit），不是需要随项目同步升级的库。

- 它提供一套基础主题：浏览器 reset、设计变量、常用组件的 CSS 和 Go 模板、组件数据结构、一个很小的 JS。
- 新项目把它**拷贝**过去，从此由项目完全拥有，在此基础上随意修改。
- thema-ui 自己用来开发和调试这套基础主题，并提供一个压缩工具供项目在部署前使用。

目标：新项目不用从零写按钮、表单、卡片、布局；用组件拼页面，用变量和项目 CSS 换风格。

非目标：
- 不追求修复自动同步到所有项目（见第 14 节）。
- 不做主题继承、组合、发行版、升级命令。
- 不做 JS 框架，不引入 Node。

---

## 2. 使用方式

```
thema-ui（在这里开发、调试基础主题）
   │  cp -R
   ▼
项目（拷贝后完全拥有，随意修改）
   ├── themes/<项目>/        ← 来自 thema-ui/theme
   └── internal/ui、locale   ← 来自 thema-ui/ui、locale
```

一句话规则：**拷过去就是项目的。** 项目里的模板、CSS、JS、文案、Go 结构体都可以改，不需要保持和 thema-ui 一致。

为什么 Go 的 `ui`、`locale` 包也拷贝而不是 import：模板读取的是 `ui.Card`、`ui.Field` 等结构体的字段。模板归项目、会被修改，结构体也应该归项目，一起改；否则 thema-ui 改了字段、项目升级依赖后，项目里改过的模板可能在渲染时出错。

---

## 3. thema-ui 仓库结构

```
thema-ui/
├── theme/                     基础主题（拷到项目 themes/<项目>/）
│   ├── bundle.json            打包清单：哪些 CSS / JS 按什么顺序合并
│   ├── templates/
│   │   ├── layouts/           document-open / document-close / head
│   │   └── components/        global · layout · navigation · content · form · feedback · overlay
│   ├── assets/
│   │   ├── src/css/           default · common · pages · tokens · style
│   │   └── src/js/            thema-ui.js
│   └── locales/               ui.* 文案（en, ja, ko, zh-Hant, zh-Hans）
├── ui/                        组件数据结构（拷到项目 internal/ui）
├── locale/                    语言列表、URL 前缀、语言切换数据（拷到项目 internal/locale）
├── cmd/preview/               调试用预览服务
├── cmd/build/                 合并 + 压缩工具（项目部署前使用）
├── examples/                  示例风格：lodge.css · forest.css · ink.css
├── docs/                      design.md（本文）· components.md
├── CHANGELOG.md               每次改动记录，方便项目手动搬运
└── README.md
```

不提交 `dist/`、不打 release 压缩包。

---

## 4. 主题结构与文件职责

| 路径 | 内容 |
| --- | --- |
| `theme.json` | Thema 主题清单。thema-ui 里没有；项目拷贝后自己建 |
| `bundle.json` | 打包清单，`cmd/build` 读取。项目可改：加自己的 CSS、改输出名 |
| `templates/layouts/` | 文档外壳：`<html>`、`<head>`、页头页尾、脚本 |
| `templates/components/` | 组件模板，按组分目录 |
| `templates/pages/` | 页面。thema-ui 里没有；项目自己写 |
| `assets/src/css/` | CSS 源码（第 5 节） |
| `assets/src/js/` | JS 源码（第 9 节） |
| `assets/css/` | 项目自己的 CSS。thema-ui 里没有 |
| `assets/dist/` | `cmd/build` 生成的合并压缩文件。不手改 |
| `locales/` | 每种语言一个文件（Thema 的限制）。拷贝后 `ui.*` 和项目文案在同一文件里 |

---

## 5. CSS 分层

一层一层叠加，每一层加上去之前页面都是完整可用的：

| 顺序 | 文件 | 职责 | 原则 |
| --- | --- | --- | --- |
| 1 | `default.css` | 浏览器 reset | 只消除浏览器差异；不含颜色、字体、组件。参考 modern-normalize、Tailwind Preflight、Josh Comeau、Andy Bell |
| 2 | `common.css` | 结构与功能 | 开头用 `:where(:root)` 给出**全部变量的中性默认值**（灰黑白、系统字体、适中圆角），然后是排版、工具类、全部组件。和 1 一起加载就完整可用 |
| 3 | `pages.css` | 页面骨架 | `.ui-main`（网站页）、`body.ui-app` + `.ui-sheet`（手机优先页） |
| 4 | `tokens.css` | thema-ui 的风格 | 只写和中性默认值不同的变量（暖白底、墨色文字、靛蓝主色、圆角与胶囊按钮、阴影、头图色） |
| 5 | `style.css` | 美化 | 品牌字体、阴影、过渡、hover / focus；只读取颜色变量，不定义颜色 |
| 6 | 项目 CSS | 项目定制 | 例如 `assets/css/kurmon-checkin.css`：覆盖变量 + 少量微调 |

约定：
- 组件的颜色、圆角、尺寸都读变量，项目优先改变量，其次才写覆盖规则。
- 可改的变量全部列在 `common.css` 开头；`tokens.css` 和项目 CSS 在 `:root` 里覆盖其中一部分。默认值的优先级为零，所以覆盖不依赖加载顺序。
- 测试 `TestCSSVariables` 检查每一层读到的变量都有定义，`tokens.css` 和示例不引入拼错的变量名。
- CSS 变量里的图片地址用根路径，例如 `--hero-image: url("/assets/<主题>/images/hero.jpg")`。相对路径会按“使用变量的 CSS 文件”解析，`src/css/` 和 `dist/` 深度不同会导致图片丢失。

---

## 6. 类名规范

```
ui-<块>__<部件>--<变体>     状态：is-*     当前项：aria-current
```

```html
<button class="ui-button ui-button--primary is-loading">…</button>
<div class="ui-field is-invalid">
  <label class="ui-field__label">…</label>
  <input class="ui-input" aria-invalid="true">
  <p class="ui-field__error">…</p>
</div>
```

- 基础主题的类名都以 `ui-` 开头。
- 项目新增的类名不用 `ui-` 前缀，避免和基础主题混淆。
- AI 生成页面时，要求它使用这些类名和组件模板，而不是自行设计按钮、表单。

---

## 7. 模板与组件

### 约定

- 每个组件：`templates/components/<组>/<名>.html`，定义 `{{ define "components/<组>/<名>" }}`，数据是 `ui` 包里的一个结构体。
- 所有可见文字都是 `ui.Text`：
  - `ui.K("key")`：翻译 key
  - `ui.S("文字")`：原样输出（例如数据库内容）
  - `ui.KV("key", "21:00")`：带 `{{value}}` 插值的翻译
- 模板里用 `{{ template "text" .X }}` 输出文字，`{{ template "icon" "名字" }}` 输出图标。
- 需要整页上下文的组件读根对象 `ui.Document`：site-header、hero、site-footer、breadcrumb、switch-lang、logo、skip-link。其他组件只读自己的结构体。
- 页面只做组合：`layouts/document-open` → 组件 → `layouts/document-close`。

### 组件清单

| 组 | 组件 | 状态 |
| --- | --- | --- |
| global | text · icon · logo · skip-link · scroll-progress | 稳定 |
| layout | site-header · site-footer · hero · page-header · section | 稳定 |
| navigation | link · switch-lang · breadcrumb · back-link | 稳定 |
| navigation | tabs · steps · pagination | 实验 |
| content | card · card-grid · list-card · info-strip · action-card · action-grid · features · media · accordion · cta · empty-state | 稳定 |
| content | code-display | 实验 |
| form | button · field · choice | 稳定 |
| form | file | 实验 |
| feedback | notice · badge · spinner | 稳定 |
| overlay | dialog | 实验 |

稳定：结构和类名已定，改动要记入 CHANGELOG。实验：还可能调整。

---

## 8. Go 包

| 包 | 内容 | 在项目中 |
| --- | --- | --- |
| `ui` | `Document`、`Text`、各组件结构体 | 拷到 `internal/ui`，随项目修改 |
| `locale` | 语言列表、`/{locale}/` 路径、`Resolve`、语言切换选项 | 拷到 `internal/locale` |

`ui.Document` 是所有页面的根对象，项目在自己的页面结构体里嵌入它。`Document.CSS`、`Document.JS` 决定加载哪些资源（第 11 节）。

---

## 9. JavaScript

`thema-ui.js` 约 100 行，是可选增强，不是框架：

- `[data-dialog-open="id"]` 打开 `<dialog>`
- `[data-copy="id"]` 复制文本并提示
- `details[data-menu]`：点外部 / 按 Esc 关闭，同时只开一个
- `form[data-lock]`：提交时禁用按钮、显示加载状态
- 文件输入显示所选文件名
- `[data-scroll-progress]`：顶部滚动进度条（`ui.Document.ScrollProgress = true` 时由 document-open 输出）

不加载它，所有组件仍然可用（dialog 需要它才能打开，滚动进度条保持隐藏）。

---

## 10. 多语言

- 语言：en、ja、ko、zh-Hant、zh-Hans；URL 前缀 `/en/`、`/ja/`…（与 kurmon-web 一致）。
- 基础主题只使用 `ui.*` 命名空间的 key（跳转、菜单、关闭、复制、分页等约 14 个）。
- 项目文案用自己的命名空间，写在同一个语言文件里。

---

## 11. 构建与压缩

压缩是**项目部署时**的事，由项目自己运行：

```sh
go run github.com/ninepeach/thema-ui/cmd/build@v0.1.0 -theme ./themes/<项目>          # 生成 assets/dist/
go run github.com/ninepeach/thema-ui/cmd/build@v0.1.0 -theme ./themes/<项目> -check   # 检查是否最新
```

- 按 `bundle.json` 的顺序合并 CSS、JS，输出 `<name>.css`、`<name>.min.css`、`<name>.js`、`<name>.min.js`。
- 项目把自己的 CSS 加进 `bundle.json`，上线只加载一个 CSS 文件。
- 使用纯 Go 的压缩库，不需要 Node。
- 部署方式二选一：提交 `dist/`，在 CI 里跑 `-check`；或者在 Dockerfile 构建阶段生成，不提交。

页面加载哪份由 Go 决定：

```go
bundle, _ := themaui.ReadBundle("./themes/kurmon-checkin")
doc.CSS, doc.JS = bundle.Assets(devMode)   // 开发：逐个加载 src；生产：dist/<name>.min.*
```

`ReadBundle` / `Assets` 很小，也一起拷进项目，不依赖 thema-ui 运行时。

---

## 12. 预览与示例

`go run ./cmd/preview`，访问 `http://localhost:8090/`：

- 首页：下拉选择**示例风格**、页面、语言。
- 页面：website（网站）、app（手机优先页）、components（所有组件和状态）。
- 示例风格：`examples/lodge.css`（深蓝、圆角、照片头图）、`forest.css`（墨绿、小圆角、衬线标题）、`ink.css`（黑白、方角、无头图）。每个只是一个项目 CSS 文件，用来示范“只改变量和少量规则就能换风格”，也可以直接拷到项目里作为起点。
- 工具栏：
  - 资源：`src` 逐个加载源码，`dist` 加载临时生成的压缩版（检查压缩结果）。
  - 层：`base` 只加载 default + common + pages（中性默认），再加 tokens，再加 style，再加示例 CSS。
  - 宽度：`/frames/` 同时显示 340 / 390 / 768 / 1280。
- 修改 `theme/` 或 `examples/` 后刷新即可看到。

---

## 13. 新项目起步步骤

```sh
# 1. 拷主题和 Go 包
cp -R thema-ui/theme              myproject/themes/myproject
cp -R thema-ui/ui                 myproject/internal/ui
cp -R thema-ui/locale             myproject/internal/locale
# 2. 新建 themes/myproject/theme.json，改 import 路径
# 3. 新建 themes/myproject/assets/css/myproject.css（可从 examples/ 拷一个起步），加进 bundle.json
# 4. 写 templates/pages/*，用组件拼页面；locales 里加项目文案
# 5. 部署前：go run github.com/ninepeach/thema-ui/cmd/build@vX -theme ./themes/myproject
```

---

## 14. 维护

- thema-ui 每次改动写进 `CHANGELOG.md`：改了哪个文件、为什么、项目是否需要跟进。
- 项目需要某个修复时，对照 CHANGELOG 手动搬运。
- thema-ui 用 git tag 标版本（`v0.1.0`…），项目记录自己是从哪个版本拷贝的（写在 `theme.json` 旁的说明或项目 README 里）。
- 重新评估的时机：共用的项目到了三四个、手动搬运变得频繁时，再考虑“发行版 + 升级”方式。届时目录结构不需要大改。

---

## 15. 已决定不做的事

| 方案 | 不做的原因 |
| --- | --- |
| 主题继承 / 组合（项目只放差异，构建时合成） | 管理太复杂；Thema 不支持多目录，需要额外的组合步骤 |
| 发行版 + 升级命令（`dist` 主题、`cmd/install`、版本对齐） | 不符合“拷贝后修改”的用法；文案合并、版本对齐、旧文件清理都要额外工具 |
| “冻结、不许改”的基础文件 | 同上 |
| 从模板导出静态 HTML | 项目都用 Go + Thema，预览页更实用，导出物没有使用方 |
| 风格变体（variants）生成多套发行版 | 示例风格只需要一个 CSS 文件，放在 `examples/` 即可 |
| JS 框架、Node 构建链 | 不需要；压缩用纯 Go |

---

## 16. 从当前代码到目标的改动

当前 `~/Desktop/dev/thema-ui` 与本文的差异，确认后按此修改：

| 改动 | 说明 |
| --- | --- |
| 删除 `cmd/install`、`themaui.Install`、`themaui.Diff` 及其测试 | 拷贝用 `cp` |
| 删除已提交的 `theme/assets/dist/`、`VERSION` 驱动的 release 打包 | 预览页的 dist 模式改为临时生成 |
| `cmd/build` 去掉 `-release` | 只保留合并、压缩、`-check` |
| 根包 `themaui` 只保留 `FS`、`ReadBundle`、`Bundle.Assets` | 供预览和项目使用，可拷贝 |
| 新增 `examples/lodge.css`、`forest.css`、`ink.css` | 取代预览里的 `demo.css` 配色开关 |
| 预览首页加下拉选择（示例 / 页面 / 语言），工具栏加示例切换 | 去掉 palette 开关 |
| 新增 `CHANGELOG.md` | 从 v0.1.0 开始记录 |
| 更新 `README.md`、`docs/components.md` | 按本文描述使用方式 |

kurmon-checkin 随后按第 13 节重新起步。

---

## 17. 已知限制与待验证

- 只在 Chromium 上测过；iPhone Safari、Firefox 上的原生控件（date、select、checkbox）、dialog 需要实际测试。
- Thema 一次只加载一个主题目录、每种语言一个文案文件，这是“拷贝即拥有”的前提，也是不做继承的原因之一。
- 修复不会自动进入已有项目（第 14 节）。
