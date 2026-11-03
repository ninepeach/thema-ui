# thema-ui

A starter kit for [Thema](https://github.com/ninepeach/thema) sites: a
browser reset, design tokens, CSS for every common component, Go
`html/template` components with their Go data structs, and a tiny JS helper.

thema-ui is **not a library to keep in sync**. A new project copies it with
`cp -R` and owns the copy from then on — templates, CSS, JS, locale text and
Go structs can all be changed. thema-ui itself is where the base theme is
developed and debugged, and it ships the build tool projects run before
deploying. See [docs/design.md](docs/design.md) for the full design.

```
theme/                          base theme → project themes/<project>/
├── bundle.json                 which CSS/JS files are merged, in order
├── templates/layouts/          document-open / document-close / head
├── templates/components/       global · layout · navigation · content · form · feedback · overlay
├── assets/src/css/             default · common · pages · tokens · style
├── assets/src/js/thema-ui.js   optional enhancements (~100 lines)
└── locales/                    ui.* keys (en, ja, ko, zh-Hant, zh-Hans)
ui/                             component data (ui.Card, ui.Field, …) → project internal/ui
locale/                         languages, /{locale}/ URLs, switcher data → project internal/locale
cmd/preview/                    preview server for debugging
cmd/build/                      merge + minify (run by projects before deploying)
examples/                       example styles: lodge.css · forest.css · ink.css
docs/                           design.md · components.md
CHANGELOG.md                    every change, so projects can port fixes by hand
```

## Develop

```sh
go run ./cmd/preview        # http://localhost:8090/ — live from theme/ and examples/
go test ./...
```

The home page picks an **example style**, a page (website · app · components)
and a language. Toolbar: **example** `lodge` / `forest` / `ink` ·
**assets** `src` (files one by one) / `dist` (built on the fly, to check the
minified output) · **layers** `base` (default + common + pages, neutral
defaults) → `tokens` → `style` → `example`.
`/frames/` shows 340 / 390 / 768 / 1280 px side by side. Edit `theme/` or
`examples/` and refresh.

Record every change in `CHANGELOG.md`; tag versions with git (`v0.1.0` …).
`theme/assets/dist/` is not committed.

## Start a project

```sh
# 1. copy the theme and the Go packages
cp -R thema-ui/theme   myproject/themes/myproject
cp -R thema-ui/ui      myproject/internal/ui
cp -R thema-ui/locale  myproject/internal/locale
# 2. add themes/myproject/theme.json; fix the import paths of internal/ui
# 3. add themes/myproject/assets/css/myproject.css (start from examples/), list it in bundle.json
# 4. write templates/pages/*, composing components; add project keys to locales/
# 5. before deploying:
go run github.com/ninepeach/thema-ui/cmd/build@v0.1.0 -theme ./themes/myproject
```

Note in the project's README which thema-ui version it was copied from; port
later fixes by hand from `CHANGELOG.md`.

## Build and deploy

`cmd/build` merges the files listed in `bundle.json` in order and writes
`assets/dist/<name>.css`, `<name>.min.css`, `<name>.js`, `<name>.min.js`
(pure Go, no Node). Add your own CSS to `bundle.json` to ship one file.

```sh
go run github.com/ninepeach/thema-ui/cmd/build@v0.1.0 -theme ./themes/myproject          # write assets/dist/
go run github.com/ninepeach/thema-ui/cmd/build@v0.1.0 -theme ./themes/myproject -check   # is it up to date?
```

Either commit `dist/` and run `-check` in CI, or generate it in the
Dockerfile build stage without committing it.

The page decides what to load. Copy `bundle.go` (`ReadBundle`, `Assets`)
into the project, or import it:

```go
bundle, _ := themaui.ReadBundle("./themes/myproject")
doc.CSS, doc.JS = bundle.Assets(devMode) // dev: src files one by one; prod: dist/<name>.min.*
```

See [docs/components.md](docs/components.md) for class naming and the component list.
