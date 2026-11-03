# 组件与约定

## 文件

| 文件 | 内容 |
| --- | --- |
| `src/css/default.css` | 浏览器 reset：只消除浏览器差异，不含颜色和组件 |
| `src/css/common.css` | 开头是全部变量的中性默认值，然后是排版、工具类、全部组件。和 default 一起就能完整使用 |
| `src/css/tokens.css` | thema-ui 的风格：只写和中性默认值不同的变量。项目在自己的 CSS 里覆盖变量 |
| `src/css/pages.css` | 页面骨架：`.ui-main`（网站页）、`body.ui-app` + `.ui-sheet`（手机优先页） |
| `src/css/style.css` | 美化：字体、阴影、过渡、hover / focus。不定义颜色 |
| `src/js/thema-ui.js` | 可选增强：打开 dialog、复制、菜单外点关闭、提交锁定、显示文件名 |
| `bundle.json` | 打包清单：哪些 CSS / JS 按什么顺序合并。项目可以把自己的 CSS 加进来 |
| `assets/css/<项目>.css` | 项目 CSS：覆盖变量 + 少量微调。可从 `examples/`（lodge · forest · ink）拷一个起步 |
| `dist/*` | 项目部署前由 `cmd/build` 按 `bundle.json` 生成，不要手改。thema-ui 仓库里不提交 |

## 类名规范

`ui-<块>__<部件>--<变体>`，状态用 `is-*`。

```html
<button class="ui-button ui-button--primary is-loading">…</button>
<div class="ui-field is-invalid">
  <label class="ui-field__label">…</label>
  <input class="ui-input" aria-invalid="true">
  <p class="ui-field__error">…</p>
</div>
```

- 块：`ui-button`、`ui-card`、`ui-list-card` …（都以 `ui-` 开头，不会和项目类名冲突）
- 部件：`ui-card__title`、`ui-hero__lead`
- 变体：`ui-button--primary`、`ui-hero--compact`、`ui-notice--danger`
- 状态：`is-loading`、`is-invalid`、`is-done`；当前项用 `aria-current`
- 项目自己的类名不要用 `ui-` 前缀

## 模板与数据

- 每个组件：`templates/components/<组>/<名>.html` 里 `{{ define "components/<组>/<名>" }}`，数据是 `ui` 包里的一个结构体。
- 文字一律是 `ui.Text`：`ui.K("key")` 翻译，`ui.S("文字")` 原样输出，`ui.KV("key", "21:00")` 带 `{{value}}`。模板里用 `{{ template "text" .X }}`，图标用 `{{ template "icon" "name" }}`。
- 需要整页上下文的组件（site-header、hero、site-footer、breadcrumb、switch-lang、logo、skip-link）读根对象 `ui.Document`。
- 组件的结构以模板和 `ui` 包为准；`go run ./cmd/preview` 的 components 页展示所有组件和状态。

## 组件清单

稳定：结构和类名已定，改动要记入 `CHANGELOG.md`。实验：可能还会调整。

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

## 注意

- 写在 CSS 变量里的图片地址要用根路径，例如 `--hero-image: url("/assets/myproject/images/hero.jpg")`。相对路径会按“使用这个变量的 CSS 文件”解析（`src/css/` 和 `dist/` 深度不同），不是按定义它的文件。
- 原生控件（date、select、checkbox）在 iPhone Safari 上要实际测试。
